package sdui

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"

	ui "github.com/mrchypark/ely-htmx"
)

// ValueType restricts schema and bound values to ordinary scalar data.
// Integer normalizes to an exact platform int; Number uses finite float64.
type ValueType string

const (
	String  ValueType = "string"
	Boolean ValueType = "boolean"
	Integer ValueType = "integer"
	Number  ValueType = "number"
)

// Property maps a schema property to a trusted template key. Enum is for
// string design tokens. Min/Max apply to numbers. A nil Default means absent.
type Property struct {
	Type     ValueType
	Target   string
	Required bool
	Default  any
	Enum     []string
	Min, Max *float64
}

// ContainerKind is a fixed host-rendered structural primitive.
type ContainerKind string

const (
	Stack  ContainerKind = "stack"
	Region ContainerKind = "region"
	Form   ContainerKind = "form"
)

// Definition registers one of: a basic template, a catalogue component, or a
// structural primitive. Templates and definitions are trusted host code; none
// of these renderer selectors can be supplied by a screen.
type Definition struct {
	Name, Template, Chapter, Component, Family string
	Container                                  ContainerKind
	Properties                                 map[string]Property
	Events                                     []string
}

// Registry is an explicit opt-in surface, independent of the whole catalogue.
// Register may be called concurrently; compiled plans retain a snapshot.
type Registry struct {
	mu          sync.RWMutex
	templates   *template.Template
	definitions map[string]Definition
}

func NewRegistry(tpl *template.Template) *Registry {
	return &Registry{templates: tpl, definitions: map[string]Definition{}}
}

var knownFamilies = map[string]bool{"foundation": true, "interaction": true, "data": true, "workspace": true, "communication": true, "native": true, "forms": true}

func reservedProperty(name string) bool {
	key := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(name, "_", ""), "-", ""))
	switch key {
	case "id", "attrs", "class", "classname", "style", "html", "innerhtml", "outerhtml", "body", "script", "javascript", "url", "href", "src", "srcset", "srcdoc", "poster", "background", "cite", "codebase", "longdesc", "manifest", "ping", "xlinkhref", "action", "formaction", "form", "formtarget", "target", "template", "chapter", "component", "componentname", "kind", "family":
		return true
	}
	return strings.HasPrefix(key, "on") || strings.HasPrefix(key, "hx") || strings.HasPrefix(key, "data")
}

