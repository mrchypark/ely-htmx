package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestCommunicationAPIRendersEveryPrimaryWithCallerData(t *testing.T) {
	tpl, err := NewTemplates()
	if err != nil {
		t.Fatal(err)
	}
	definitions, err := ComponentDefinitions()
	if err != nil {
		t.Fatal(err)
	}
	for _, definition := range definitions {
		if definition.Family != "communication" {
			continue
		}
		t.Run(definition.Chapter+"/"+definition.Name, func(t *testing.T) {
			var out bytes.Buffer
			props := map[string]any{"Title": `<img src=x onerror=1>`, "Content": "Caller content", "Value": "12", "Items": []map[string]any{{"Value": "one", "Label": "Caller item", "Content": "Caller detail"}}, "Options": []map[string]any{{"Value": "one", "Label": "One", "Selected": true}}}
			if err := RenderComponent(tpl, &out, Component{Chapter: definition.Chapter, Name: definition.Name, ID: "audit-" + definition.Chapter, Props: props}); err != nil {
				t.Fatal(err)
			}
			got := out.String()
			if strings.TrimSpace(got) == "" || strings.Contains(got, "<img src=x") || !strings.Contains(got, `id="audit-`+definition.Chapter+`"`) {
				t.Fatalf("unsafe or empty renderer: %s", got)
			}
		})
	}
}

func TestCommunicationAPIUsesCallerItemsAndDisabledState(t *testing.T) {
	tpl, err := NewTemplates()
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	err = RenderComponent(tpl, &out, Component{Chapter: "chat", Name: "ConversationList", ID: "audit-chat-disabled", Props: map[string]any{
		"Title": "Caller conversation", "Disabled": true, "Items": []map[string]any{{"Value": "one", "Label": "<caller item>", "Content": "One instance only"}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{`id="audit-chat-disabled"`, "&lt;caller item&gt;", "One instance only", "disabled"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in %s", want, got)
		}
	}
	if strings.Contains(got, "<caller item>") {
		t.Fatal("caller item was interpreted as HTML")
	}
}

func TestCommunicationAPIStructuredStoriesUseNativeControls(t *testing.T) {
	tpl, err := NewTemplates()
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		chapter, name string
		props         map[string]any
		wants         []string
	}{
		{"calendar", "CalendarMonthView", map[string]any{"Title": "October", "Items": []map[string]any{{"Date": "2026-10-12", "Day": "12", "Events": []map[string]any{{"Value": "event-1", "Label": "Review"}}}}}, []string{"ely-calendar-grid", `data-ely-date="2026-10-12"`, `data-ely-action="open-event"`}},
		{"project", "KanbanBoard", map[string]any{"Lanes": []map[string]any{{"Value": "todo", "Label": "To do", "Cards": []map[string]any{{"Value": "task-1", "Title": "Review", "Priority": "High", "Due": "2026-10-14"}}}}}, []string{`data-lane-id="todo"`, `data-card-id="task-1"`, "High", "2026-10-14"}},
		{"mail", "MailComposer", map[string]any{"To": "to@example.test", "Subject": "Review", "Value": "Body"}, []string{`type="email"`, `name="subject"`, `data-slot="textarea"`, "required"}},
		{"agent", "ToolApprovalDialog", map[string]any{"Title": "Permission", "Content": "Read document", "Payload": "doc-1"}, []string{"<dialog", `data-ely-action="approve"`, `data-ely-action="deny"`}},
	}
	for _, tc := range cases {
		var out bytes.Buffer
		if err := RenderComponent(tpl, &out, Component{Chapter: tc.chapter, Name: tc.name, ID: "structured-" + tc.chapter, Props: tc.props}); err != nil {
			t.Fatal(err)
		}
		for _, want := range tc.wants {
			if !strings.Contains(out.String(), want) {
				t.Errorf("%s/%s missing %q", tc.chapter, tc.name, want)
			}
		}
	}
}
