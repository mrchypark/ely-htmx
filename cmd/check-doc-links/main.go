// Command check-doc-links checks rendered documentation links without treating code examples as links.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
	"golang.org/x/net/html"
)

type document struct {
	links []string
	ids   map[string]bool
}

func htmlData(raw []byte, d *document) {
	z := html.NewTokenizer(bytes.NewReader(raw))
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return
		}
		if tt != html.StartTagToken && tt != html.SelfClosingTagToken {
			continue
		}
		t := z.Token()
		for _, a := range t.Attr {
			switch a.Key {
			case "href", "src":
				d.links = append(d.links, a.Val)
			case "id":
				d.ids[a.Val] = true
			case "name":
				if t.Data == "a" {
					d.ids[a.Val] = true
				}
			}
		}
	}
}
func slug(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if r == ' ' || r == '\t' {
			b.WriteByte('-')
		} else if r == '-' || r == '_' || unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsMark(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}
func parse(path string, src []byte) document {
	d := document{ids: map[string]bool{}}
	if strings.HasSuffix(path, ".html") {
		htmlData(src, &d)
		return d
	}
	root := goldmark.New(goldmark.WithExtensions(extension.GFM)).Parser().Parse(text.NewReader(src))
	used := map[string]bool{}
	ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n := n.(type) {
		case *ast.Link:
			d.links = append(d.links, string(n.Destination))
		case *ast.Image:
			d.links = append(d.links, string(n.Destination))
		case *ast.AutoLink:
			d.links = append(d.links, string(n.URL(src)))
		case *ast.Heading:
			s := slug(string(n.Text(src)))
			id := s
			for i := 1; used[id]; i++ {
				id = fmt.Sprintf("%s-%d", s, i)
			}
			used[id] = true
			d.ids[id] = true
		case *ast.RawHTML:
			var b bytes.Buffer
			for i := 0; i < n.Segments.Len(); i++ {
				s := n.Segments.At(i)
				b.Write(s.Value(src))
			}
			htmlData(b.Bytes(), &d)
		case *ast.HTMLBlock:
			var b bytes.Buffer
			for i := 0; i < n.Lines().Len(); i++ {
				s := n.Lines().At(i)
				b.Write(s.Value(src))
			}
			if n.HasClosure() {
				b.Write(n.ClosureLine.Value(src))
			}
			htmlData(b.Bytes(), &d)
		}
		return ast.WalkContinue, nil
	})
	return d
}
func localTarget(source string, u *url.URL, pages string) (string, bool) {
	if u.Host == "github.com" && strings.HasPrefix(u.Path, "/mrchypark/elysdui/") {
		for _, prefix := range []string{"/mrchypark/elysdui/blob/main/", "/mrchypark/elysdui/tree/main/"} {
			if strings.HasPrefix(u.Path, prefix) {
				return strings.TrimPrefix(u.Path, prefix), true
			}
		}
	}
	if u.Host == "mrchypark.github.io" && strings.HasPrefix(u.Path, "/elysdui/") && pages != "" {
		return filepath.Join(pages, strings.TrimPrefix(u.Path, "/elysdui/")), true
	}
	if u.Scheme != "" || u.Host != "" {
		return "", false
	}
	if source == "site/index.html" && pages != "" && u.Path != "" {
		return filepath.Join(pages, strings.TrimPrefix(u.Path, "/elysdui/")), true
	}
	// Leading slashes in GitHub-rendered Markdown refer to github.com, not the repository root.
	if strings.HasPrefix(u.Path, "/") {
		return "", false
	}
	if u.Path == "" {
		return source, true
	}
	return filepath.Clean(filepath.Join(filepath.Dir(source), u.Path)), true
}
func checkFile(path, fragment string) error {
	st, err := os.Stat(path)
	if err != nil {
		return err
	}
	if st.IsDir() {
		for _, name := range []string{"index.html", "README.md", "readme.md"} {
			candidate := filepath.Join(path, name)
			if _, err := os.Stat(candidate); err == nil {
				return checkFile(candidate, fragment)
			}
		}
		if fragment != "" {
			return fmt.Errorf("directory has no document for anchor #%s", fragment)
		}
		return nil
	}
	if fragment == "" {
		return nil
	}
	if !strings.HasSuffix(path, ".md") && !strings.HasSuffix(path, ".html") {
		return nil
	} // source-line fragments are GitHub UI, not document anchors
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	d := parse(path, src)
	if !d.ids[fragment] {
		return fmt.Errorf("missing anchor #%s in %s", fragment, path)
	}
	return nil
}

