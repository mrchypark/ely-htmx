package sdui

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Request contains authenticated host context and parsed input. Files are valid
// only until the action returns; applications must copy uploads they retain.
type Request struct {
	Principal   string
	Values      url.Values
	Files       map[string][]*multipart.FileHeader
	HTTPRequest *http.Request
}
type Loader func(context.Context, Request) (map[string]any, error)
type Result struct {
	Data   map[string]any
	Status int
}
type Action struct {
	Targets   []string
	Authorize func(context.Context, Request) error
	Run       func(context.Context, Request) (Result, error)
}

// Config requires host-owned authentication, authorization, and CSRF issuance.
// The in-memory replay cache is a bounded single-process convenience, not a
// durable business transaction or a replacement for database idempotency.
type Config struct {
	Endpoint              string
	Authenticate          func(*http.Request) (string, error)
	Authorize             func(context.Context, string) error
	CSRF                  func(*http.Request) (string, error)
	Loaders               map[string]Loader
	Actions               map[string]Action
	MaxBodyBytes          int64
	IdempotencyTTL        time.Duration
	MaxIdempotencyEntries int
	MaxReplayBytes        int64
}
type replay struct {
	fingerprint [32]byte
	done        chan struct{}
	body        []byte
	status      int
	expires     time.Time
	reserved    int64
}
type replayCache struct {
	mu          sync.Mutex
	replays     map[string]*replay
	replayBytes int64
}

type Handler struct {
	plan     *Plan
	config   Config
	revision string
	*replayCache
}

var requestKey = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)

