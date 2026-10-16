package playground

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestPreviewThemeIsBoundedAndHashAuthorized(t *testing.T) {
	colors := map[string]string{}
	for _, role := range previewColorRoles {
		colors[role] = "#123456"
	}
	theme := &previewTheme{Mode: "dark", Colors: colors}
	read := func(string) ([]byte, error) { return []byte("body{font-family:Inter,system-ui,sans-serif}"), nil }
	out, err := previewDocumentWithTheme(`<p>Safe preview</p>`, `<link rel="stylesheet" href="/assets/app.css">`, read, theme)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`data-theme="dark"`, `--color-paper:#123456;`, `default-src 'none'`, `img-src 'none'`, `Safe preview`} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %s", want)
		}
	}
	start := strings.Index(out, "<style>") + len("<style>")
	end := strings.Index(out, "</style>")
	hash := sha256.Sum256([]byte(out[start:end]))
	if !strings.Contains(out, "style-src 'sha256-"+base64.StdEncoding.EncodeToString(hash[:])+"'") {
		t.Fatal("theme CSS must be covered by the exact CSP hash")
	}
	for _, bad := range []string{"red", "#123", "#123456;}", "url(https://example.test)", "</style><script>alert(1)</script>"} {
		colors["color-paper"] = bad
		if _, _, err := theme.styles(); err == nil {
			t.Fatalf("accepted %q", bad)
		}
	}
	colors["color-paper"] = "#123456"
	colors["unknown"] = "#123456"
	if _, _, err := theme.styles(); err == nil {
		t.Fatal("accepted unknown color role")
	}
	delete(colors, "unknown")
	delete(colors, "focus")
	if _, _, err := theme.styles(); err == nil {
		t.Fatal("accepted incomplete color roles")
	}
	for _, mode := range []string{"", "system", `dark" onload="alert(1)`} {
		if _, _, err := (&previewTheme{Mode: mode}).styles(); err == nil {
			t.Fatalf("accepted mode %q", mode)
		}
	}
}

func TestPreviewDefaultRemainsCompatible(t *testing.T) {
	mode, css, err := (*previewTheme)(nil).styles()
	if err != nil || mode != "light" || css != "" {
		t.Fatalf("unexpected default: %s %s %v", mode, css, err)
	}
	mode, css, err = (&previewTheme{Mode: "dark"}).styles()
	if err != nil || mode != "dark" || css != "" {
		t.Fatal("mode-only request should use built-in neutral CSS")
	}
}

func TestPreviewAPIThemeAndMeasurements(t *testing.T) {
	server := setup(t)
	screen := json.RawMessage(`{"version":"1","id":"theme-check","root":{"id":"title","type":"text","props":{"text":"Theme test"}}}`)
	var previous map[string]any
	for _, mode := range []string{"light", "dark"} {
		body, _ := json.Marshal(map[string]any{"screen": screen, "theme": map[string]any{"mode": mode}})
		response := request(t, server, "POST", "/api/preview", body, true)
		if response.Code != 200 {
			t.Fatalf("%d %s", response.Code, response.Body)
		}
		var result map[string]any
		if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(result["html"].(string), `data-theme="`+mode+`"`) {
			t.Fatal("API lost theme mode")
		}
		if previous != nil {
			before, _ := json.Marshal(previous["sizes"])
			after, _ := json.Marshal(result["sizes"])
			if string(before) != string(after) {
				t.Fatal("theme chrome changed authored markup measurements")
			}
		}
		previous = result
	}
	bad, _ := json.Marshal(map[string]any{"screen": screen, "theme": map[string]any{"mode": "dark", "css": "body{display:none}"}})
	if response := request(t, server, "POST", "/api/preview", bad, true); response.Code != 422 {
		t.Fatalf("unrecognized theme field accepted: %d", response.Code)
	}
}
