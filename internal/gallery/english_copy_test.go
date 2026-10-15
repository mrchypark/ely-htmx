package uiexamples

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGalleryDefaultEnglishCopy(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		path string
		want []string
	}{
		{"/components", []string{"From small details", "Search components", "All categories", "No matching examples.", "Open playground"}},
		{"/components/usage", []string{"Form examples", "Research topics", "Inspect submitted values", "Submitted values are displayed only on this page."}},
		{"/components/themes", []string{"Find your palette.", "Theme settings", "Apply theme", "Restore default", "Import JSON"}},
		{"/components/api", []string{"Advanced component API", "Actions from these examples appear here."}},
		{"/components/audit?family=forms", []string{"Component behavior checks", "Run all behavior checks", "instance pairs waiting", `aria-label="Check results"`}},
		{"/components/forms/input-textfield-clearableinput", []string{"Try the live component below.", `aria-label="Component list"`, `aria-label="Previous and next examples"`}},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			server.ServeHTTP(response, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d", response.Code)
			}
			for _, want := range append(tc.want, `<html lang="en">`, `aria-label="Explore examples"`) {
				if !strings.Contains(response.Body.String(), want) {
					t.Errorf("missing English UI copy %q", want)
				}
			}
		})
	}
}

func TestGalleryInputErrorsAreEnglish(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for path, want := range map[string]string{
		"/components":      "Search queries must be 200 characters or fewer.",
		"/components/demo": "Input must be 200 characters or fewer.",
	} {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path+"?q="+strings.Repeat("a", 201), nil))
		if response.Code != http.StatusBadRequest || strings.TrimSpace(response.Body.String()) != want {
			t.Errorf("%s: status=%d body=%q", path, response.Code, response.Body.String())
		}
	}
}
