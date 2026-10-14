package sdui

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"sort"
	"strings"
	"unicode"

	ui "github.com/mrchypark/ely-htmx"
)

// RenderContext is host-owned request data. Endpoint is a same-origin absolute
// path (no query, fragment or authority); it cannot come from the schema.
// Bound values are read from nested map[string]any values without reflection,
// expression evaluation, method calls, array indexing or coercion to HTML.
type RenderContext struct {
	Data     map[string]any
	CSRF     string
	Endpoint string
}

// boundedBuffer caps every intermediate fragment, not only the final result.
// Its WriteString override also covers host writes outside html/template.
type boundedBuffer struct{ bytes.Buffer }

func (b *boundedBuffer) Write(data []byte) (int, error) {
	if len(data) > MaxRenderBytes-b.Len() {
		return 0, fmt.Errorf("sdui: rendered output exceeds %d bytes", MaxRenderBytes)
	}
	return b.Buffer.Write(data)
}
func (b *boundedBuffer) WriteString(data string) (int, error) { return b.Write([]byte(data)) }

type compiledNode struct {
	id, controlID string
	definition    Definition
	props         map[string]any
	bindings      map[string]string
	children      []*compiledNode
	event         *Event
}

// Plan is a validated, immutable snapshot. It is safe to render concurrently
// when the host does not mutate RenderContext.Data during rendering.
type Plan struct {
	id, loader string
	root       *compiledNode
	regions    map[string]*compiledNode
	events     []Event
	families   []string
	templates  *template.Template
}

func (p *Plan) ScreenID() string         { return p.id }
func (p *Plan) Loader() string           { return p.loader }
func (p *Plan) HasRegion(id string) bool { _, ok := p.regions[id]; return ok }
func (p *Plan) Events() []Event          { return append([]Event(nil), p.events...) }
func (p *Plan) Families() []string       { return append([]string(nil), p.families...) }

