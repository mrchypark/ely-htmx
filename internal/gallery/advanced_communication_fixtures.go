package uiexamples

import "github.com/mrchypark/ely-htmx"

// advancedCommunicationCases gives each communication story its own caller-owned data.
func advancedCommunicationCases() []AdvancedCase {
	definitions, err := ui.ComponentDefinitions()
	if err != nil {
		return nil
	}
	indexes := map[string]int{}
	var cases []AdvancedCase
	for _, definition := range definitions {
		if definition.Family != "communication" {
			continue
		}
		indexes[definition.Chapter]++
		id := "audit-" + definition.Chapter + "-" + itoa(indexes[definition.Chapter])
		cases = append(cases, AdvancedCase{Family: "communication", Component: ui.Component{Chapter: definition.Chapter, Name: definition.Name, ID: id, Props: map[string]any{
			"Title": definition.Name + " for Northwind review", "Label": "Caller label", "Content": "Caller-controlled communication content for " + id, "Value": "25", "Name": "message", "Placeholder": "Write a caller-owned message", "Required": indexes[definition.Chapter]%4 == 0, "Disabled": indexes[definition.Chapter]%13 == 0 && definition.Name != "TemperatureSlider" && definition.Name != "Pomodoro Timer",
			"Author": "Rina Chen", "Time": "09:30", "Count": 2, "Reaction": "♥", "Status": "Awaiting host operation", "Progress": 45, "Max": 100, "Start": "2026-10-12T09:00", "End": "2026-10-12T10:00", "From": "rina@example.test", "To": "team@example.test", "Subject": "Northwind review", "Quote": "Use the recorded result.", "Payload": `{"task":"review"}`, "Path": "notes/review.md", "Diff": "- draft\n+ reviewed", "Alt": "Caller supplied preview", "URL": "",
			"Items":   []map[string]any{{"Value": id + "-first", "Label": "First caller item", "Title": "First caller item", "Content": "Independent first value"}, {"Value": id + "-second", "Label": "Second caller item", "Title": "Second caller item", "Content": "Independent second value"}},
			"Options": []map[string]any{{"Value": "open", "Label": "Open", "Selected": true}, {"Value": "closed", "Label": "Closed"}},
			"Context": []map[string]any{{"Value": "project:northwind", "Label": "Northwind"}, {"Value": "source:42", "Label": "Source 42"}}, "Commands": []map[string]any{{"Value": "/summarize", "Label": "/summarize"}},
			"Lanes": []map[string]any{{"Value": "todo", "Label": "To do", "Count": 1, "Cards": []map[string]any{{"Value": id + "-card", "Title": "Caller task", "Content": "Move only after host confirmation", "Priority": "High", "Due": "2026-10-14"}}}, {"Value": "done", "Label": "Done", "Count": 0}},
			"Due":   "2026-10-14", "Priority": "High", "Date": "2026-10-12", "Day": "12", "Until": "2026-10-12T13:00", "SendAt": "2026-10-12T14:00", "Events": []map[string]any{{"Value": id + "-event", "Label": "Caller event", "Start": "09:00", "End": "10:00"}},
		}}})
	}
	return cases
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	var digits [20]byte
	i := len(digits)
	for value > 0 {
		i--
		digits[i] = byte('0' + value%10)
		value /= 10
	}
	return string(digits[i:])
}
