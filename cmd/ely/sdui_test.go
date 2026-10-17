package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSDUIValidateFileAndStdin(t *testing.T) {
	screen := `{"version":"1","id":"screen","loader":"hostLoader","root":{"id":"root","type":"region","children":[{"id":"save","type":"button","props":{"label":"Save"},"event":{"on":"click","action":"hostAction","target":"root"}}]}}`
	for _, prefix := range [][]string{{"validate"}, {"sdui", "validate"}} {
		var stdout, stderr bytes.Buffer
		args := append(append([]string(nil), prefix...), "-registry", "core", "-")
		if code := runSDUI(args, strings.NewReader(screen), &stdout, &stderr); code != 0 {
			t.Fatalf("validate code=%d: %s", code, stderr.String())
		}
		if !strings.Contains(stdout.String(), "Valid screen screen") || !strings.Contains(stdout.String(), "host loader/action registration, authorization and runtime data were not checked") {
			t.Fatal(stdout.String())
		}
	}
	path := filepath.Join(t.TempDir(), "screen.json")
	if err := os.WriteFile(path, []byte(screen), 0600); err != nil {
		t.Fatal(err)
	}
	if code := runSDUI([]string{"validate", "-registry", "core", path}, strings.NewReader("invalid"), io.Discard, io.Discard); code != 0 {
		t.Fatal("file validation failed")
	}
	for _, input := range []string{`null`, `{"version":"1","version":"1"}`, `{"version":"1","id":"screen","root":{"id":"root","type":"text","props":{"text":3}}}`, `{"version":"1","id":"screen","root":{"id":"root","type":"unknown"}}`} {
		var stdout, stderr bytes.Buffer
		if code := runSDUI([]string{"validate", "-registry", "core", "-"}, strings.NewReader(input), &stdout, &stderr); code != 1 || stdout.Len() != 0 || stderr.Len() == 0 {
			t.Fatalf("invalid validation code %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
	}
}

func TestSDUISchemaAndCatalog(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := runSDUI([]string{"schema", "-registry", "core"}, nil, &stdout, &stderr); code != 0 {
		t.Fatalf("schema: %d %s", code, stderr.String())
	}
	var schema map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &schema); err != nil {
		t.Fatal(err)
	}
	if schema["$schema"] != "https://json-schema.org/draft/2020-12/schema" {
		t.Fatal("wrong schema dialect")
	}
	stdout.Reset()
	stderr.Reset()
	if code := runSDUI([]string{"catalog", "-registry", "core"}, nil, &stdout, &stderr); code != 0 {
		t.Fatalf("catalog: %d %s", code, stderr.String())
	}
	var catalog map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &catalog); err != nil {
		t.Fatal(err)
	}
	if len(catalog["definitions"].([]any)) != 10 {
		t.Fatal("wrong core catalogue")
	}
	stdout.Reset()
	stderr.Reset()
	if code := runSDUI([]string{"catalog", "-registry", "core", "-type", "input"}, nil, &stdout, &stderr); code != 0 {
		t.Fatalf("filtered catalog: %d %s", code, stderr.String())
	}
	var definition map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &definition); err != nil {
		t.Fatal(err)
	}
	if definition["name"] != "input" || definition["properties"] == nil {
		t.Fatal("filtered contract missing")
	}
	if code := runSDUI([]string{"catalog", "-registry", "core", "-type", "unknown"}, nil, io.Discard, io.Discard); code != 1 {
		t.Fatal("unknown type succeeded")
	}
}

func TestSDUICommandUsageAndOutputFailure(t *testing.T) {
	for _, args := range [][]string{{}, {"sdui"}, {"unknown"}, {"schema", "extra"}, {"validate", "-registry", "core"}, {"validate", "-unknown"}} {
		if code := runSDUI(args, nil, io.Discard, io.Discard); code != 2 {
			t.Errorf("usage %v returned %d", args, code)
		}
	}
	if code := runSDUI([]string{"schema", "-registry", "unknown"}, nil, io.Discard, io.Discard); code != 1 {
		t.Fatal("unknown registry succeeded")
	}
	if code := runSDUI([]string{"validate", "-registry", "core", "/definitely-missing/screen.json"}, nil, io.Discard, io.Discard); code != 1 {
		t.Fatal("missing file succeeded")
	}
	if code := runSDUI([]string{"schema", "-help"}, nil, io.Discard, io.Discard); code != 0 {
		t.Fatal("help failed")
	}
	if code := runSDUI([]string{"schema", "-registry", "core"}, nil, errorOutput{}, io.Discard); code != 1 {
		t.Fatal("writer failure ignored")
	}
}

type errorOutput struct{}

func (errorOutput) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }
