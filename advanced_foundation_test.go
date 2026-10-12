package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestFoundationAPIRendersEveryAssignedPrimary(t *testing.T) {
	tpl, err := NewTemplates()
	if err != nil {
		t.Fatal(err)
	}
	chapters, err := Catalogue()
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, chapter := range chapters {
		if chapter.Slug != "primitives" && chapter.Slug != "typography" && chapter.Slug != "layout" {
			continue
		}
		for index, story := range chapter.Stories {
			count++
			t.Run(chapter.Slug+"/"+story.Name, func(t *testing.T) {
				var out bytes.Buffer
				props := map[string]any{"ID": "foundation-test", "Kind": story.Name, "Variant": story.Name, "Title": "Caller title", "Text": "Caller text", "Label": "Caller label", "Value": "<caller>", "Name": "value", "Key": "K", "Checked": index%2 == 0, "Items": []string{"one", "two"}}
				if err := tpl.ExecuteTemplate(&out, "ely-api-"+chapter.Slug, props); err != nil {
					t.Fatal(err)
				}
				if strings.TrimSpace(out.String()) == "" || !strings.Contains(out.String(), "Caller") {
					t.Fatalf("%s did not render caller data: %q", story.Name, out.String())
				}
			})
		}
	}
	if count != 98 {
		t.Fatalf("assigned primary count = %d, want 98; update Foundation fixtures and checks", count)
	}
}

func TestFoundationAPIEscapesCallerTextAndKeepsNativeState(t *testing.T) {
	tpl, err := NewTemplates()
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	err = tpl.ExecuteTemplate(&out, "ely-api-typography", map[string]any{"ID": "safe", "Kind": "Label", "Title": "<script>", "Label": "Name", "Name": "title", "Value": "<img>", "Required": true, "Disabled": true})
	if err != nil {
		t.Fatal(err)
	}
	body := out.String()
	for _, want := range []string{"&lt;script&gt;", "&lt;img&gt;", "required", "disabled", `name="title"`} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %q in %s", want, body)
		}
	}
}

func TestFoundationInstancesKeepIDsAndDisabledLayoutState(t *testing.T) {
	tpl, err := NewTemplates()
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	for _, component := range []Component{
		{Chapter: "layout", Name: "SplitPane", ID: "foundation-one", Props: map[string]any{"Title": "<first>", "Label": "Resize", "Before": "A", "After": "B", "Disabled": true}},
		{Chapter: "layout", Name: "SplitPane", ID: "foundation-two", Props: map[string]any{"Title": "<second>", "Label": "Resize", "Before": "C", "After": "D", "Disabled": false}},
	} {
		if err := RenderComponent(tpl, &out, component); err != nil {
			t.Fatal(err)
		}
	}
	body := out.String()
	for _, want := range []string{`id="foundation-one"`, `id="foundation-two"`, "&lt;first&gt;", "&lt;second&gt;", "disabled"} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %q in %s", want, body)
		}
	}
}

func TestFoundationAssetsUseScopedRegistration(t *testing.T) {
	data, err := Files.ReadFile("assets/ely-api-foundation.js")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"import {register, emit}", "register('foundation'", "foundationMounted", "emit(root, 'change'", "Intl.NumberFormat", "Intl.RelativeTimeFormat", "data-foundation-dropzone", "drawCanvas", "requestAnimationFrame", "ely:cleanup"} {
		if !strings.Contains(string(data), want) {
			t.Errorf("asset lacks %q", want)
		}
	}
}
