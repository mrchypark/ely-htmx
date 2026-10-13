package uiexamples

import "testing"

func TestInteractionOverlaySourceFixtures(t *testing.T) {
	count := 0
	for _, fixture := range advancedInteractionCases() {
		c := fixture.Component
		if c.Chapter != "overlays" {
			continue
		}
		count++
		if c.Props["Disabled"] != false {
			t.Errorf("%s is disabled", c.Name)
		}
		switch c.Name {
		case "Dialog":
			items := c.Props["Items"].([]map[string]any)
			if len(items) != 5 {
				t.Errorf("got %d dialogs, want five source examples", len(items))
			}
			if items[3]["Disallow"] != "/" || items[3]["TrimRequired"] != true {
				t.Error("source prompt validation missing")
			}
		case "Lightbox":
			for _, item := range c.Props["Items"].([]map[string]any) {
				if item["Src"] != "/assets/ely-api-interaction.dunes.jpg" && item["Src"] != "/assets/ely-api-interaction.atrium.jpg" {
					t.Errorf("not an original source photo: %v", item["Src"])
				}
			}
		case "Spotlight":
			if len(c.Props["Steps"].([]map[string]any)) != 3 || c.Props["CoachIndex"] != 1 {
				t.Error("tour or coachmark source data missing")
			}
		case "FloatingToolbar":
			if c.Props["Selection"] != true {
				t.Error("toolbar must follow text selection")
			}
		case "Callout":
			if c.Props["Toggle"] != true {
				t.Error("Peek must toggle in the document flow")
			}
		}
	}
	if count != 8 {
		t.Fatalf("got %d overlay fixtures", count)
	}
}
