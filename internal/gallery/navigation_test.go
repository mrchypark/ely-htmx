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
		if strings.Count(nav, "<a ") != 5 || strings.Count(nav, `aria-current="page"`) != 1 {
			t.Fatalf("%s inconsistent links or current page: %s", path, nav)
		}
		active := path
		if strings.HasPrefix(path, "/components/forms/") {
			active = "/components"
		}
		if !strings.Contains(nav, `href="`+active+`" aria-current="page"`) {
			t.Errorf("%s has incorrect current page", path)
		}
		for _, label := range []string{"컴포넌트", "폼 예제", "테마 제안", "고급 API", "동작 검증"} {
			if !strings.Contains(nav, ">"+label+"</a>") {
				t.Errorf("%s lacks %s", path, label)
			}
		}
	}
}
