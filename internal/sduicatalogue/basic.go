package sduicatalogue

import "strings"

// BasicNames is the deliberately reviewed public template list in COMPONENTS.md.
// Assets, branding hooks, and private implementation helpers are not components.
func BasicNames() []string {
	return strings.Fields(`input textarea select radio-group checkbox switch slider button field-label field-description datalist badge progress separator kbd kbd-group clearable-input password-input pin-input tag-input key-value-input file-dropzone date-picker color-picker hotkey-input signature-pad combobox multi-select selectable-list dialog tabs card stack grid code-block avatar alert empty-state accordion breadcrumb timeline menu table data-table chart calendar map-view coordinate-display`)
}
func BasicFixtures() []Fixture {
	options := []map[string]any{{"Value": "alpha", "Label": "Caller Alpha", "Selected": true}, {"Value": "beta", "Label": "Caller Beta", "Disabled": true}}
	out := []Fixture{}
	for _, name := range BasicNames() {
		p := map[string]any{"Label": "Caller " + name, "Title": "Caller title", "Description": "Caller description", "Content": "Caller content <escaped>", "Text": "Caller text", "Name": "caller_value", "Value": "alpha", "Placeholder": "Caller placeholder", "AriaLabel": "Caller control", "Disabled": false, "Required": false, "Options": options}
		switch name {
		case "input":
			p["Type"] = "text"
		case "textarea":
			p["Rows"] = 3
		case "checkbox", "switch":
			p["Checked"] = true
		case "slider", "progress":
			p["Value"], p["Min"], p["Max"] = "35", 0, 100
		case "button":
			p["Variant"], p["Size"] = "outline", "sm"
		case "field-label":
			p["For"] = "caller-field"
		case "field-description":
			p["Error"] = false
		case "badge":
			p["Variant"] = "secondary"
		case "separator":
			p["Orientation"] = "horizontal"
		case "kbd-group":
			p["Keys"] = []string{"Ctrl", "K"}
		case "pin-input":
			p["Items"], p["Value"] = []int{1, 2, 3, 4}, "1234"
		case "tag-input":
			p["Value"] = "alpha,beta"
		case "key-value-input":
			p["Items"] = []map[string]any{{"Key": "first", "Value": "Caller first"}, {"Key": "second", "Value": "Caller second"}}
		case "file-dropzone":
			p["Accept"], p["Multiple"] = ".txt", true
		case "date-picker", "calendar":
			p["Value"], p["Min"], p["Max"] = "2026-10-14", "2026-01-01", "2026-12-31"
		case "color-picker":
			p["Value"] = "#335cff"
		case "hotkey-input":
			p["Value"] = "Ctrl+K"
		case "signature-pad":
			p["Value"] = "Caller signature"
		case "selectable-list":
			p["Size"] = 4
		case "tabs":
			p["Items"] = []map[string]any{{"Value": "alpha", "Label": "Caller first tab", "Content": "Caller first panel"}, {"Value": "beta", "Label": "Caller second tab", "Content": "Caller second panel"}}
		case "stack":
			p["Horizontal"] = true
		case "avatar":
			p["Src"] = "/assets/sample.svg"
		case "alert":
			p["Urgent"] = true
		case "accordion":
			p["Items"] = []map[string]any{{"Label": "Caller first section", "Content": "Caller first body", "Open": true}, {"Label": "Caller second section", "Content": "Caller second body"}}
		case "breadcrumb", "menu":
			p["Items"] = []map[string]any{{"Value": "first", "Label": "Caller first link", "Href": "#first"}, {"Value": "second", "Label": "Caller second link", "Href": "/second", "Disabled": true}}
		case "timeline":
			p["Items"] = []map[string]any{{"Time": "09:00", "Label": "Caller first event", "Content": "Caller first detail"}, {"Time": "10:00", "Label": "Caller second event", "Content": "Caller second detail"}}
		case "table", "data-table":
			p["Columns"] = []map[string]any{{"Label": "Name"}, {"Label": "Status"}}
			p["Rows"] = []map[string]any{{"Cells": []string{"Caller Alice", "Ready"}}, {"Cells": []string{"Caller Bob", "Pending"}}}
			p["Caption"] = "Caller table"
		case "chart":
			p["Type"] = "bar"
			p["Points"] = []map[string]any{{"X": 20, "Y": 40, "Width": 16, "Height": 60, "Label": "Caller first point"}, {"X": 60, "Y": 20, "Width": 16, "Height": 80, "Label": "Caller second point"}}
		case "map-view":
			p["Mode"] = "markers"
			p["Center"] = map[string]any{"Latitude": 37.5665, "Longitude": 126.978}
			p["Points"] = []map[string]any{{"Latitude": 37.5665, "Longitude": 126.978, "Label": "Caller Seoul", "Value": "seoul"}}
			p["Routes"] = []map[string]any{{"Name": "Caller route", "Points": []map[string]any{{"Latitude": 37.5665, "Longitude": 126.978}, {"Latitude": 35.1796, "Longitude": 129.0756}}}}
		case "coordinate-display":
			p["Latitude"], p["Longitude"], p["Format"] = 37.5, 127.0, "decimal"
		}
		out = append(out, Fixture{Name: "basic/" + name, Template: name, Source: "basic:" + name, Props: p})
	}
	return out
}
