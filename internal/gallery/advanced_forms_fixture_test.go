package uiexamples

import "testing"

func TestGalleryEmojiCataloguePreservesUpstreamData(t *testing.T) {
	if len(galleryEmojiOptions) != 1914 {
		t.Fatalf("emoji count: %d", len(galleryEmojiOptions))
	}
	seen := map[string]bool{}
	groups := map[string]bool{}
	joined := false
	for _, option := range galleryEmojiOptions {
		value := option["Value"].(string)
		if value == "" || seen[value] || option["Glyph"] != value || option["Label"] == "" {
			t.Fatalf("invalid emoji: %v", option)
		}
		seen[value] = true
		groups[option["Group"].(string)] = true
		if value == "👩‍💻" {
			joined = true
		}
	}
	if len(groups) != 9 || !joined || !seen["✨"] {
		t.Fatal("groups or complete joined emoji lost")
	}
}

func TestRegionalCataloguesPreserveUpstreamChoices(t *testing.T) {
	for key, count := range map[string]int{"languages": 184, "countries": 249, "currencies": 154} {
		options := galleryRegionalOptions[key]
		if len(options) != count {
			t.Fatalf("%s count: %d", key, len(options))
		}
		seen := map[string]bool{}
		for _, option := range options {
			v := option["Value"].(string)
			if v == "" || seen[v] || option["Label"] == "" {
				t.Fatalf("invalid %s choice: %v", key, option)
			}
			seen[v] = true
		}
	}
}
