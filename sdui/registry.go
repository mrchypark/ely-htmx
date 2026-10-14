package sdui

import (
	"fmt"
	"html/template"
	"sort"
	"strings"
	"sync"

	ui "github.com/mrchypark/ely-htmx"
)

// ValueType describes ordinary, recursively typed JSON data. Integer preserves
// exact platform integers; Number normalizes to a finite float64.
type ValueType string

const (
	String  ValueType = "string"
	Boolean ValueType = "boolean"
	Integer ValueType = "integer"
	Number  ValueType = "number"
	Object  ValueType = "object"
	Array   ValueType = "array"
)

// Property maps a wire property to a trusted template key. Objects are closed
// unless AdditionalProperties declares the dictionary value type. OneOf uses
// exactly-one-match semantics and cannot also specify Type. Defaults are copied.
type Property struct {
	Type                 ValueType           `json:"type,omitempty"`
	Target               string              `json:"target,omitempty"`
	Required             bool                `json:"required,omitempty"`
	Default              any                 `json:"default,omitempty"`
	Enum                 []string            `json:"enum,omitempty"`
	Min                  *float64            `json:"min,omitempty"`
	Max                  *float64            `json:"max,omitempty"`
	Items                *Property           `json:"items,omitempty"`
	Properties           map[string]Property `json:"properties,omitempty"`
	AdditionalProperties *Property           `json:"additionalProperties,omitempty"`
	OneOf                []Property          `json:"oneOf,omitempty"`
	Nullable             bool                `json:"nullable,omitempty"`
	Format               string              `json:"format,omitempty"`
	MaxLength            int                 `json:"maxLength,omitempty"`
	MaxItems             int                 `json:"maxItems,omitempty"`
	MaxProperties        int                 `json:"maxProperties,omitempty"`
	urlPolicy            URLPolicy
	sandboxHTML          bool
	inertData            bool
}

// URLPolicy can authorize an external HTTP(S) URL after the runtime's URL
// safety checks. It must be safe to call concurrently and must not mutate data.
// The format distinguishes links, images and media resources. Local absolute
// paths and fragments do not need authorization; other schemes never reach it.
type URLPolicy func(format, rawURL string) bool

type RegistryOption func(*Registry)

// WithURLPolicy explicitly permits host-reviewed external origins. Without it,
// URLs must be local absolute paths, fragments, or empty optional values.
func WithURLPolicy(policy URLPolicy) RegistryOption {
	return func(r *Registry) { r.urlPolicy = policy }
}

// SlotDefinition grants a named template slot only compiled child markup.
// It is not a property passthrough and cannot be filled with schema strings.
type SlotDefinition struct {
	Target       string   `json:"target,omitempty"`
	Required     bool     `json:"required,omitempty"`
	AllowedTypes []string `json:"allowedTypes,omitempty"`
	MaxNodes     int      `json:"maxNodes,omitempty"`
}

type Slot = SlotDefinition

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
	Name         string                    `json:"name"`
	Template     string                    `json:"template,omitempty"`
	Chapter      string                    `json:"chapter,omitempty"`
	Component    string                    `json:"component,omitempty"`
	Family       string                    `json:"family,omitempty"`
	Container    ContainerKind             `json:"container,omitempty"`
	Properties   map[string]Property       `json:"properties,omitempty"`
	Events       []string                  `json:"events,omitempty"`
	Assets       []string                  `json:"assets,omitempty"`
	ContainsForm bool                      `json:"containsForm,omitempty"`
	Slots        map[string]SlotDefinition `json:"slots,omitempty"`
	// SandboxHTML grants the reviewed, fixed sandbox="" iframe srcdoc sink of
	// misc/WebView and misc/IframeEmbed. It never marks input as trusted HTML.
	SandboxHTML bool `json:"sandboxHTML,omitempty"`
}

// Registry is an explicit opt-in surface, independent of the whole catalogue.
// Register may be called concurrently; compiled plans retain a snapshot.
type Registry struct {
	mu          sync.RWMutex
	templates   *template.Template
	definitions map[string]Definition
	urlPolicy   URLPolicy
}

func NewRegistry(tpl *template.Template, options ...RegistryOption) *Registry {
	r := &Registry{templates: tpl, definitions: map[string]Definition{}}
	for _, option := range options {
		if option != nil {
			option(r)
		}
	}
	return r
}

var knownFamilies = map[string]bool{"foundation": true, "interaction": true, "data": true, "workspace": true, "communication": true, "native": true, "forms": true}

