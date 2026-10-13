package uiexamples

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMessagingExtendedSourceStories(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for slug, wants := range map[string][]string{
		"readreceipt":                {"Sending", "Sent", "Delivered", "Seen by", "<small>+1</small>"},
		"memberlist-userprofilecard": {"Chloé Martin", "Lead designer", "Drawing the stair", "Europe/Paris", "Farah Haddad", "Owner"},
		"onlinestatus":               {`data-presence="online"`, `data-presence="away"`, `data-presence="busy"`, `data-presence="offline"`},
		"statussetter":               {"No status set.", "In a meeting", "Commuting", "Out sick", "On leave", "30 minutes", "This week"},
		"callcontrols-voicecallbar-huddleindicator": {"Huddle in # design", "00:00", "data-ely-call-join", `data-ely-call-toggle="camera"`, `data-ely-call-toggle="share"`},
		"participanttile-videocallgrid":             {"Ana Lima", "Ben Ito", "Chloé Martin", "Dev Rao", "atrium-stair.jpg", "atrium-olive.jpg"},
		"screenshareview":                           {"Ana Lima is presenting", "screen.jpg", "Ben Ito", "Chloé Martin", `data-ely-ratio="1.6"`},
		"incomingcalldialog":                        {">Ring</button>", "Ana Lima", "Incoming video call", "Decline", "Accept"},
		"stickerpicker":                             {"Atrium", "Garden", "arch.png", "sprout.png", "Pick a sticker"},
		"gifpicker":                                 {"Light on the stair", "The olive tree", "Dunes at dusk", "The atrium, before and after", "Searching GIFs"},
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
	for _, name := range []string{"arch.png", "stair.png", "column.png", "window.png", "olive.png", "sun.png", "stone.png", "sprout.png", "stair-light.gif", "olive-breath.gif", "dunes-drift.gif", "atrium-fade.gif", "atrium-stair.jpg", "atrium-olive.jpg", "screen.jpg"} {
		t.Run(name, func(t *testing.T) {
			r := httptest.NewRecorder()
			s.ServeHTTP(r, httptest.NewRequest("GET", "/assets/communication/"+name, nil))
			if r.Code != 200 || !strings.HasPrefix(r.Header().Get("Content-Type"), "image/") || r.Body.Len() == 0 {
				t.Fatalf("media unavailable: status=%d type=%s bytes=%d", r.Code, r.Header().Get("Content-Type"), r.Body.Len())
			}
		})
	}
}
