package uiexamples

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMailSourceStories(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for slug, wants := range map[string][]string{
		"mailboxlist-folderlist":               {"Mailboxes", "Folders", "Labels", "Snoozed", "Spam", "Invoices", "<b>12</b>", "<b>4</b>"},
		"maillist-mailitem":                    {"Invoice 2026-114", "The brief, final", "A visit from the press", "Welcome to the studio", "3 attachments", "Sep 8, 2025"},
		"inboxzerostate-emptystate":            {"All caught up", "Nothing waits in your inbox.", "Open the archive"},
		"mailreader-quotedtext":                {"Morning both,", "ana@atrium.studio", "2.4 MB", "3.1 MB", "On Thursday, Ben Ito wrote:", "reply-all"},
		"mailthreadview":                       {"Plaster samples", "Leave notes on the backs.", "The lime reads warmest", "Dusk works.", "The tadelakt is sealed.", "data-ely-mail-thread-message open"},
		"mailcomposer-recipientinput-taginput": {"Cc Bcc", "Add people", "atrium-plan.pdf", "Chloé Martin", "Lead designer · Atrium Studio", "data-ely-mail-compose-schedule"},
		"signatureeditor":                      {"Under each message", "Chloé Martin", "Lead designer · Atrium Studio"},
		"labelpicker":                          {"Clients", "Studio", "Press", "Finance", "Site visits", "Travel", "Find or make a label"},
		"snoozepicker":                         {"Snooze", "Pick a time…", "Snooze until"},
		"schedulesend":                         {">Send</button>", "Pick a time…", "Schedule send"},
	} {
		t.Run(slug, func(t *testing.T) {
			r := httptest.NewRecorder()
			s.ServeHTTP(r, httptest.NewRequest("GET", "/components/mail/"+slug, nil))
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
