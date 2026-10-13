package uiexamples

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTraceSourceStory(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRecorder()
	s.ServeHTTP(r, httptest.NewRequest("GET", "/components/agent/traceviewer-evalresulttable", nil))
	if r.Code != 200 {
		t.Fatalf("status %d", r.Code)
	}
	body := r.Body.String()
	for _, want := range []string{`data-ely-trace-duration="4600"`, `data-trace-label="edit" aria-pressed="true"`, `data-trace-label="run_tests" aria-pressed="false"`, `class="ely-trace-controls" hidden`, `data-trace-start="950" data-trace-end="1100"`, "Fix a failing test", "Rename across files", "Add a settings page", "Explain a stack trace", "Migrate a schema", "96%", "91%", "48%", "88%", "62%", "18,400", "9,800", "31,200", "4,100", "22,700", "42.1 s", "18.6 s", "96.4 s", "7.9 s", "71.3 s"} {
		if !strings.Contains(body, want) {
			t.Errorf("source story missing %q", want)
		}
	}
	if strings.Contains(body, "<h2>Agent trace</h2>") || strings.Contains(body, "<h2>Evaluation results</h2>") {
		t.Error("source story has extra inner headings")
	}
}
