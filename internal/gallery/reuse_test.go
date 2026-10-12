package uiexamples

import (
	"io/fs"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestGalleryRecipesCallSharedComponents(t *testing.T) {
	paths, err := fs.Glob(Files, "templates/demo-*.html")
	if err != nil {
		t.Fatal(err)
	}
	controls := regexp.MustCompile(`<(input|textarea|select|button)\b`)
	for _, path := range paths {
		data, err := Files.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if controls.Match(data) {
			t.Errorf("gallery owns control implementation: %s", path)
		}
	}
	server, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"password-input", "pin-input", "tag-input", "key-value-input", "combobox", "multi-select", "calendar", "file-dropzone", "signature-pad", "map-view", "data-table"} {
		if _, err := server.templates.New(name).Parse(`<output data-reused="` + name + `"></output>`); err != nil {
			t.Fatal(err)
		}
	}
	for path, component := range map[string]string{
		"/components/forms/passwordinput":           "password-input",
		"/components/forms/pininput-otpinput":       "pin-input",
		"/components/forms/taginput-tokeninput":     "tag-input",
		"/components/forms/keyvalueinput-listinput": "key-value-input",
		"/components/forms/combobox-autocomplete":   "combobox",
		"/components/forms/multiselect":             "multi-select",
		"/components/forms/calendar":                "calendar",
		"/components/forms/fileinput-dropzone":      "file-dropzone",
		"/components/forms/signaturepad":            "signature-pad",
		"/components/maps/locationpicker":           "map-view",
		"/components/tables/datatable-tabletoolbar-filterbuilder-sortbuilder-columnvisibility-cellrenderer-rowselection-tabledensity-aggregationfooter-conditionalformatting": "data-table",
	} {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		if response.Code != 200 || !strings.Contains(response.Body.String(), `data-reused="`+component+`"`) {
			t.Errorf("%s does not reuse %s (%d)", path, component, response.Code)
		}
	}
}

func TestUsagePageServesIndependentConfiguredInstances(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	server.ServeHTTP(response, httptest.NewRequest("GET", "/components/usage", nil))
	if response.Code != 200 {
		t.Fatalf("usage=%d: %s", response.Code, response.Body.String())
	}
	for _, want := range []string{`id="usage-pin-a"`, `id="usage-pin-b"`, `id="usage-combo"`, `id="usage-combo-b"`, `data-ui-control-template="checkbox"`, `data-ui-table-sort`, `data-usage-form`} {
		if !strings.Contains(response.Body.String(), want) {
			t.Errorf("usage lacks %s", want)
		}
	}
}
