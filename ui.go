// Package ui provides shared Go templates, theme styles, and browser assets.
package ui

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

//go:embed templates/*.html assets/* components.json LICENSE THIRD_PARTY_UI.md
var Files embed.FS

// Option supplies a submitted value, a label, and optional choice icons.
type Option struct {
	Value    string
	Label    string
	Selected bool
	Disabled bool
	Icon     string
	Glyph    string
	Note     string
	Keywords string
	Group    string
}

var extraAttributeName = regexp.MustCompile(`^(data|aria|hx)-[a-z0-9]+(-[a-z0-9]+)*$`)

func attrName(name string) (template.HTMLAttr, error) {
	if !extraAttributeName.MatchString(name) {
		return "", fmt.Errorf("unsupported component attribute %q", name)
	}
	// Only the validated name bypasses escaping. Attribute values remain escaped.
	return template.HTMLAttr(name), nil
}

func NewTemplates() (*template.Template, error) {
	tpl := template.New("").Funcs(Funcs())
	tpl.Funcs(template.FuncMap{"component": func(value any) (template.HTML, error) {
		c, err := componentFrom(value)
		if err != nil {
			return "", err
		}
		var output bytes.Buffer
		if err := RenderComponent(tpl, &output, c); err != nil {
			return "", err
		}
		return template.HTML(output.String()), nil
	}})
	return tpl.ParseFS(Files, "templates/*.html")
}

// BindAssets registers shared asset paths without taking over application assets.
func BindAssets(mux *http.ServeMux) error {
	assets, err := fs.Sub(Files, "assets")
	if err != nil {
		return err
	}
	handler := http.StripPrefix("/assets/", http.FileServer(http.FS(assets)))
	return fs.WalkDir(assets, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(assets, name)
		if err != nil {
			return err
		}
		etag := fmt.Sprintf(`"%x"`, sha256.Sum256(data))
		mux.Handle("GET /assets/"+name, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, no-cache")
			w.Header().Set("ETag", etag)
			handler.ServeHTTP(w, r)
		}))
		return nil
	})
}

func Funcs() template.FuncMap {
	return template.FuncMap{
		"qrSVG": qrSVG, "barcodeSVG": barcodeSVG, "markdownHTML": markdownHTML,
		"withProps": func(source map[string]any, values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("withProps requires key/value pairs")
			}
			out := make(map[string]any, len(source)+len(values)/2)
			for key, value := range source {
				out[key] = value
			}
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("withProps key must be string")
				}
				out[key] = values[i+1]
			}
			return out, nil
		},
		"seq": func(count int) ([]int, error) {
			if count < 0 || count > 1000 {
				return nil, fmt.Errorf("sequence length must be between 0 and 1000")
			}
			items := make([]int, count)
			for i := range items {
				items[i] = i + 1
			}
			return items, nil
		},
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict requires key/value pairs")
			}
			result := map[string]any{}
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict key must be string")
				}
				result[key] = values[i+1]
			}
			return result, nil
		},
		"list":      func(v ...any) []any { return v },
		"add":       func(a, b int) int { return a + b },
		"join":      strings.Join,
		"attrName":  attrName,
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"deref": func(v any) any {
			r := reflect.ValueOf(v)
			if r.IsValid() && r.Kind() == reflect.Pointer {
				if r.IsNil() {
					return nil
				}
				return r.Elem().Interface()
			}
			return v
		},
	}
}