func NewHandler(plan *Plan, c Config) (*Handler, error) {
	if plan == nil || plan.root == nil || plan.templates == nil || c.Authenticate == nil || c.Authorize == nil || c.CSRF == nil {
		return nil, fmt.Errorf("plan, authentication, authorization and CSRF callbacks are required")
	}
	if !validEndpoint(c.Endpoint) {
		return nil, fmt.Errorf("endpoint must be a local absolute path")
	}
	if c.MaxBodyBytes == 0 {
		c.MaxBodyBytes = 1 << 20
	}
	if c.MaxBodyBytes < 1 || c.MaxBodyBytes > 64<<20 {
		return nil, fmt.Errorf("MaxBodyBytes must be 1..64 MiB")
	}
	if c.IdempotencyTTL == 0 {
		c.IdempotencyTTL = 10 * time.Minute
	}
	if c.IdempotencyTTL < time.Second || c.IdempotencyTTL > 24*time.Hour {
		return nil, fmt.Errorf("IdempotencyTTL must be 1 second..24 hours")
	}
	if c.MaxIdempotencyEntries == 0 {
		c.MaxIdempotencyEntries = 1024
	}
	if c.MaxIdempotencyEntries < 1 || c.MaxIdempotencyEntries > 100000 {
		return nil, fmt.Errorf("MaxIdempotencyEntries must be 1..100000")
	}
	if c.MaxReplayBytes == 0 {
		c.MaxReplayBytes = 64 << 20
	}
	if c.MaxReplayBytes < MaxRenderBytes || c.MaxReplayBytes > 512<<20 {
		return nil, fmt.Errorf("MaxReplayBytes must be 2..512 MiB")
	}
	actions := make(map[string]Action, len(c.Actions))
	allowed := make(map[string][]string, len(c.Actions))
	for name, a := range c.Actions {
		if a.Run == nil || a.Authorize == nil {
			return nil, fmt.Errorf("action %q needs Run and Authorize", name)
		}
		a.Targets = append([]string(nil), a.Targets...)
		actions[name] = a
		allowed[name] = a.Targets
	}
	if err := plan.ValidateActions(allowed); err != nil {
		return nil, err
	}
	c.Actions = actions
	loaders := make(map[string]Loader, len(c.Loaders))
	for n, l := range c.Loaders {
		loaders[n] = l
	}
	c.Loaders = loaders
	if plan.Loader() != "" && c.Loaders[plan.Loader()] == nil {
		return nil, fmt.Errorf("unknown loader %q", plan.Loader())
	}
	return &Handler{plan: plan, config: c, replayCache: &replayCache{replays: make(map[string]*replay)}}, nil
}
func validEndpoint(path string) bool {
	return path != "" && validateEndpoint(path) == nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Vary", "Cookie, Authorization")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if h.revision != "" {
		w.Header().Set("X-SDUI-Revision", h.revision)
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", 405)
		return
	}
	principal, err := h.config.Authenticate(r)
	if err != nil || principal == "" {
		http.Error(w, "authentication required", 401)
		return
	}
	if err = h.config.Authorize(r.Context(), principal); err != nil {
		http.Error(w, "forbidden", 403)
		return
	}
	csrf, err := h.config.CSRF(r)
	if err != nil || csrf == "" {
		http.Error(w, "CSRF unavailable", 403)
		return
	}
	if r.Method == http.MethodGet {
		h.page(w, r, principal, csrf)
		return
	}
	if r.URL.Path != h.config.Endpoint {
		http.NotFound(w, r)
		return
	}
	h.action(w, r, principal, csrf)
}
func (h *Handler) page(w http.ResponseWriter, r *http.Request, principal, csrf string) {
	data := map[string]any{}
	if load := h.config.Loaders[h.plan.Loader()]; load != nil {
		var err error
		data, err = load(r.Context(), Request{Principal: principal, Values: r.URL.Query(), HTTPRequest: r})
		if err != nil {
			http.Error(w, "invalid page request", 400)
			return
		}
	}
	var b bytes.Buffer
	b.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><title>SDUI</title>")
	if err := h.plan.RenderAssets(&b); err != nil {
		http.Error(w, "render failed", 500)
		return
	}
	b.WriteString("</head><body>")
	if err := h.plan.Render(&b, RenderContext{Data: data, CSRF: csrf, Endpoint: h.config.Endpoint, Revision: h.revision}); err != nil {
		http.Error(w, "render failed", 500)
		return
	}
	b.WriteString("</body></html>")
	writeHTML(w, 200, b.Bytes())
}
func (h *Handler) action(w http.ResponseWriter, r *http.Request, principal, csrf string) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, h.config.MaxBodyBytes))
	if err != nil {
		http.Error(w, "request too large", 413)
		return
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		err = r.ParseMultipartForm(h.config.MaxBodyBytes)
	} else {
		err = r.ParseForm()
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}
	if err != nil {
		http.Error(w, "invalid form", 400)
		return
	}
	// Reserved fields must be unambiguous, and query strings never supply them.
	one := func(key string) string {
		v := r.PostForm[key]
		if len(v) != 1 {
			return ""
		}
		return v[0]
	}
	if subtle.ConstantTimeCompare([]byte(one("_csrf")), []byte(csrf)) != 1 {
		http.Error(w, "invalid CSRF token", 403)
		return
	}
	if h.revision != "" && one("_revision") != h.revision {
		w.Header().Set("X-SDUI-Reload", "required")
		w.Header().Set("X-SDUI-Revision", h.revision)
		http.Error(w, "screen revision changed; reload before submitting", http.StatusConflict)
		return
	}
	name, target := one("_action"), one("_target")
	action, ok := h.config.Actions[name]
	if !ok || !h.allowed(name, target) {
		http.Error(w, "unknown action or target", 400)
		return
	}
	req := Request{Principal: principal, Values: r.PostForm, HTTPRequest: r}
	if r.MultipartForm != nil {
		req.Files = r.MultipartForm.File
	}
	if err := action.Authorize(r.Context(), req); err != nil {
		http.Error(w, "forbidden", 403)
		return
	}
	key := r.Header.Get("Idempotency-Key")
	if len(r.Header.Values("Idempotency-Key")) != 1 || !requestKey.MatchString(key) {
		http.Error(w, "valid Idempotency-Key required", 400)
		return
	}
	// Multipart boundaries vary across retries; hash canonical field/file content.
	fingerprint, err := fingerprintForm(r)
	if err != nil {
		http.Error(w, "invalid upload", 400)
		return
	}
	cacheKey := fmt.Sprintf("%d:%s%d:%s%d:%s%s", len(principal), principal, len(name), name, len(target), target, key)
	h.mu.Lock()
	now := time.Now()
	for k, v := range h.replays {
		select {
		case <-v.done:
			if now.After(v.expires) {
				h.replayBytes -= v.reserved
				delete(h.replays, k)
			}
		default:
		}
	}
	if existing := h.replays[cacheKey]; existing != nil {
		h.mu.Unlock()
		if existing.fingerprint != fingerprint {
			http.Error(w, "idempotency key reused with different input", 409)
			return
		}
		select {
		case <-existing.done:
			writeHTML(w, existing.status, existing.body)
		case <-r.Context().Done():
		}
		return
	}
	if len(h.replays) >= h.config.MaxIdempotencyEntries || h.replayBytes+MaxRenderBytes > h.config.MaxReplayBytes {
		h.mu.Unlock()
		http.Error(w, "replay capacity reached; retry later", 503)
		return
	}
	entry := &replay{fingerprint: fingerprint, done: make(chan struct{}), expires: now.Add(h.config.IdempotencyTTL), reserved: MaxRenderBytes}
	h.replays[cacheKey] = entry
	h.replayBytes += entry.reserved
	h.mu.Unlock()
	// Finalize every outcome, including a host panic, so duplicate requests cannot
	// wait forever. Panics become generic failures without exposing host internals.
	status := 500
	response := []byte("action failed\n")
	defer func() {
		_ = recover()
		h.mu.Lock()
		h.replayBytes += int64(len(response)) - entry.reserved
		entry.reserved = int64(len(response))
		entry.status = status
		entry.body = append([]byte(nil), response...)
		entry.expires = time.Now().Add(h.config.IdempotencyTTL)
		close(entry.done)
		h.mu.Unlock()
		writeHTML(w, status, response)
	}()
	result, err := action.Run(r.Context(), req)
	if err != nil {
		return
	}
	if result.Status == 0 {
		result.Status = 200
	}
	if result.Status != 200 && result.Status != 422 {
		return
	}
	var b bytes.Buffer
	if err = h.plan.RenderRegion(&b, target, RenderContext{Data: result.Data, CSRF: csrf, Endpoint: h.config.Endpoint, Revision: h.revision}); err != nil {
		return
	}
	status, response = result.Status, b.Bytes()
}
func (h *Handler) allowed(action, target string) bool {
	for _, e := range h.plan.Events() {
		if e.Action == action && e.Target == target {
			return true
		}
	}
	return false
}
func writeHTML(w http.ResponseWriter, status int, b []byte) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(b)
}
func fingerprintForm(r *http.Request) ([32]byte, error) {
	h := sha256.New()
	fmt.Fprintf(h, "%d:%s", len(r.URL.RawQuery), r.URL.RawQuery)
	io.WriteString(h, r.PostForm.Encode())
	if r.MultipartForm != nil {
		names := make(url.Values)
		for name := range r.MultipartForm.File {
			names[name] = []string{""}
		}
		for _, pair := range strings.Split(names.Encode(), "&") {
			name, _, _ := strings.Cut(pair, "=")
			name, _ = url.QueryUnescape(name)
			for _, f := range r.MultipartForm.File[name] {
				fmt.Fprintf(h, "\x00%d:%s%d:%s%d:", len(name), name, len(f.Filename), f.Filename, f.Size)
				metadata := url.Values(f.Header).Encode()
				fmt.Fprintf(h, "%d:%s", len(metadata), metadata)
				file, err := f.Open()
				if err != nil {
					return [32]byte{}, err
				}
				_, err = io.Copy(h, file)
				file.Close()
				if err != nil {
					return [32]byte{}, err
				}
			}
		}
	}
	var sum [32]byte
	copy(sum[:], h.Sum(nil))
	return sum, nil
}
