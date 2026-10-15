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
		nav := match[1]
		if strings.Contains(nav, "/components/audit") || !strings.Contains(w.Body.String(), `aria-label="개발자 도구"`) {
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
		for _, label := range []string{"컴포넌트", "폼 예제", "테마 제안", "플레이그라운드"} {
			if !strings.Contains(nav, ">"+label+"</a>") {
				t.Errorf("%s lacks %s", path, label)
			}
		}
	}
}