// ValidateActions checks every schema event against the host allowlist. Each
// action maps to the exact region IDs it may replace. HTTP handlers must call
// this before serving a plan; compilation alone cannot know host capabilities.
func (p *Plan) ValidateActions(allowed map[string][]string) error {
	for _, event := range p.events {
		targets, ok := allowed[event.Action]
		if !ok {
			return fmt.Errorf("sdui: action %q is not registered", event.Action)
		}
		found := false
		for _, target := range targets {
			if target == event.Target {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("sdui: action %q cannot replace region %q", event.Action, event.Target)
		}
	}
	return nil
}

// Compile checks structure, scalars, bindings, events and in-screen targets.
// No templates execute and no loader or action runs during compilation.
func Compile(screen Screen, r *Registry) (*Plan, error) {
	if screen.Version != Version {
		return nil, fmt.Errorf("sdui: unsupported schema version %q", screen.Version)
	}
	if !validNodeID(screen.ID) {
		return nil, fmt.Errorf("sdui: invalid screen ID %q", screen.ID)
	}
	if screen.Loader != "" && !registryKey.MatchString(screen.Loader) {
		return nil, fmt.Errorf("sdui: invalid loader key %q", screen.Loader)
	}
	if r == nil || r.templates == nil {
		return nil, fmt.Errorf("sdui: registry requires templates")
	}
	r.mu.RLock()
	definitions := make(map[string]Definition, len(r.definitions))
	for name, d := range r.definitions {
		definitions[name] = d
	}
	tpl := r.templates
	r.mu.RUnlock()
	p := &Plan{id: screen.ID, loader: screen.Loader, regions: map[string]*compiledNode{}, templates: tpl}
	seen := map[string]bool{screen.ID: true}
	families := map[string]bool{}
	count := 0
	var compileNode func(Node, int, bool) (*compiledNode, error)
	compileNode = func(node Node, depth int, inForm bool) (*compiledNode, error) {
		count++
		if count > MaxNodes {
			return nil, fmt.Errorf("sdui: node count exceeds %d", MaxNodes)
		}
		if depth > MaxDepth {
			return nil, fmt.Errorf("sdui: tree depth exceeds %d", MaxDepth)
		}
		if !validNodeID(node.ID) {
			return nil, fmt.Errorf("sdui: invalid node ID %q", node.ID)
		}
		if seen[node.ID] {
			return nil, fmt.Errorf("sdui: duplicate ID %q", node.ID)
		}
		seen[node.ID] = true
		definition, ok := definitions[node.Type]
		if !ok {
			return nil, fmt.Errorf("sdui: node %s has unknown type %q", node.ID, node.Type)
		}
		if definition.Container == "" && len(node.Children) > 0 {
			return nil, fmt.Errorf("sdui: leaf %s cannot contain children", node.ID)
		}
		if definition.Container == Form && inForm {
			return nil, fmt.Errorf("sdui: nested form %s", node.ID)
		}
		if len(node.Props) > 64 || len(node.Bindings) > 64 {
			return nil, fmt.Errorf("sdui: too many properties on %s", node.ID)
		}
		c := &compiledNode{id: node.ID, controlID: fmt.Sprintf("sdui-internal-%x-%d", sha256.Sum256([]byte(screen.ID)), count), definition: definition, props: map[string]any{}, bindings: map[string]string{}}
		for key, value := range node.Props {
			property, ok := definition.Properties[key]
			if !ok || reservedProperty(key) {
				return nil, fmt.Errorf("sdui: node %s has unknown or reserved property %q", node.ID, key)
			}
			normalized, err := normalize(property, value)
			if err != nil {
				return nil, fmt.Errorf("sdui: node %s property %s: %w", node.ID, key, err)
			}
			c.props[key] = normalized
		}
		for key, path := range node.Bindings {
			if _, ok := definition.Properties[key]; !ok || reservedProperty(key) {
				return nil, fmt.Errorf("sdui: node %s has unknown binding property %q", node.ID, key)
			}
			if _, exists := node.Props[key]; exists {
				return nil, fmt.Errorf("sdui: node %s property %s has both a literal and binding", node.ID, key)
			}
			if len(path) > 512 || !bindingPath.MatchString(path) {
				return nil, fmt.Errorf("sdui: invalid binding path %q", path)
			}
			for _, part := range strings.Split(path, ".") {
				if part == "constructor" || part == "prototype" || part == "__proto__" {
					return nil, fmt.Errorf("sdui: reserved binding path %q", path)
				}
			}
			c.bindings[key] = path
		}
		for key, property := range definition.Properties {
			_, literal := c.props[key]
			_, bound := c.bindings[key]
			if !literal && !bound {
				if property.Default != nil {
					c.props[key] = property.Default
				} else if property.Required {
					return nil, fmt.Errorf("sdui: node %s requires property %s", node.ID, key)
				}
			}
		}
		if node.Event != nil {
			event := *node.Event
			found := false
			for _, on := range definition.Events {
				if on == event.On {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("sdui: node %s does not support event %q", node.ID, event.On)
			}
			if !registryKey.MatchString(event.Action) || !validNodeID(event.Target) {
				return nil, fmt.Errorf("sdui: invalid action or target on %s", node.ID)
			}
			c.event = &event
			p.events = append(p.events, event)
		}
		if definition.Container == Region {
			p.regions[node.ID] = c
		}
		if definition.Family != "" {
			families[definition.Family] = true
		}
		for _, child := range node.Children {
			compiled, err := compileNode(child, depth+1, inForm || definition.Container == Form)
			if err != nil {
				return nil, err
			}
			c.children = append(c.children, compiled)
		}
		return c, nil
	}
	root, err := compileNode(screen.Root, 1, false)
	if err != nil {
		return nil, err
	}
	p.root = root
	for _, event := range p.events {
		if !p.HasRegion(event.Target) {
			return nil, fmt.Errorf("sdui: target %q is not a screen region", event.Target)
		}
	}
	// Workspace components reuse forms styling. The stylesheet closure is
	// collected once for the whole screen, including currently empty regions.
	for family := range families {
		p.families = append(p.families, family)
	}
	sort.Strings(p.families)
	return p, nil
}

func validNodeID(id string) bool {
	return schemaID.MatchString(id) && !strings.HasPrefix(id, "sdui-internal-")
}

var structureTemplates = template.Must(template.New("sdui").Parse(`
{{define "event"}}{{if .}} data-sdui-on="{{.On}}" data-sdui-action="{{.Action}}" data-sdui-target="{{.Target}}"{{end}}{{end}}
{{define "screen"}}<div id="{{.ID}}" data-sdui-screen="{{.ID}}"{{if .Endpoint}} data-sdui-endpoint="{{.Endpoint}}"{{end}}{{if .CSRF}} data-sdui-csrf="{{.CSRF}}"{{end}}>{{.Body}}</div>{{end}}
{{define "stack"}}<div id="{{.ID}}" class="sdui-stack" data-sdui-gap="{{.Props.Gap}}" data-sdui-direction="{{.Props.Direction}}" data-sdui-align="{{.Props.Align}}"{{template "event" .Event}}>{{.Body}}</div>{{end}}
{{define "region"}}<section id="{{.ID}}" data-sdui-region="{{.ID}}"{{if .Props.Label}} aria-label="{{.Props.Label}}"{{end}}{{template "event" .Event}}>{{.Body}}</section>{{end}}
{{define "form"}}<form id="{{.ID}}" method="post"{{if .Props.Label}} aria-label="{{.Props.Label}}"{{end}}{{template "event" .Event}}>{{if .CSRF}}<input type="hidden" name="_csrf" value="{{.CSRF}}">{{end}}{{.Body}}</form>{{end}}
{{define "leaf"}}<div id="{{.ID}}" data-sdui-node="{{.Type}}"{{template "event" .Event}}>{{.Body}}</div>{{end}}
`))

// Render renders only screen markup. Call RenderAssets once in the document
// head; partial responses must never execute or duplicate loader scripts.
func (p *Plan) Render(w io.Writer, ctx RenderContext) error {
	if p == nil || p.root == nil || p.templates == nil {
		return fmt.Errorf("sdui: rendering requires a compiled plan")
	}
	if err := validateEndpoint(ctx.Endpoint); err != nil {
		return err
	}
	var body, output boundedBuffer
	if err := p.renderNode(&body, p.root, ctx); err != nil {
		return err
	}
	if err := structureTemplates.ExecuteTemplate(&output, "screen", map[string]any{"ID": p.id, "Endpoint": ctx.Endpoint, "CSRF": ctx.CSRF, "Body": template.HTML(body.String())}); err != nil {
		return err
	}
	_, err := w.Write(output.Bytes())
	return err
}

// RenderRegion returns the complete replacement region element. All rendering
// and data validation finish before writing, so failed bindings cannot produce
// a successful partial fragment. Writer errors are returned to the host.
func (p *Plan) RenderRegion(w io.Writer, id string, ctx RenderContext) error {
	if p == nil || p.root == nil || p.templates == nil {
		return fmt.Errorf("sdui: rendering requires a compiled plan")
	}
	node, ok := p.regions[id]
	if !ok {
		return fmt.Errorf("sdui: unknown region %q", id)
	}
	var output boundedBuffer
	if err := p.renderNode(&output, node, ctx); err != nil {
		return err
	}
	_, err := w.Write(output.Bytes())
	return err
}

func (p *Plan) renderNode(w io.Writer, node *compiledNode, ctx RenderContext) error {
	props := make(map[string]any, len(node.props)+len(node.bindings)+1)
	for key, value := range node.props {
		props[node.definition.Properties[key].Target] = value
	}
	for key, path := range node.bindings {
		value, err := resolveBinding(ctx.Data, path)
		if err != nil {
			return fmt.Errorf("sdui: node %s: %w", node.id, err)
		}
		property := node.definition.Properties[key]
		value, err = normalize(property, value)
		if err != nil {
			return fmt.Errorf("sdui: node %s binding %s: %w", node.id, key, err)
		}
		props[property.Target] = value
	}
	props["ID"] = node.controlID
	var body boundedBuffer
	renderer := "leaf"
	if node.definition.Container != "" {
		renderer = string(node.definition.Container)
		for _, child := range node.children {
			if err := p.renderNode(&body, child, ctx); err != nil {
				return err
			}
		}
	} else if node.definition.Component != "" {
		if err := ui.RenderComponent(p.templates, &body, ui.Component{Chapter: node.definition.Chapter, Name: node.definition.Component, ID: node.controlID, Props: props}); err != nil {
			return err
		}
	} else if err := p.templates.ExecuteTemplate(&body, node.definition.Template, props); err != nil {
		return err
	}
	return structureTemplates.ExecuteTemplate(w, renderer, map[string]any{"ID": node.id, "Type": node.definition.Name, "Props": props, "Event": node.event, "CSRF": ctx.CSRF, "Body": template.HTML(body.String())})
}

func resolveBinding(data map[string]any, path string) (any, error) {
	var value any = data
	for _, key := range strings.Split(path, ".") {
		object, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("binding %q does not traverse an object", path)
		}
		value, ok = object[key]
		if !ok {
			return nil, fmt.Errorf("binding %q is missing", path)
		}
	}
	return value, nil
}

func validateEndpoint(endpoint string) error {
	if endpoint == "" {
		return nil
	}
	if !strings.HasPrefix(endpoint, "/") || strings.Contains(endpoint, "//") || strings.ContainsAny(endpoint, "\\?#") || strings.IndexFunc(endpoint, func(r rune) bool { return unicode.IsControl(r) || unicode.IsSpace(r) }) >= 0 {
		return fmt.Errorf("sdui: endpoint must be a same-origin absolute path")
	}
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Host != "" || parsed.Scheme != "" || parsed.RawQuery != "" || parsed.Fragment != "" || parsed.Path == "" || parsed.Path != endpoint {
		return fmt.Errorf("sdui: invalid endpoint path")
	}
	// Encoded slashes, backslashes and control bytes must not obscure authority
	// or redirect-like paths after a proxy decodes them.
	if strings.Contains(parsed.Path, "//") || strings.Contains(parsed.Path, "\\") || strings.IndexFunc(parsed.Path, unicode.IsControl) >= 0 {
		return fmt.Errorf("sdui: invalid encoded endpoint path")
	}
	return nil
}

// RenderAssets emits the fixed basic assets, the union of registered leaf
// family assets and their CSS dependencies, and the event bridge. Schemas can
// never select an asset URL. Call once per page, before any region swaps.
func (p *Plan) RenderAssets(w io.Writer) error {
	if p == nil || p.root == nil || p.templates == nil {
		return fmt.Errorf("sdui: rendering requires a compiled plan")
	}
	var output boundedBuffer
	output.WriteString(`<link rel="stylesheet" href="/assets/app.css"><link rel="stylesheet" href="/assets/sdui.css">`)
	if err := p.templates.ExecuteTemplate(&output, "ui-component-assets", nil); err != nil {
		return err
	}
	css := map[string]bool{}
	if len(p.families) > 0 {
		css["ely-api.css"] = true
	}
	for _, family := range p.families {
		css["ely-api-"+family+".css"] = true
		switch family {
		case "workspace":
			css["ely-api-forms.css"] = true
		case "communication":
			css["ely-api-calendar-parity.css"] = true
		case "native":
			css["ely-api-native-tooling.css"] = true
		}
	}
	names := make([]string, 0, len(css))
	for name := range css {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintf(&output, `<link rel="stylesheet" href="/assets/%s">`, name)
	}
	for _, family := range p.families {
		if family == "interaction" {
			output.WriteString(`<script src="/assets/lottie-light.min.js" defer></script>`)
		}
		fmt.Fprintf(&output, `<script type="module" src="/assets/ely-api-%s-entry.js"></script>`, family)
	}
	output.WriteString(`<script src="/assets/sdui.js" defer></script>`)
	_, err := w.Write(output.Bytes())
	return err
}