var sourceLine = regexp.MustCompile(`^L[1-9][0-9]*(?:-L[1-9][0-9]*)?$`)

func main() {
	external := flag.Bool("external", false, "also GET public HTTP(S) links (blocked/timeouts reported separately)")
	pages := flag.String("pages", "", "optional built Pages directory, e.g. dist")
	flag.Parse()
	output, err := exec.Command("git", "ls-files", "-z").Output()
	if err != nil {
		panic(err)
	}
	paths := strings.Split(string(output), "\x00")
	sort.Strings(paths)
	remote := map[string]bool{}
	failures, count, docs := 0, 0, 0
	for _, path := range paths {
		if !strings.HasSuffix(path, ".md") && path != "site/index.html" {
			continue
		}
		docs++
		src, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		d := parse(path, src)
		for _, raw := range d.links {
			count++
			u, err := url.Parse(raw)
			if err != nil {
				fmt.Printf("BROKEN %s: %s: %v\n", path, raw, err)
				failures++
				continue
			}
			if path != "site/index.html" && u.Scheme == "" && u.Host == "" && strings.HasPrefix(u.Path, "/") {
				u.Scheme, u.Host = "https", "github.com"
			}
			if u.Scheme == "mailto" {
				continue
			}
			if path == "site/index.html" && *pages == "" && u.Scheme == "" && u.Path != "" {
				continue
			} // generated routes are checked against --pages output
			target, local := localTarget(path, u, *pages)
			if local {
				fragment := u.Fragment
				if u.Host == "github.com" && sourceLine.MatchString(fragment) {
					fragment = ""
				}
				if err := checkFile(target, fragment); err != nil {
					fmt.Printf("BROKEN %s: %s: %v\n", path, raw, err)
					failures++
				}
				continue
			}
			if u.Scheme == "http" || u.Scheme == "https" {
				host := u.Hostname()
				if host == "localhost" || host == "127.0.0.1" || host == "10.0.2.2" || strings.HasSuffix(host, ".example") || host == "example.com" {
					continue
				}
				u.Fragment = ""
				remote[u.String()] = true
			}
		}
	}
	fmt.Printf("Checked %d links in %d tracked documents; %d local failures; %d distinct public URLs.\n", count, docs, failures, len(remote))
	if *external {
		client := http.Client{Timeout: 15 * time.Second}
		urls := []string{}
		for u := range remote {
			urls = append(urls, u)
		}
		sort.Strings(urls)
		for _, u := range urls {
			req, _ := http.NewRequest("GET", u, nil)
			req.Header.Set("User-Agent", "elysdui-documentation-link-check/1.0")
			res, err := client.Do(req)
			if err != nil {
				fmt.Printf("UNVERIFIED %s: %v\n", u, err)
				continue
			}
			io.Copy(io.Discard, io.LimitReader(res.Body, 4096))
			res.Body.Close()
			status := "OK"
			if res.StatusCode == 404 || res.StatusCode == 410 {
				status = "BROKEN"
				failures++
			} else if res.StatusCode < 200 || res.StatusCode >= 400 {
				status = "UNVERIFIED"
			}
			fmt.Printf("%s %d %s -> %s\n", status, res.StatusCode, u, res.Request.URL)
		}
	}
	if failures > 0 {
		os.Exit(1)
	}
}