func reservedProperty(name string) bool {
	key := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(name, "_", ""), "-", ""))
	switch key {
	case "id", "attrs", "class", "classname", "style", "html", "innerhtml", "outerhtml", "body", "script", "javascript", "url", "href", "src", "srcset", "srcdoc", "poster", "background", "cite", "codebase", "longdesc", "manifest", "ping", "xlinkhref", "action", "formaction", "form", "formtarget", "target", "template", "chapter", "component", "componentname", "kind", "family":
		return true
	}
	return strings.HasPrefix(key, "on") || strings.HasPrefix(key, "hx")
}

// Definitions returns sorted, independent copies for editors and tooling.
// Mutating a returned property, default, or slot never changes the registry.
func (r *Registry) Definitions() []Definition {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Definition, 0, len(r.definitions))
	for _, definition := range r.definitions {
		result = append(result, cloneDefinition(definition))
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func validDefinitionName(name string) bool {
	if name == "" || len(name) > 192 || strings.TrimSpace(name) != name {
		return false
	}
	if !((name[0] >= 'A' && name[0] <= 'Z') || (name[0] >= 'a' && name[0] <= 'z')) {
		return false
	}
	for _, r := range name {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || strings.ContainsRune(" _-./()+&,:", r)) {
			return false
		}
	}
	return true
}

// Register validates a trusted definition and deeply copies its defaults and
// property metadata. Duplicate names and dangerous passthrough props fail.
func (r *Registry) Register(d Definition) error {
	if r == nil || r.templates == nil {
		return fmt.Errorf("sdui: registry requires templates")
	}
	if !validDefinitionName(d.Name) {
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
	if d.SandboxHTML {
		component, err := ui.ResolveComponent(d.Chapter, d.Component)
		if err != nil || component.Chapter != "misc" || (component.Name != "WebView" && component.Name != "IframeEmbed") {
			return fmt.Errorf("sdui: sandbox-html requires a reviewed sandboxed document component")
		}
	}
	properties := make(map[string]Property, len(d.Properties))
	targets := map[string]bool{}
	if len(d.Properties) > MaxObjectProperties {
		return fmt.Errorf("sdui: too many declared properties")
	}
	for key, p := range d.Properties {
		if !propertyAllowed(key, p) {
			return fmt.Errorf("sdui: reserved or invalid property %q", key)
		}
		if p.Target == "" {
			p.Target = key
		}
		if !propertyAllowed(p.Target, p) {
			return fmt.Errorf("sdui: reserved or invalid target property %q", p.Target)
		}
		if targets[p.Target] {
			return fmt.Errorf("sdui: duplicate target property %q", p.Target)
		}
		targets[p.Target] = true
		inert := false
		if p.Format == "inert-data" && p.Target == "Devtools" && d.Component != "" {
			component, err := ui.ResolveComponent(d.Chapter, d.Component)
			inert = err == nil && component.Chapter == "devtools" && p.Type == Object
		}
		validated, err := validatePropertyContext(p, r.urlPolicy, d.SandboxHTML, inert)
		if err != nil {
			return fmt.Errorf("sdui: property %s: %w", key, err)
		}
		properties[key] = validated
	}
	slots := make(map[string]SlotDefinition, len(d.Slots))
	for key, slot := range d.Slots {
		if !propertyName.MatchString(key) {
			return fmt.Errorf("sdui: invalid slot %q", key)
		}
		if slot.Target == "" {
			slot.Target = key
		}
		if !propertyName.MatchString(slot.Target) || hardReservedProperty(slot.Target, true) || targets[slot.Target] {
			return fmt.Errorf("sdui: invalid or conflicting slot target %q", slot.Target)
		}
		if slot.MaxNodes < 0 || slot.MaxNodes > MaxNodes {
			return fmt.Errorf("sdui: invalid slot node limit")
		}
		for _, name := range slot.AllowedTypes {
			if !validDefinitionName(name) {
				return fmt.Errorf("sdui: invalid slot type %q", name)
			}
		}
		slot.AllowedTypes = append([]string(nil), slot.AllowedTypes...)
		targets[slot.Target] = true
		slots[key] = slot
	}
	d.Slots = slots
	d.Properties = properties
	d.Assets = append([]string(nil), d.Assets...)
	seenAssets := map[string]bool{}
	for _, asset := range d.Assets {
		if asset != "map" || seenAssets[asset] {
			return fmt.Errorf("sdui: unsupported or duplicate asset %q", asset)
		}
		seenAssets[asset] = true
	}
	d.Events = append([]string(nil), d.Events...)
	for _, event := range d.Events {
		if event != "click" && event != "change" && event != "submit" && event != "ely:action" && event != "ely:change" {
			return fmt.Errorf("sdui: unsupported event %q", event)
		}
		if event == "submit" && d.Container != Form && !d.ContainsForm {
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
