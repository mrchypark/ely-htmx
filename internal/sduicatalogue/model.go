// Package sduicatalogue builds reviewed, bounded SDUI contracts from the actual
// catalogue renderers and build-time fixtures. It is never imported by sdui.
package sduicatalogue

type Property struct {
	Type                 string              `json:"type,omitempty"`
	Target               string              `json:"target,omitempty"`
	Required             bool                `json:"required,omitempty"`
	Nullable             bool                `json:"nullable,omitempty"`
	Properties           map[string]Property `json:"properties,omitempty"`
	Items                *Property           `json:"items,omitempty"`
	AdditionalProperties *Property           `json:"additionalProperties,omitempty"`
	OneOf                []Property          `json:"oneOf,omitempty"`
	Format               string              `json:"format,omitempty"`
	Enum                 []string            `json:"enum,omitempty"`
	MaxLength            int                 `json:"maxLength,omitempty"`
	MaxItems             int                 `json:"maxItems,omitempty"`
	MaxProperties        int                 `json:"maxProperties,omitempty"`
}

type Contract struct {
	Name         string              `json:"name"`
	Template     string              `json:"template,omitempty"`
	Chapter      string              `json:"chapter,omitempty"`
	Component    string              `json:"component,omitempty"`
	Primary      string              `json:"primary,omitempty"`
	Family       string              `json:"family,omitempty"`
	ContainsForm bool                `json:"containsForm,omitempty"`
	SandboxHTML  bool                `json:"sandboxHTML,omitempty"`
	Properties   map[string]Property `json:"properties"`
	Assets       []string            `json:"assets,omitempty"`
	Events       []string            `json:"events,omitempty"`
	Excluded     []Exclusion         `json:"excluded,omitempty"`
}

type Exclusion struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}
type Manifest struct {
	Version         string     `json:"version"`
	SourceHash      string     `json:"sourceHash"`
	Primaries       int        `json:"primaries"`
	QualifiedNames  int        `json:"qualifiedNames"`
	Chapters        int        `json:"chapters"`
	BasicTemplates  int        `json:"basicTemplates"`
	SourceFixtures  int        `json:"sourceFixtures"`
	PreviewFixtures int        `json:"previewFixtures"`
	Contracts       []Contract `json:"contracts"`
}

// Fixture is a caller-owned, template-keyed input used only by the generator.
type Fixture struct {
	Name, Template, Source string
	Props                  map[string]any
}
type TestFixture struct {
	Name     string         `json:"name"`
	Source   string         `json:"source"`
	Sources  []string       `json:"sources,omitempty"`
	Props    map[string]any `json:"props"`
	Excluded []Exclusion    `json:"excluded,omitempty"`
}