// Register validates a trusted definition and copies its scalar defaults and
// property metadata. Duplicate names and dangerous passthrough props fail.
func (r *Registry) Register(d Definition) error {
	if r == nil || r.templates == nil {
		return fmt.Errorf("sdui: registry requires templates")
	}
	if !schemaID.MatchString(d.Name) {
		return fmt.Errorf("sdui: invalid definition name %q", d.Name)
	}
	kinds := 0
	if d.Template != "" {
		kinds++
	}
	if d.Component != "" {
		kinds++
	}
	if d.Container != "" {
		kinds++
	}
	if kinds != 1 {
		return fmt.Errorf("sdui: definition %q must select exactly one renderer", d.Name)
	}
	if d.Container != "" && d.Container != Stack && d.Container != Region && d.Container != Form {
		return fmt.Errorf("sdui: unknown container %q", d.Container)
	}
	if d.Template != "" && r.templates.Lookup(d.Template) == nil {
		return fmt.Errorf("sdui: unknown template %q", d.Template)
	}
	if d.Component != "" {
		component, err := ui.ResolveComponent(d.Chapter, d.Component)
		if err != nil {
			return fmt.Errorf("sdui: %w", err)
		}
		if d.Family != "" && d.Family != component.Family {
			return fmt.Errorf("sdui: component family mismatch")
		}
		d.Family = component.Family
	}
	if d.Family != "" && !knownFamilies[d.Family] {
		return fmt.Errorf("sdui: unknown asset family %q", d.Family)
	}
	properties := make(map[string]Property, len(d.Properties))
	targets := map[string]bool{}
	for key, p := range d.Properties {
		if !propertyName.MatchString(key) || reservedProperty(key) {
			return fmt.Errorf("sdui: reserved or invalid property %q", key)
		}
		if p.Target == "" {
			p.Target = key
		}
		if !propertyName.MatchString(p.Target) || reservedProperty(p.Target) {
			return fmt.Errorf("sdui: reserved or invalid target property %q", p.Target)
		}
		if targets[p.Target] {
			return fmt.Errorf("sdui: duplicate target property %q", p.Target)
		}
		targets[p.Target] = true
		if p.Type != String && p.Type != Boolean && p.Type != Integer && p.Type != Number {
			return fmt.Errorf("sdui: unknown property type %q", p.Type)
		}
		if len(p.Enum) > 0 && p.Type != String {
			return fmt.Errorf("sdui: enum requires string type")
		}
		p.Enum = append([]string(nil), p.Enum...)
		if (p.Min != nil || p.Max != nil) && p.Type != Integer && p.Type != Number {
			return fmt.Errorf("sdui: bounds require numeric type")
		}
		if p.Min != nil {
			value := *p.Min
			if math.IsNaN(value) || math.IsInf(value, 0) {
				return fmt.Errorf("sdui: invalid minimum")
			}
			p.Min = &value
		}
		if p.Max != nil {
			value := *p.Max
			if math.IsNaN(value) || math.IsInf(value, 0) {
				return fmt.Errorf("sdui: invalid maximum")
			}
			p.Max = &value
		}
		if p.Min != nil && p.Max != nil && *p.Min > *p.Max {
			return fmt.Errorf("sdui: reversed property bounds")
		}
		if p.Default != nil {
			value, err := normalize(p, p.Default)
			if err != nil {
				return fmt.Errorf("sdui: default %s: %w", key, err)
			}
			p.Default = value
		}
		properties[key] = p
	}
	d.Properties = properties
	d.Events = append([]string(nil), d.Events...)
	for _, event := range d.Events {
		if event != "click" && event != "change" && event != "submit" && event != "ely:action" && event != "ely:change" {
			return fmt.Errorf("sdui: unsupported event %q", event)
		}
		if event == "submit" && d.Container != Form {
			return fmt.Errorf("sdui: submit requires form container")
		}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.definitions[d.Name]; exists {
		return fmt.Errorf("sdui: duplicate definition %q", d.Name)
	}
	r.definitions[d.Name] = d
	return nil
}

func normalize(p Property, value any) (any, error) {
	var number float64
	switch p.Type {
	case String:
		text, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expected string")
		}
		if len(text) > MaxStringBytes {
			return nil, fmt.Errorf("string exceeds limit")
		}
		if len(p.Enum) > 0 {
			found := false
			for _, token := range p.Enum {
				if token == text {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("unsupported token %q", text)
			}
		}
		if p.Target == "Name" && !registryKey.MatchString(text) {
			return nil, fmt.Errorf("field name must start with a letter and use at most 128 letters, digits, dots, underscores or hyphens; leading underscores are reserved")
		}
		return text, nil
	case Boolean:
		b, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("expected boolean")
		}
		return b, nil
	case Integer:
		integer, err := exactInteger(value)
		if err != nil {
			return nil, err
		}
		// Compare integer bounds without a float64 round-trip, including values
		// above 2^53. Bounds themselves are the exact supplied float64 values.
		exact := new(big.Rat).SetInt64(int64(integer))
		if p.Min != nil && exact.Cmp(new(big.Rat).SetFloat64(*p.Min)) < 0 {
			return nil, fmt.Errorf("value is below minimum")
		}
		if p.Max != nil && exact.Cmp(new(big.Rat).SetFloat64(*p.Max)) > 0 {
			return nil, fmt.Errorf("value exceeds maximum")
		}
		return integer, nil
	case Number:
		switch v := value.(type) {
		case json.Number:
			if len(v) > 128 || !jsonNumeric.MatchString(string(v)) {
				return nil, fmt.Errorf("invalid JSON number")
			}
			f, err := v.Float64()
			if err != nil {
				return nil, fmt.Errorf("invalid number")
			}
			number = f
		case float64:
			number = v
		case float32:
			number = float64(v)
		default:
			integer, err := exactInteger(value)
			if err != nil {
				return nil, fmt.Errorf("expected number")
			}
			number = float64(integer)
		}
		if math.IsNaN(number) || math.IsInf(number, 0) {
			return nil, fmt.Errorf("number must be finite")
		}
		value = number
	default:
		return nil, fmt.Errorf("unknown property type")
	}
	if p.Min != nil && number < *p.Min {
		return nil, fmt.Errorf("value is below minimum")
	}
	if p.Max != nil && number > *p.Max {
		return nil, fmt.Errorf("value exceeds maximum")
	}
	return value, nil
}

