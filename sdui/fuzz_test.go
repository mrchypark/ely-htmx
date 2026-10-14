package sdui

import (
	"bytes"
	"testing"
)

func FuzzDecodeCompileRender(f *testing.F) {
	for _, seed := range []string{
		`{"version":"1","id":"screen","root":{"id":"message","type":"text","props":{"text":"Hello <script>"}}}`,
		`{"version":"1","id":"screen","root":{"id":"region","type":"region","children":[]}}`,
		`{"version":"1","Version":"2"}`,
		`{"root":{"props":{"value":1e10000}}}`,
		`null`, `[]`, `{"root":`,
	} {
		f.Add([]byte(seed))
	}
	registry, err := DefaultRegistry()
	if err != nil {
		f.Fatal(err)
	}
	f.Fuzz(func(t *testing.T, input []byte) {
		screen, err := Decode(bytes.NewReader(input))
		if err != nil {
			return
		}
		plan, err := Compile(screen, registry)
		if err != nil {
			return
		}
		var output bytes.Buffer
		if err := plan.Render(&output, RenderContext{Data: map[string]any{}, Endpoint: "/action", CSRF: "test"}); err != nil && output.Len() != 0 {
			t.Fatal("failed rendering wrote partial output")
		}
		if output.Len() > MaxRenderBytes {
			t.Fatal("render exceeded budget")
		}
	})
}
