package sdui

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"

	ui "github.com/mrchypark/ely-htmx"
)

type renderBudget struct {
	nodes  int
	values *valueBudget
	ids    map[string]bool
}

func (budget *renderBudget) valueBudget() *valueBudget {
	if budget.values == nil {
		budget.values = newValueBudget()
	}
	return budget.values
}

func applyParameters(node *compiledNode, scope renderScope, budget *renderBudget) (renderScope, error) {
	params := make(map[string]any, len(node.use.params)+len(node.use.bindings))
	for _, key := range sortedKeys(node.use.params) {
		value, err := copyNormalizedValue(node.use.params[key], budget.valueBudget(), 1)
		if err != nil {
			return scope, fmt.Errorf("sdui: use %s parameter %s: %w", node.id, key, err)
		}
		params[key] = value
	}
	for _, key := range sortedKeys(node.use.bindings) {
		ref := node.use.bindings[key]
		value, present, err := ref.lookup(scope)
		if err != nil {
			return scope, fmt.Errorf("sdui: use %s: %w", node.id, err)
		}
		if !present {
			return scope, fmt.Errorf("sdui: use %s parameter binding %q is missing", node.id, ref.path)
		}
		normalized, err := normalizeWithBudget(node.use.properties[key], value, budget.valueBudget())
		if err != nil {
			return scope, fmt.Errorf("sdui: use %s parameter %s: %w", node.id, key, err)
		}
		params[key] = normalized
	}
	scope.params = params
	scope.vars = nil
	return scope, nil
}

func (p *Plan) renderNode(w io.Writer, node *compiledNode, ctx RenderContext, scope renderScope, budget *renderBudget) error {
	return forEachNodeInstance(node, scope, budget, func(instance renderScope) error {
		return p.renderInstance(w, node, ctx, instance, budget)
	})
}

// forEachNodeInstance is shared by the HTML and native-neutral projections.
// Stable identities and pre-allocation expansion limits must be identical.
func forEachNodeInstance(node *compiledNode, scope renderScope, budget *renderBudget, visit func(renderScope) error) error {
	if node.repeat == nil {
		return visit(scope)
	}
	value, present, err := node.repeat.items.lookup(scope)
	if err != nil {
		return fmt.Errorf("sdui: node %s repeat: %w", node.id, err)
	}
	if !present {
		return fmt.Errorf("sdui: node %s repeat binding %q is missing", node.id, node.repeat.items.path)
	}
	// Check before copying supported typed slices, to bound allocations too.
	count, ok := arrayLength(value)
	if !ok {
		return fmt.Errorf("sdui: node %s repeat requires an array", node.id)
	}
	if count > MaxNodes || count > MaxNodes-budget.nodes {
		return fmt.Errorf("sdui: expanded node count exceeds %d", MaxNodes)
	}
	items, _ := dataArray(value)
	seen := map[string]bool{}
	for _, item := range items {
		value, present, err := walkData(item, node.repeat.key, node.repeat.keyPath)
		if err != nil {
			return fmt.Errorf("sdui: node %s repeat key: %w", node.id, err)
		}
		if !present {
			return fmt.Errorf("sdui: node %s repeat key %q is missing", node.id, node.repeat.keyPath)
		}
		key, err := repeatKey(value)
		if err != nil {
			return fmt.Errorf("sdui: node %s: %w", node.id, err)
		}
		if seen[key] {
			return fmt.Errorf("sdui: node %s has duplicate repeat key", node.id)
		}
		seen[key] = true
		vars := make(map[string]any, len(scope.vars)+1)
		for key, value := range scope.vars {
			vars[key] = value
		}
		vars[node.repeat.as] = item
		inner := scope
		inner.vars = vars
		inner.identity = fmt.Sprintf("%x", sha256.Sum256([]byte(scope.identity+"\x00"+node.id+"\x00"+key)))
		if err := visit(inner); err != nil {
			return err
		}
	}
	return nil
}

