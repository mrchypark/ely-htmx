package main

import (
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseRenderedLinks(t *testing.T) {
	source := "# Hello, **world**!\n\n# Hello, world!\n\n[reference][target]\n\n[target]: <file%20name.md#intro>\n\n![image](image.png)\n\n<https://example.org/>\n\n<img\n src=\"mark.svg\" alt=\"mark\">\n\n<a id=\"explicit\"></a>\n\n`[code](missing.md)`\n\n```html\n<a href=\"missing.html\">sample</a>\n```\n"
	got := parse("README.md", []byte(source))
	want := []string{"file%20name.md#intro", "image.png", "https://example.org/", "mark.svg"}
	if !reflect.DeepEqual(got.links, want) {
		t.Fatalf("links = %#v, want %#v", got.links, want)
	}
	for _, id := range []string{"hello-world", "hello-world-1", "explicit"} {
		if !got.ids[id] {
			t.Errorf("missing %q in %#v", id, got.ids)
		}
	}
}
func TestLocalTargets(t *testing.T) {
	for _, tc := range []struct {
		source, raw, pages, want string
		local                    bool
	}{
		{"docs/a.md", "../README.md#start", "", "README.md", true},
		{"docs/a.md", "file%20name.md", "", "docs/file name.md", true},
		{"README.md", "https://github.com/mrchypark/elysdui/blob/main/native/README.md", "", "native/README.md", true},
		{"README.md", "https://mrchypark.github.io/elysdui/go-htmx/playground/", "dist", "dist/go-htmx/playground", true},
		{"site/index.html", "go-htmx/", "dist", "dist/go-htmx", true},
		{"site/index.html", "/elysdui/go-htmx/", "dist", "dist/go-htmx", true},
		{"docs/a.md", "/README.md", "", "", false},
		{"site/index.html", "#renderers", "dist", "site/index.html", true},
		{"README.md", "https://example.org/a", "", "", false},
	} {
		u, err := url.Parse(tc.raw)
		if err != nil {
			t.Fatal(err)
		}
		got, local := localTarget(tc.source, u, tc.pages)
		if got != tc.want || local != tc.local {
			t.Errorf("%s: got %q, %v; want %q, %v", tc.raw, got, local, tc.want, tc.local)
		}
	}
}
func TestMissingLinksAndAnchors(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "README.md")
	if err := os.WriteFile(p, []byte("# Intro\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := checkFile(p, "intro"); err != nil {
		t.Fatal(err)
	}
	if err := checkFile(dir, "intro"); err != nil {
		t.Fatal(err)
	}
	if err := checkFile(p, "missing"); err == nil {
		t.Error("accepted missing anchor")
	}
	if err := checkFile(filepath.Join(dir, "missing.md"), ""); err == nil {
		t.Error("accepted missing file")
	}
}

func TestGitHubSourceLineAnchors(t *testing.T) {
	for _, fragment := range []string{"L1", "L20-L30"} {
		if !sourceLine.MatchString(fragment) {
			t.Errorf("rejected %s", fragment)
		}
	}
	for _, fragment := range []string{"intro", "L0", "L20-extra"} {
		if sourceLine.MatchString(fragment) {
			t.Errorf("accepted %s", fragment)
		}
	}
}
