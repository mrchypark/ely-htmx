package playground

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSharedEnglishChrome(t *testing.T) {
	h, err := NewHandler()
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/", "/assets/site-chrome.css", "/assets/theme.js", "/assets/ely/fonts/Inter-Regular.ttf"} {
		req := httptest.NewRequest("GET", "http://localhost"+path, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatalf("%s: %d", path, w.Code)
		}
		if path == "/" {
			for _, want := range []string{`lang="en"`, `Ely HTMX · Playground`, `site-header`, `Explore examples`, `Form examples`, `Themes`, `data-theme-toggle`, `/assets/site-chrome.css`} {
				if !strings.Contains(w.Body.String(), want) {
					t.Errorf("missing %s", want)
				}
			}
		}
	}
}