func arrayLength(value any) (int, bool) {
	switch values := value.(type) {
	case []any:
		return len(values), true
	case []map[string]any:
		return len(values), true
	case []string:
		return len(values), true
	default:
		return 0, false
	}
}

func (p *Plan) renderInstance(w io.Writer, node *compiledNode, ctx RenderContext, scope renderScope, budget *renderBudget) error {
	budget.nodes++
	if budget.nodes > MaxNodes {
		return fmt.Errorf("sdui: expanded node count exceeds %d", MaxNodes)
	}
	show, err := node.when.evaluate(scope)
	if err != nil {
		return fmt.Errorf("sdui: node %s condition: %w", node.id, err)
	}
	if !show {
		return nil
	}
	if node.use != nil {
		inner, err := applyParameters(node, scope, budget)
		if err != nil {
			return err
		}
		return p.renderNode(w, node.use.root, ctx, inner, budget)
	}
	id := dynamicID(node.id, scope.identity)
	controlID := dynamicControlID(node.controlID, scope.identity)
	if budget.ids == nil {
		budget.ids = map[string]bool{p.id: true}
	}
	if budget.ids[id] || budget.ids[controlID] {
		return fmt.Errorf("sdui: duplicate rendered ID %q", id)
	}
	budget.ids[id] = true
	budget.ids[controlID] = true
	props := make(map[string]any, len(node.props)+len(node.bindings)+len(node.slots)+1)
	for _, key := range sortedKeys(node.props) {
		value, err := copyNormalizedValue(node.props[key], budget.valueBudget(), 1)
		if err != nil {
			return fmt.Errorf("sdui: node %s property %s: %w", node.id, key, err)
		}
		props[node.definition.Properties[key].Target] = value
	}
	for _, key := range sortedKeys(node.bindings) {
		ref := node.bindings[key]
		value, present, err := ref.lookup(scope)
		if err != nil {
			return fmt.Errorf("sdui: node %s: %w", node.id, err)
		}
		if !present {
			return fmt.Errorf("sdui: node %s binding %q is missing", node.id, ref.path)
		}
		property := node.definition.Properties[key]
		normalized, err := normalizeWithBudget(property, value, budget.valueBudget())
		if err != nil {
			return fmt.Errorf("sdui: node %s binding %s: %w", node.id, key, err)
		}
		props[property.Target] = normalized
	}
	props["ID"] = controlID
	for _, name := range sortedKeys(node.slots) {
		var content boundedBuffer
		before := budget.nodes
		for _, child := range node.slots[name] {
			if err := p.renderNode(&content, child, ctx, scope, budget); err != nil {
				return err
			}
		}
		slot := node.definition.Slots[name]
		if slot.MaxNodes > 0 && budget.nodes-before > slot.MaxNodes {
			return fmt.Errorf("sdui: node %s slot %s exceeds expanded node limit", node.id, name)
		}
		props[slot.Target] = template.HTML(content.String()) // Only markup rendered above crosses this boundary.
	}
	var body boundedBuffer
	renderer := "leaf"
	if node.definition.Container != "" {
		renderer = string(node.definition.Container)
		for _, child := range node.children {
			if err := p.renderNode(&body, child, ctx, scope, budget); err != nil {
				return err
			}
		}
	} else if node.definition.Component != "" {
		if err := ui.RenderComponentLimited(p.templates, &body, ui.Component{Chapter: node.definition.Chapter, Name: node.definition.Component, ID: controlID, Props: props}, MaxRenderBytes); err != nil {
			return err
		}
	} else if err := p.templates.ExecuteTemplate(&body, node.definition.Template, props); err != nil {
		return err
	}
	return structureTemplates.ExecuteTemplate(w, renderer, map[string]any{"ID": id, "Type": node.definition.Name, "Props": props, "Event": eventJSON(node.events), "CSRF": ctx.CSRF, "Body": template.HTML(body.String())})
}
func dynamicControlID(id, identity string) string {
	if identity == "" {
		return id
	}
	return fmt.Sprintf("sdui-internal-%x", sha256.Sum256([]byte(id+"\x00"+identity)))
}
