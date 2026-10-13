package uiexamples

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestLocalFormStoriesWireValidationAndHandler(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		chapter, slug string
		required      int
	}{
		{"account", "loginform-oauthbuttons", 2},
		{"account", "signupform", 4}, // Source terms checkbox is required with name, email, and password.
		{"account", "magiclinkform-forgotpassword", 2},
		{"account", "profilecard-profileeditor", 2},
		{"onboarding", "onboardingwizard-welcometour", 1}, // Source invite step is optional.
		{"onboarding", "feedbackwidget-contactsupport", 0},
	} {
		t.Run(tc.slug, func(t *testing.T) {
			w := httptest.NewRecorder()
			server.ServeHTTP(w, httptest.NewRequest(http.MethodGet, componentHref(tc.chapter, tc.slug), nil))
			body := w.Body.String()
			if w.Code != http.StatusOK || !strings.Contains(body, `data-ely-family="forms"`) || !strings.Contains(body, `src="/assets/ely-api-all.js"`) || (!strings.Contains(body, "data-ely-submit") && !strings.Contains(body, "data-ely-wizard")) {
				t.Fatal("local form is not wired to its handler")
			}
			if tc.slug == "signupform" && !regexp.MustCompile(`<input\b[^>]*name="terms"[^>]*required`).MatchString(body) {
				t.Fatal("source signup must require terms consent")
			}
			if tc.slug == "onboardingwizard-welcometour" {
				invite := regexp.MustCompile(`<input\b[^>]*name="email"[^>]*>`).FindString(body)
				if invite == "" || strings.Contains(invite, "required") {
					t.Fatal("source invitation input must exist and remain optional")
				}
				if !regexp.MustCompile(`<input\b[^>]*name="workspace"[^>]*required`).MatchString(body) {
					t.Fatal("workspace name must remain required")
				}
			}
			inputs := regexp.MustCompile(`<input\b[^>]*\srequired(?:\s|>)`).FindAllString(body, -1)
			if len(inputs) != tc.required {
				t.Fatalf("required inputs = %d, want %d", len(inputs), tc.required)
			}
		})
	}
}

func TestInteractiveStoryTargetsAreRendered(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for path, hooks := range map[string][]string{
		"/components/settings/settingslayout-settingssection-settingsrow-settingssearch": {"data-ely-settings-search", "data-ely-settings-section", "save-settings"},
		"/components/primitives/focusscope-focustrap":                                    {`data-foundation-scope`, `class="ely-foundation__focus-scope"`, ">One</button>", ">Two</button>", ">Three</button>"},
		"/components/finance/drawingtools-drawingtoolbar":                                {"data-ely-drawing", "ely-data__drawing"},
		"/components/generative/videogenerationtimeline":                                 {"data-native-shots", "data-native-shot-ruler", "data-native-shot-playhead", "data-native-shot-caption", `data-native-shot="s4"`, `data-state="failed"`, `data-ely-action="extend-video"`},
		"/components/chat/promptinput-composer-attachmentbutton-attachmentchip-attachmentpreview-pasteimagepreview-dragdropoverlay-contextchips-slashcommandmenu-contextmentionmenu-sendbutton-inputhint": {`data-ely-submit="send-message"`, `name="message"`, `data-ely-prompt-editor`, "data-ely-command"},
	} {
		t.Run(path, func(t *testing.T) {
			w := httptest.NewRecorder()
			server.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d", w.Code)
			}
			for _, hook := range hooks {
				if !strings.Contains(w.Body.String(), hook) {
					t.Errorf("missing interactive target %q", hook)
				}
			}
		})
	}
}
