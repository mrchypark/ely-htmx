package uiexamples

import "strings"

// Only gallery-owned rows opt into mode changes. Plain caller Color values stay fixed.
func nativeThemeModeFixtures(chapter, name string, props map[string]any) {
	if chapter != "theme" && !(chapter == "i18n" && name == "ColorBlindSafe Palette") {
		return
	}
	items, _ := props["Items"].([]map[string]any)
	keys := map[string][]string{
		"Surfaces":               strings.Fields("bg surface sunken hover active border border_strong"),
		"Ink":                    strings.Fields("fg fg_muted fg_subtle fg_disabled focus link"),
		"Signals":                strings.Fields("success warning danger info"),
		"Series":                 strings.Fields("chart_0 chart_1 chart_2 chart_3 chart_4 chart_5 chart_6 chart_7"),
		"Syntax":                 strings.Fields("syntax_keyword syntax_function syntax_punctuation syntax_operator syntax_type_name syntax_variable syntax_property syntax_constant syntax_number syntax_comment syntax_tag syntax_attribute syntax_string"),
		"ColorBlindSafe Palette": strings.Fields("chart_0 chart_1 chart_2 chart_3 chart_4 chart_5 chart_6 chart_7"),
	}
	if name == "Ink" {
		filtered := make([]map[string]any, 0, 6)
		for _, item := range items {
			if item["Label"] != "accent" {
				filtered = append(filtered, item)
			}
		}
		items = filtered
		props["Items"] = items
	}
	for i, key := range keys[name] {
		if i >= len(items) {
			break
		}
		pair := nativeToolingPalette[key]
		items[i]["Color"], items[i]["DarkColor"] = strings.ToUpper(pair[0]), strings.ToUpper(pair[1])
		if name == "Signals" {
			quiet := nativeToolingPalette[key+"_subtle"]
			items[i]["Content"], items[i]["DarkContent"] = strings.ToUpper(quiet[0]), strings.ToUpper(quiet[1])
		}
	}
	if name == "ColorBlindSafe Palette" {
		for i, value := range []string{"#3471E3", "#57B6B2", "#F6A537", "#DE545A", "#9D8DEC", "#C4ECB7", "#FFD93D", "#52DAFE"} {
			items[i]["DarkSafeColor"] = value
		}
	}
	if name == "Surfaces" {
		props["PreviewExamples"] = []map[string]any{{"Title": ""}}
	}
	if name == "Motion" {
		props["HideReplay"] = true
	}
	if name == "ThemeEditor" || name == "ThemeProvider" {
		props["DensitySegments"] = true
		colors, _ := props["Colors"].([]map[string]any)
		var tokens strings.Builder
		for _, color := range colors {
			key := strings.ReplaceAll(strings.TrimPrefix(color["Token"].(string), "--"), "-", "_")
			pair, ok := nativeToolingPalette[key]
			if !ok {
				continue
			}
			color["Color"], color["DarkColor"] = strings.ToUpper(pair[0]), strings.ToUpper(pair[1])
			tokens.WriteString(color["Token"].(string) + ": " + color["Color"].(string) + ";\n")
		}
		props["Tokens"] = tokens.String()
	}
}