func exactInteger(value any) (int, error) {
	var text string
	switch v := value.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		text = strconv.FormatInt(v, 10)
	case uint:
		text = strconv.FormatUint(uint64(v), 10)
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		text = strconv.FormatUint(uint64(v), 10)
	case uint64:
		text = strconv.FormatUint(v, 10)
	case json.Number:
		if len(v) > 128 || !jsonNumeric.MatchString(string(v)) {
			return 0, fmt.Errorf("invalid JSON number")
		}
		text = string(v)
	case float32:
		return exactInteger(float64(v))
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) || math.Trunc(v) != v || math.Abs(v) > 9007199254740991 {
			return 0, fmt.Errorf("expected exact integer")
		}
		text = strconv.FormatFloat(v, 'f', 0, 64)
	default:
		return 0, fmt.Errorf("expected integer")
	}
	if n, err := strconv.Atoi(text); err == nil {
		return n, nil
	}
	// Rat preserves decimals exactly; cap exponent and length before allocation.
	if len(text) > 128 {
		return 0, fmt.Errorf("integer exceeds limit")
	}
	if index := strings.IndexAny(text, "eE"); index >= 0 {
		exp, err := strconv.Atoi(text[index+1:])
		if err != nil || exp < -308 || exp > 308 {
			return 0, fmt.Errorf("integer exponent exceeds limit")
		}
	}
	rational, ok := new(big.Rat).SetString(text)
	if !ok || !rational.IsInt() || !rational.Num().IsInt64() {
		return 0, fmt.Errorf("expected exact integer")
	}
	number := rational.Num().Int64()
	result := int(number)
	if int64(result) != number {
		return 0, fmt.Errorf("integer exceeds platform limit")
	}
	return result, nil
}

// DefaultRegistry is a small safe surface, not an implicit catalogue import.
func DefaultRegistry() (*Registry, error) {
	tpl, err := ui.NewTemplates()
	if err != nil {
		return nil, err
	}
	tpl, err = tpl.Parse(`{{define "sdui-text"}}<p>{{.Text}}</p>{{end}}`)
	if err != nil {
		return nil, err
	}
	r := NewRegistry(tpl)
	str := func(target string, required bool) Property {
		return Property{Type: String, Target: target, Required: required}
	}
	boolean := func(target string) Property { return Property{Type: Boolean, Target: target, Default: false} }
	enum := func(target, def string, values ...string) Property {
		return Property{Type: String, Target: target, Default: def, Enum: values}
	}
	definitions := []Definition{
		{Name: "stack", Container: Stack, Properties: map[string]Property{"gap": enum("Gap", "md", "none", "xs", "sm", "md", "lg"), "direction": enum("Direction", "column", "column", "row"), "align": enum("Align", "stretch", "start", "center", "end", "stretch")}},
		{Name: "region", Container: Region, Properties: map[string]Property{"label": str("Label", false)}},
		{Name: "form", Container: Form, Events: []string{"submit"}, Properties: map[string]Property{"label": str("Label", false)}},
		{Name: "text", Template: "sdui-text", Properties: map[string]Property{"text": str("Text", true)}},
		{Name: "notice", Template: "alert", Properties: map[string]Property{"text": str("Content", true), "title": str("Title", false), "urgent": boolean("Urgent")}},
		{Name: "input", Template: "input", Events: []string{"change"}, Properties: map[string]Property{"name": str("Name", true), "value": str("Value", false), "label": str("AriaLabel", true), "placeholder": str("Placeholder", false), "type": enum("Type", "text", "text", "email", "search", "tel", "number", "password", "date", "file", "hidden"), "required": boolean("Required"), "disabled": boolean("Disabled"), "readonly": boolean("Readonly"), "multiple": boolean("Multiple")}},
		{Name: "button", Template: "button", Events: []string{"click"}, Properties: map[string]Property{"label": str("Label", true), "type": enum("Type", "button", "button", "submit", "reset"), "variant": enum("Variant", "default", "default", "primary", "secondary", "outline", "ghost", "destructive"), "size": enum("Size", "default", "default", "xs", "sm", "lg"), "name": str("Name", false), "value": str("Value", false), "disabled": boolean("Disabled")}},
		{Name: "switch", Chapter: "forms", Component: "Switch", Events: []string{"change", "ely:change"}, Properties: map[string]Property{"label": str("Label", true), "name": str("Name", true), "value": str("Value", false), "checked": boolean("Checked"), "disabled": boolean("Disabled")}},
	}
	for _, definition := range definitions {
		if err := r.Register(definition); err != nil {
			return nil, err
		}
	}
	return r, nil
}
