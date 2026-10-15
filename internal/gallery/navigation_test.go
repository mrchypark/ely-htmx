package uiexamples

import (
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestGalleryNavigationIsConsistent(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}
	navPattern := regexp.MustCompile(`(?s)<nav class="header-nav"[^>]*>(.*?)</nav>`)
	for _, path := range []string{"/components", "/components/usage", "/components/themes", "/components/api", "/components/audit", "/components/forms/pininput-otpinput"} {
		w := httptest.NewRecorder()
		s.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		match := navPattern.FindStringSubmatch(w.Body.String())
		if w.Code != 200 || len(match) != 2 {
			t.Fatalf("%s lacks navigation (%d)", path, w.Code)
		}
		if !strings.Contains(w.Body.String(), `<html lang="en">`) {
			t.Errorf("%s must declare English as its default language", path)
		}
		if !strings.Contains(match[0], `aria-label="Explore examples"`) {
			t.Errorf("%s lacks the shared English navigation label", path)
		}
		nav := match[1]
		if strings.Contains(nav, "/components/audit") || !strings.Contains(w.Body.String(), `aria-label="Developer tools"`) {
			t.Fatal("developer tools must remain secondary")
		}
		if strings.Count(nav, "<a ") != 4 {
			t.Fatalf("%s inconsistent links or current page: %s", path, nav)
		}
		active := path
		if strings.HasPrefix(path, "/components/forms/") {
			active = "/components"
		}
		if active != "/components/api" && active != "/components/audit" && !strings.Contains(nav, `href="`+active+`" aria-current="page"`) {
			t.Errorf("%s has incorrect current page", path)
		}
		for _, label := range []string{"Components", "Form examples", "Themes", "Playground"} {
			if !strings.Contains(nav, ">"+label+"</a>") {
				t.Errorf("%s lacks %s", path, label)
			}
		}
	}
}
