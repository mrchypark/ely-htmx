// Package httpencoding implements HTTP content coding for explicitly public data.
// It is not response middleware: callers must not compress secrets by default.
package httpencoding

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

// Encode encodes one complete representation. identity returns data unchanged.
func Encode(data []byte, coding string) ([]byte, error) {
	o, _ := Normalize(nil, "static")
	return EncodeWithOptions(data, coding, o)
}

// EncodeWithOptions encodes using validated, normalized options.
func EncodeWithOptions(data []byte, coding string, options Options) ([]byte, error) {
	var b bytes.Buffer
	switch coding {
	case "identity":
		return data, nil
	case "gzip":
		w, err := gzip.NewWriterLevel(&b, options.GzipLevel)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(data); err != nil {
			return nil, err
		}
		if err := w.Close(); err != nil {
			return nil, err
		}
	case "br":
		w := brotli.NewWriterLevel(&b, options.BrotliLevel)
		if _, err := w.Write(data); err != nil {
			return nil, err
		}
		if err := w.Close(); err != nil {
			return nil, err
		}
	case "zstd":
		w, err := zstd.NewWriter(nil, zstd.WithEncoderConcurrency(1), zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(options.ZstdLevel)))
		if err != nil {
			return nil, err
		}
		defer w.Close()
		return w.EncodeAll(data, nil), nil
	default:
		return nil, fmt.Errorf("unsupported content coding %q", coding)
	}
	return b.Bytes(), nil
}

// Negotiate chooses zstd, br, gzip or identity, in server preference order when
// client qualities tie. An absent/empty header conservatively chooses identity.
func Negotiate(header string) (string, bool) {
	return NegotiateAvailable(header, "zstd", "br", "gzip", "identity")
}

// NegotiateAvailable follows RFC 9110 section 12.5.3. available is ordered by
// server preference. Identity is an implicit fallback unless explicitly ranked
// or excluded; *;q=0 also excludes implicit identity. Malformed qualities are
// treated as zero and conflicting duplicates use the least permissive quality.
func NegotiateAvailable(header string, available ...string) (string, bool) {
	qualities := make(map[string]int)
	for _, item := range strings.Split(header, ",") {
		parts := strings.Split(item, ";")
		coding := strings.ToLower(strings.TrimSpace(parts[0]))
		if coding == "" {
			continue
		}
		q := 1000
		if len(parts) > 1 {
			key, value, found := strings.Cut(strings.TrimSpace(parts[1]), "=")
			if len(parts) != 2 || !found || !strings.EqualFold(strings.TrimSpace(key), "q") {
				q = 0
			} else {
				q = quality(strings.TrimSpace(value))
			}
		}
		if previous, ok := qualities[coding]; !ok || q < previous {
			qualities[coding] = q
		}
	}
	best, bestQ := "", 0
	identity := false
	for _, coding := range available {
		if coding == "identity" {
			identity = true
		}
		q, explicit := qualities[coding]
		if !explicit {
			if coding == "identity" {
				continue
			}
			q = qualities["*"]
		}
		if q > bestQ {
			best, bestQ = coding, q
		}
	}
	if bestQ > 0 {
		return best, true
	}
	if identity {
		q, explicit := qualities["identity"]
		wildcard, wildcardSet := qualities["*"]
		if (explicit && q > 0) || (!explicit && (!wildcardSet || wildcard > 0)) {
			return "identity", true
		}
	}
	return "", false
}

func quality(s string) int {
	if s == "0" {
		return 0
	}
	if s == "1" {
		return 1000
	}
	if len(s) < 2 || len(s) > 5 || (s[0] != '0' && s[0] != '1') || s[1] != '.' {
		return 0
	}
	digits := s[2:]
	for _, c := range digits {
		if c < '0' || c > '9' || (s[0] == '1' && c != '0') {
			return 0
		}
	}
	if s[0] == '1' {
		return 1000
	}
	n, _ := strconv.Atoi(digits + strings.Repeat("0", 3-len(digits)))
	return n
}

// AddVary preserves existing Vary values, including Vary: *.
func AddVary(h http.Header, name string) {
	for _, line := range h.Values("Vary") {
		for _, token := range strings.Split(line, ",") {
			if strings.EqualFold(strings.TrimSpace(token), name) || strings.TrimSpace(token) == "*" {
				return
			}
		}
	}
	h.Add("Vary", name)
}
