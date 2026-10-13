package uiexamples

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCommunicationCompletionStories(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for slug, hooks := range map[string][]string{
		"chat/citationbadge-inlinecitation-sourcecard-sourcelist-sourcespanel": {"ely-chat-citation-answer", "3 sources", "Gamma and the eye"},
		"chat/searchprogress-webresultcard-documentchunkpreview":               {"Designing dark themes", "lift-notes.pdf", "<mark>more lift</mark>"},
		"chat/modelselector-modeselector-tooltogglemenu-prompttemplatemenu":    {"ely-chat-model-modes", "<summary>Tools</summary>", "<summary>Prompts</summary>"},
		"chat/streamingtext-streamingmarkdown-streamingcursor":                 {"data-ely-stream-status", "replay-stream"},
		"agent/artifactpanel-artifactversionswitcher-livepreview":              {"data-ely-artifact-source", "data-ely-artifact-version"},
		"agent/sandboxstatus": {"data-sandbox-state", "<meter"},
		"agent/mcpserverlist-connectorlist-toolregistry-panel": {"configure-server", "Read only", "Open setup", "data-ely-tool-access"},
		"agent/memorypanel-checkpointlist-rewind":              {"data-ely-memory-count", "data-ely-checkpoint"},
		"agent/interruptbutton-humaninputrequest":              {"data-ely-interrupt", "data-ely-answer-choice"},
	} {
		t.Run(slug, func(t *testing.T) {
			response := httptest.NewRecorder()
			server.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/components/"+slug, nil))
			if response.Code != http.StatusOK {
				t.Fatalf("status %d", response.Code)
			}
			for _, hook := range hooks {
				if !strings.Contains(response.Body.String(), hook) {
					t.Errorf("missing %q", hook)
				}
			}
		})
	}
}
