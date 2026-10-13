package uiexamples

import (
	_ "embed"
	"encoding/json"
	"math"
)

// Source fixtures from Ely 49ed1f8 examples/gallery/pages/{editor,terminal,git,debug}.
// Keep these separate from the caller-owned audit fixtures, which exercise API isolation.
//
//go:embed assets/workspace-source.json
var workspaceSourceJSON []byte

func workspaceSourcePreviews() map[string]map[string]any {
	var scenes map[string]map[string]any
	if err := json.Unmarshal(workspaceSourceJSON, &scenes); err != nil {
		panic(err)
	}
	for _, scene := range scenes {
		workspaceSourceNumbers(scene)
	}
	return scenes
}

// Template arithmetic helpers take ints; retain fractional timeline timestamps.
func workspaceSourceNumbers(value any) any {
	switch v := value.(type) {
	case map[string]any:
		for key, item := range v {
			v[key] = workspaceSourceNumbers(item)
		}
	case []any:
		for i, item := range v {
			v[i] = workspaceSourceNumbers(item)
		}
	case float64:
		if math.Trunc(v) == v {
			return int(v)
		}
	}
	return value
}
