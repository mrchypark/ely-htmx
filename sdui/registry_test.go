package sdui

import (
	"encoding/json"
	"html/template"
	"math"
	"testing"
)

func TestPropertyNormalization(t *testing.T) {
	min, max := float64(0), float64(100)
	largeMax := float64(9007199254740992)
	cases := []struct {
		name  string
		p     Property
		value any
		want  any
		fails bool
	}{
		{"string", Property{Type: String}, "hello", "hello", false},
		{"trusted HTML rejected", Property{Type: String}, template.HTML("<b>x</b>"), nil, true},
		{"trusted URL rejected", Property{Type: String}, template.URL("https://example.com"), nil, true},
		{"bool", Property{Type: Boolean}, false, false, false},
		{"string bool rejected", Property{Type: Boolean}, "false", nil, true},
		{"int", Property{Type: Integer}, 42, 42, false},
		{"JSON integer", Property{Type: Integer}, json.Number("42"), 42, false},
		{"JSON decimal integer", Property{Type: Integer}, json.Number("42.00"), 42, false},
		{"JSON scientific integer", Property{Type: Integer}, json.Number("4.2e1"), 42, false},
		{"exact fractional rejected", Property{Type: Integer}, json.Number("1.000000000000000000001"), nil, true},
		{"fraction rejected", Property{Type: Integer}, 3.2, nil, true},
		{"float integer", Property{Type: Integer}, 3.0, 3, false},
		{"precision loss rejected", Property{Type: Integer}, float64(9007199254740992), nil, true},
		{"huge exponent rejected", Property{Type: Integer}, json.Number("1e1000000000"), nil, true},
		{"tiny exponent rejected", Property{Type: Integer}, json.Number("1e-1000000000"), nil, true},
		{"overflow rejected", Property{Type: Integer}, json.Number("9223372036854775808"), nil, true},
		{"non JSON fraction rejected", Property{Type: Integer}, json.Number("1/1"), nil, true},
		{"number", Property{Type: Number}, json.Number("1.25"), 1.25, false},
		{"nan", Property{Type: Number}, math.NaN(), nil, true},
		{"infinity", Property{Type: Number}, math.Inf(1), nil, true},
		{"bounds", Property{Type: Integer, Min: &min, Max: &max}, 50, 50, false},
		{"min", Property{Type: Integer, Min: &min}, -1, nil, true},
		{"max", Property{Type: Number, Max: &max}, 101, nil, true},
		{"exact integer bound", Property{Type: Integer, Max: &largeMax}, json.Number("9007199254740993"), nil, true},
		{"enum", Property{Type: String, Enum: []string{"sm", "md"}}, "md", "md", false},
		{"enum rejected", Property{Type: String, Enum: []string{"sm", "md"}}, "<style>", nil, true},
		{"name reserved", Property{Type: String, Target: "Name"}, "_csrf", nil, true},
		{"name invalid", Property{Type: String, Target: "Name"}, "a[b]", nil, true},
		{"name allowed", Property{Type: String, Target: "Name"}, "profile.name", "profile.name", false},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := normalize(test.p, test.value)
			if test.fails {
				if err == nil {
					t.Fatalf("accepted %#v", test.value)
				}
				return
			}
			if err != nil || got != test.want {
				t.Fatalf("got %#v, %v; want %#v", got, err, test.want)
			}
		})
	}
}

func TestRegistryValidationAndSnapshot(t *testing.T) {
	tpl := template.Must(template.New("leaf").Parse(`{{.Text}}`))
	valid := func() Definition {
		return Definition{Name: "custom", Template: "leaf", Properties: map[string]Property{"text": {Type: String, Target: "Text", Default: "safe"}}}
	}
	tests := map[string]func(*Definition){
		"unknown template": func(d *Definition) { d.Template = "missing" }, "multiple renderers": func(d *Definition) { d.Container = Stack },
		"missing renderer": func(d *Definition) { d.Template = "" }, "unknown container": func(d *Definition) { d.Template = ""; d.Container = "iframe" },
		"unknown family": func(d *Definition) { d.Family = "evil" }, "invalid name": func(d *Definition) { d.Name = "<bad>" },
		"unknown component": func(d *Definition) { d.Template = ""; d.Component = "Missing" },
		"wrong family":      func(d *Definition) { d.Template = ""; d.Component = "Switch"; d.Chapter = "forms"; d.Family = "data" },
		"reserved prop":     func(d *Definition) { d.Properties["href"] = Property{Type: String} },
		"reserved target":   func(d *Definition) { d.Properties["safe"] = Property{Type: String, Target: "Attrs"} },
		"unknown type":      func(d *Definition) { d.Properties["safe"] = Property{Type: "html"} },
		"duplicate target":  func(d *Definition) { d.Properties["safe"] = Property{Type: String, Target: "Text"} },
		"invalid default":   func(d *Definition) { d.Properties["safe"] = Property{Type: Boolean, Default: "false"} },
		"wrong enum":        func(d *Definition) { d.Properties["safe"] = Property{Type: Boolean, Enum: []string{"false"}} },
		"unsupported event": func(d *Definition) { d.Events = []string{"onclick"} }, "submit on leaf": func(d *Definition) { d.Events = []string{"submit"} },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			d := valid()
			mutate(&d)
			if err := NewRegistry(tpl).Register(d); err == nil {
				t.Fatal("invalid definition accepted")
			}
		})
	}
	r := NewRegistry(tpl)
	d := valid()
	if err := r.Register(d); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(d); err == nil {
		t.Fatal("duplicate accepted")
	}
	d.Properties["text"] = Property{Type: Boolean}
	plan, err := Compile(Screen{Version: Version, ID: "screen", Root: Node{ID: "root", Type: "custom"}}, r)
	if err != nil {
		t.Fatal(err)
	}
	if plan.root.props["text"] != "safe" {
		t.Fatal("registered definition mutated")
	}
	if err := NewRegistry(nil).Register(valid()); err == nil {
		t.Fatal("nil templates accepted")
	}
}
