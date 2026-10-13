package uiexamples

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalendarSourceStories(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for slug, wants := range map[string][]string{
		"calendarmonthview":                                        {"October 2026", "Studio stand-up", "Site visit · Lisbon", "data-calendar-span=\"3\""},
		"minicalendar-forms-calendar":                              {"October 2026", "2026-10-13", "ely-calendar-mini-grid"},
		"calendarweekview-timegrid-alldayrow-currenttimeindicator": {"Oct 12–18, 2026", "Sun 18", "Studio offsite", "2026-10-13T09:30"},
		"calendardayview-timezoneoverlay":                          {"JST", "KST", "Design review"},
		"agendaview":                                               {"Today · Tuesday, October 13", "Studio offsite", "Night shoot"},
		"eventpopover-attendeelist":                                {"Chloé Martin", "eli@atrium.studio", "5 guests", "data-calendar-popover"},
		"eventeditor-recurrenceeditor-reminderpicker":              {"Model room", "Weekly", "data-calendar-reminders", "2026-10-13T11:00"},
		"availabilitypicker":                                       {"Wed 14", "Tue 20", "2026-10-14T09:00"},
		"yearview":                                                 {"January", "December", "2026-10-13"},
	} {
		t.Run(slug, func(t *testing.T) {
			r := httptest.NewRecorder()
			s.ServeHTTP(r, httptest.NewRequest("GET", "/components/calendar/"+slug, nil))
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
