package uiexamples

import (
	_ "embed"
	"encoding/json"
)

// Source: gallery misc pages at 49ed1f84056fa13347d68801b479a4a5a8e03b53.
//
//go:embed assets/native-misc-fixtures.json
var nativeMiscJSON []byte

func nativeMiscFixtures(name string, props map[string]any) {
	var fixtures map[string]map[string]any
	if err := json.Unmarshal(nativeMiscJSON, &fixtures); err != nil {
		panic(err)
	}
	fixture, ok := fixtures[name]
	if !ok {
		return
	}
	// Gallery-only data: reusable rendering never applies these defaults to callers.
	clear(props)
	for key, value := range fixture {
		props[key] = value
	}
	// PreviewExamples is consumed by Go before it reaches HTML templates.
	if raw, ok := props["PreviewExamples"].([]any); ok {
		previews := make([]map[string]any, 0, len(raw))
		for _, value := range raw {
			previews = append(previews, value.(map[string]any))
		}
		props["PreviewExamples"] = previews
	}
}
