package playground

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportStatic(t *testing.T) {
	dir := t.TempDir()
	if err := ExportStatic(dir); err != nil {
		t.Fatal(err)
	}
	read := func(name string) []byte {
		b, e := os.ReadFile(filepath.Join(dir, name))
		if e != nil {
			t.Fatal(e)
		}
		return b
	}
	html := string(read("index.html"))
	for _, want := range []string{`src="./static-provider.js"`, `src="./app.js"`, `href="../"`, `Browser-only workspace`, `'wasm-unsafe-eval'`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	for _, bad := range []string{`{{.Token}}`, `src="/app.js"`, `href="/style.css"`} {
		if strings.Contains(html, bad) {
			t.Errorf("found %s", bad)
		}
	}
	var items []map[string]string
	if err := json.Unmarshal(read("catalogue.json"), &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 1372 {
		t.Fatalf("catalogue count %d", len(items))
	}
	for _, item := range items {
		var contract map[string]any
		if err := json.Unmarshal(read(item["file"]), &contract); err != nil {
			t.Fatal(err)
		}
		if contract["name"] != item["name"] {
			t.Fatalf("contract mismatch %v", item)
		}
	}
	var styles []string
	if err := json.Unmarshal(read("styles.json"), &styles); err != nil {
		t.Fatal(err)
	}
	if len(styles) == 0 {
		t.Fatal("missing stylesheet manifest")
	}
	for _, name := range styles {
		read(name)
	}
	var examples []map[string]any
	if err := json.Unmarshal(read("examples.json"), &examples); err != nil {
		t.Fatal(err)
	}
	if len(examples) != 5 {
		t.Fatal("missing examples")
	}
}
