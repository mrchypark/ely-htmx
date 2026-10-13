package uiexamples

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProjectSourceStories(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for path, wants := range map[string][]string{
		"calendar/eventcard-eventchip":                                       {"Design review", "11:00", "Model room", "Night shoot", "Sun 21:00", "Studio offsite", "Oct 15 – 17"},
		"project/priorityindicator-issueidbadge":                             {"Urgent", "High", "Medium", "Low", "No priority", "ARC-41", "ENG-1204"},
		"project/statusselect-assigneepicker-duedatepicker":                  {"In progress", "Ana Lima", "2026-10-12", "name=\"assignee\""},
		"project/taskitem-tasklist":                                          {"Write the client brief", "Order plaster samples", "Photograph the site model", "Revise the lighting budget", "Book the land survey", "Archive last year&#39;s renders", "ARC-17"},
		"project/subtasklist":                                                {"2 of 4", "Measure the room", "Sketch two options", "Price the fixtures", "Send to Chloé", "data-project-subtask-progress"},
		"project/kanbanboard-kanbancolumn-kanbancard-issuecard-sprint-board": {"data-lane-id=\"todo\"", "data-lane-id=\"doing\"", "data-card-id=\"brief\"", "ARC-44", "Order plaster samples"},
		"project/taskdetailpanel":                                            {"ARC-41", "Write the client brief", "The client wants the brief before the site visit", "Draft the scope", "Share for comments", "1 of 3"},
		"project/labelmanager":                                               {"Client", "Site", "Render", "Budget", "New label"},
		"project/milestoneprogress":                                          {"Concept", "14 of 14 issues", "Schematic design", "9 of 20 issues", "Permit set", "At risk", "Land survey", "2 days late"},
		"project/workloadview":                                               {"Chloé Martin", "12 / 10 pts", "Ana Lima", "7 / 10 pts", "9 / 8 pts", "Dev Rao"},
		"project/timetracker":                                                {"What are you working on?", "Client brief", "1h 35m", "Site photos", "40m", "Budget review", "25m", "2h 40m"},
		"project/pomodoro-timer":                                             {"25:00", "Focus", "data-project-rounds", "data-project-skip"},
		"project/databaseview":                                               {"data-project-layout", "Table", "Board", "Gallery", "Calendar", "Timeline", "Order plaster samples", "In review"},
	} {
		t.Run(path, func(t *testing.T) {
			r := httptest.NewRecorder()
			s.ServeHTTP(r, httptest.NewRequest("GET", "/components/"+path, nil))
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
