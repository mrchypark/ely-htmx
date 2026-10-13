package uiexamples

import "testing"

func TestInteractionMenuSourceDatasets(t *testing.T) {
	want := map[string]int{"Menu": 11, "CheckboxMenuItem": 8, "DropdownMenu": 2, "ContextMenu": 6, "MenuBar Menu": 2, "MegaMenu": 2, "SearchableMenu": 8, "PieMenu": 6}
	for _, fixture := range advancedInteractionCases() {
		c := fixture.Component
		if c.Chapter != "menus" {
			continue
		}
		if c.Props["Disabled"] != false {
			t.Errorf("%s must be enabled", c.Name)
		}
		if items := c.Props["Items"].([]map[string]any); len(items) != want[c.Name] {
			t.Errorf("%s: got %d rows, want %d", c.Name, len(items), want[c.Name])
		}
		if c.Name == "DropdownMenu" && len(c.Props["PreviewExamples"].([]map[string]any)) != 3 {
			t.Error("missing dropdown/split/overflow examples")
		}
		delete(want, c.Name)
	}
	if len(want) != 0 {
		t.Fatalf("missing menu datasets: %v", want)
	}
}

func TestInteractionActionFixturesRemainEnabled(t *testing.T) {
	for _, fixture := range advancedInteractionCases() {
		c := fixture.Component
		switch c.Chapter + "/" + c.Name {
		case "buttons/Button", "buttons/ActionBar", "feedback/Toast", "motion/AnimatePresence":
			if c.Props["Disabled"] != false {
				t.Errorf("%s/%s action fixture is disabled", c.Chapter, c.Name)
			}
		}
	}
}

func TestInteractionNavigationSourceDatasets(t *testing.T) {
	want := map[string]int{"NavigationMenu": 2, "BackForwardNavigation": 4, "Anchor": 5, "CommandPalette": 14, "QuickOpen": 24, "QuickSwitcher": 6, "SearchPalette": 12}
	for _, fixture := range advancedInteractionCases() {
		c := fixture.Component
		count, ok := want[c.Name]
		if c.Chapter != "navigation" || !ok {
			continue
		}
		items := c.Props["Items"].([]map[string]any)
		if len(items) != count {
			t.Errorf("%s: got %d items, want %d", c.Name, len(items), count)
		}
		if c.Name == "NavigationMenu" {
			for _, entry := range items {
				if len(entry["Items"].([]map[string]any)) != 4 {
					t.Error("navigation entry must contain four links")
				}
			}
		}
		if c.Name == "SearchPalette" {
			suggested := 0
			for _, item := range items {
				if item["Suggested"] == true {
					suggested++
				}
			}
			if suggested != 4 {
				t.Errorf("got %d search suggestions", suggested)
			}
		}
		delete(want, c.Name)
	}
	if len(want) != 0 {
		t.Fatalf("missing source datasets: %v", want)
	}
}
