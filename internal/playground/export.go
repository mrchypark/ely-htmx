package playground

import (
	"encoding/json"
	"fmt"
	ui "github.com/mrchypark/ely-htmx"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ExportStatic writes the browser-only editor. The caller supplies Go's matching
// wasm_exec.js and the compiled renderer.wasm beside these generated files.
func ExportStatic(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	write := func(name string, b []byte) error {
		target := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		return os.WriteFile(target, b, 0644)
	}
	for _, name := range []string{"app.js", "style.css", "static-provider.js", "renderer-worker.js"} {
		b, err := web.ReadFile("web/" + name)
		if err != nil {
			return err
		}
		if err = write(name, b); err != nil {
			return err
		}
	}
	index, err := web.ReadFile("web/index.html")
	if err != nil {
		return err
	}
	html := strings.NewReplacer(`content="{{.Token}}"`, `content=""`, `href="/style.css"`, `href="./style.css"`, `src="/app.js"`, `src="./app.js"`, `href="/" aria-label="Ely Studio home"`, `href="../" aria-label="Component gallery"`, `Local workspace`, `Browser-only workspace`, `Actual server-encoded sizes`, `Actual Go/WASM-encoded sizes`, `The server defines each profile`, `The shared Go codec defines each profile`).Replace(string(index))
	html = strings.Replace(html, `<script src="./app.js" defer></script>`, `<script src="./static-provider.js" defer></script><script src="./app.js" defer></script>`, 1)
	html = strings.Replace(html, `<meta name="color-scheme"`, `<meta http-equiv="Content-Security-Policy" content="default-src 'none'; script-src 'self' 'wasm-unsafe-eval'; worker-src 'self'; style-src 'self' 'unsafe-inline'; connect-src 'self'; frame-src 'self'; img-src 'none'; base-uri 'none'; form-action 'none'"><meta name="color-scheme"`, 1)
	html = strings.Replace(html, `SERVER-DRIVEN UI PLAYGROUND`, `<a href="../">COMPONENT GALLERY ↗</a>`, 1)
	html = strings.Replace(html, `Compose real Ely components. Inspect the result, down to the byte.`, `Compose real Ely components. Your JSON stays in this browser. The Go renderer loads on first preview; slower devices may take a moment.`, 1)
	if err = write("index.html", []byte(html)); err != nil {
		return err
	}
	registry, err := newRegistry()
	if err != nil {
		return err
	}
	items := []map[string]string{}
	for i, d := range registry.Definitions() {
		filename := fmt.Sprintf("contracts/%04d.json", i)
		items = append(items, map[string]string{"name": d.Name, "family": d.Family, "file": filename})
		b, e := json.Marshal(d)
		if e != nil {
			return e
		}
		if e = write(filename, b); e != nil {
			return e
		}
	}
	for name, value := range map[string]any{"catalogue.json": items, "examples.json": examples()} {
		b, e := json.Marshal(value)
		if e != nil {
			return e
		}
		if e = write(name, b); e != nil {
			return e
		}
	}
	styles := []string{}
	err = fs.WalkDir(ui.Files, "assets", func(name string, entry fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if entry.IsDir() || !strings.HasSuffix(name, ".css") {
			return nil
		}
		b, e := ui.Files.ReadFile(name)
		if e != nil {
			return e
		}
		if e = write(name, b); e != nil {
			return e
		}
		styles = append(styles, name)
		return nil
	})
	if err != nil {
		return err
	}
	b, err := json.Marshal(styles)
	if err != nil {
		return err
	}
	return write("styles.json", b)
}
