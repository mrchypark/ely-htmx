package uiexamples

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestThemeProposalPage(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	server.ServeHTTP(response, httptest.NewRequest("GET", "/components/themes", nil))
	if response.Code != 200 {
		t.Fatalf("status = %d", response.Code)
	}
	for _, marker := range []string{`id="theme-preset"`, `value="blue"`, `id="theme-json"`, `data-theme-preview="light"`, `data-theme-preview="dark"`, `/assets/theme.js`, `/assets/themes.js`, `data-ui-control-template="input"`} {
		if !strings.Contains(response.Body.String(), marker) {
			t.Errorf("missing %s", marker)
		}
	}
	for _, path := range []string{"/assets/theme.js", "/assets/themes.js", "/assets/themes.css"} {
		asset := httptest.NewRecorder()
		server.ServeHTTP(asset, httptest.NewRequest("GET", path, nil))
		if asset.Code != 200 {
			t.Errorf("%s status = %d", path, asset.Code)
		}
	}
}
