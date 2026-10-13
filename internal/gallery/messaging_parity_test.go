package uiexamples

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMessagingSourceStories(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for slug, wants := range map[string][]string{
		"channellist-chatlist":          {"general", "Dev Rao", `data-ely-unread>12`, `aria-selected="true"`},
		"channelitem-directmessageitem": {"atrium", "AL", "DR", `data-presence="offline"`},
		"channelheader":                 {"12 members", "3 pins", "start-call"},
		"pinnedmessages":                {"3 hours ago", "Please leave notes on the backs.", "The brief, final:"},
		"chatmessage-unreaddivider-messagereactions-typingindicator": {"09:05", "09:07", "09:21", "Thoughts on where the trees go?", "🌿"},
		"messagethread-threadpanel":                                  {"Thoughts on where the trees go?", "Two there, and one on the landing.", "10:12", "data-ely-thread-close"},
	} {
		t.Run(slug, func(t *testing.T) {
			r := httptest.NewRecorder()
			s.ServeHTTP(r, httptest.NewRequest("GET", "/components/messaging/"+slug, nil))
			if r.Code != 200 {
				t.Fatalf("status %d", r.Code)
			}
			for _, want := range wants {
				if !strings.Contains(r.Body.String(), want) {
					t.Errorf("missing %q", want)
				}
			}
		})
	}
}
