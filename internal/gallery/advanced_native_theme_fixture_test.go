package uiexamples

import "testing"

func TestNativeThemeModeFixtures(t *testing.T) {
	fixtures := map[string]map[string]any{}
	for _, f := range advancedNativeCases() {
		fixtures[f.Component.Chapter+"/"+f.Component.Name] = f.Component.Props
	}
	for name, count := range map[string]int{"Surfaces": 7, "Ink": 6, "Signals": 4, "Series": 8, "Syntax": 13} {
		p := fixtures["theme/"+name]
		items := p["Items"].([]map[string]any)
		if len(items) != count {
			t.Fatalf("%s: %d items", name, len(items))
		}
		for _, item := range items {
			if item["Color"] == nil || item["DarkColor"] == nil {
				t.Fatalf("%s missing active-mode colors", name)
			}
		}
	}
	surfaces := fixtures["theme/Surfaces"]
	if len(surfaces["PreviewExamples"].([]map[string]any)) != 1 {
		t.Fatal("Surfaces must show only the active mode")
	}
	if surfaces["Items"].([]map[string]any)[0]["DarkColor"] != "#252525" {
		t.Fatal("source dark bg changed")
	}
	if fixtures["theme/Motion"]["HideReplay"] != true {
		t.Fatal("source Motion has no action button")
	}
	editor := fixtures["theme/ThemeEditor"]
	if editor["DensitySegments"] != true {
		t.Fatal("source density segments missing")
	}
	colors := editor["Colors"].([]map[string]any)
	if colors[0]["DarkColor"] != "#252525" || colors[2]["DarkColor"] != "#3F3F3F" {
		t.Fatal("editor default dark colors differ from source")
	}
	safe := fixtures["i18n/ColorBlindSafe Palette"]["Items"].([]map[string]any)
	if safe[0]["DarkColor"] != "#79A7E2" || safe[0]["DarkSafeColor"] != "#3471E3" {
		t.Fatal("source standard/safe dark hues changed")
	}
}
