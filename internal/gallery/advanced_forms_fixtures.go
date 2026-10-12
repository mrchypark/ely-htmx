package uiexamples

import "github.com/mrchypark/ely-htmx"

func advancedFormsCases() []AdvancedCase {
	chapters, _ := ui.Catalogue()
	var cases []AdvancedCase
	for _, chapter := range chapters {
		if chapter.Slug != "forms" && chapter.Slug != "account" && chapter.Slug != "settings" && chapter.Slug != "onboarding" && chapter.Slug != "maps" {
			continue
		}
		for _, story := range chapter.Stories {
			options := []map[string]any{{"Value": "alpha", "Label": "Caller Alpha"}, {"Value": "beta", "Label": "Caller Beta"}}
			fields := []map[string]any{{"Name": "subject", "Label": "Caller subject", "Value": "Caller source", "Required": true}, {"Name": "notes", "Label": "Caller notes", "Type": "textarea", "Value": "Caller notes"}}
			props := map[string]any{
				"Title": "Caller " + story.Name, "Label": "Caller field", "Name": "caller_value", "Value": "alpha", "Text": "Caller source text", "Placeholder": "Caller placeholder", "Options": options,
				"Fields": fields, "ActionLabel": "Save caller data", "Email": "caller@example.test", "DisplayName": "Caller researcher", "PasswordLabel": "Example password", "Max": 5,
				"Items": []map[string]any{{"Value": "alpha", "Label": "Caller Alpha", "Text": "Caller first", "Key": "first", "Done": false, "Date": "2026-10-20", "Amount": "120", "Status": "Paid", "Email": "caller@example.test", "Role": "alpha", "Shortcut": "Ctrl+K", "Fields": fields}, {"Value": "beta", "Label": "Caller Beta", "Text": "Caller second", "Key": "second", "Done": true, "Fields": []map[string]any{{"Name": "method", "Label": "Caller method", "Value": "survey", "Required": true}}}},
			}
			switch story.Name {
			case "Combobox", "MultiSelect":
				props["Options"] = append(options, map[string]any{"Value": "locked", "Label": "Disabled caller choice", "Disabled": true})
			case "PinInput", "TwoFactorInput":
				props["Items"], props["Value"] = []int{1, 2, 3, 4}, "1234"
			case "TagInput":
				props["Value"] = "Caller Alpha,Caller Beta"
			case "NumberInput", "ScrubInput", "CurrencyInput", "Slider", "Knob", "Rating":
				props["Value"], props["Min"], props["Max"], props["Step"] = "3", "0", 10, "1"
			case "EmailInput":
				props["Value"] = "caller@example.test"
			case "PhoneInput":
				props["Value"], props["Mask"] = "01012345678", "###-####-####"
			case "RegexInput":
				props["Value"], props["Text"] = "Caller", "Caller source text"
			case "ExpressionInput":
				props["Value"], props["Variables"] = "price * quantity", map[string]any{"price": 8, "quantity": 3}
			case "Calendar", "DatePicker", "RelativeDatePicker", "LocationPicker":
				props["Value"], props["End"], props["BaseDate"] = "2026-10-20", "2026-10-22", "2026-10-20"
				if story.Name == "RelativeDatePicker" {
					props["Options"] = []map[string]any{{"Label": "Tomorrow", "Days": 1}, {"Label": "Next week", "Days": 7}}
				}
			case "MonthPicker":
				props["Value"] = "2026-10"
			case "CronEditor":
				props["Value"] = "*/5 * * * *"
			case "ColorPicker", "ColorSwatch":
				props["Value"], props["Options"] = "#335cff", []map[string]any{{"Label": "Blue", "Value": "#335cff"}, {"Label": "Green", "Value": "#228855"}}
			case "GradientEditor":
				props["Start"], props["End"], props["Angle"] = "#335cff", "#228855", 90
			case "CascadeSelect":
				props["Options"] = []map[string]any{{"Label": "Caller parent", "Value": "parent", "Children": options}}
			case "SignaturePad":
				props["Value"] = "Caller signature"
			case "SettingsLayout":
				props["Items"] = []map[string]any{{"Label": "Caller settings", "Fields": fields}}
			}
			if chapter.Slug == "maps" {
				points := []map[string]any{{"Label": "South Korea", "Value": "73", "Latitude": 37.56, "Longitude": 126.98}, {"Label": "Japan", "Value": "42", "Latitude": 35.67, "Longitude": 139.65}}
				props["Value"], props["Points"], props["Routes"] = "37.5665,126.9780", points, []map[string]any{{"Points": points}}
			}
			cases = append(cases, AdvancedCase{Family: "forms", Component: ui.Component{Chapter: chapter.Slug, Name: story.Name, ID: "audit-" + chapter.Slug + "-" + story.Slug, Props: props}})
		}
	}
	return cases
}
