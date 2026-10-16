package playground

import (
	"fmt"
	"regexp"
	"strings"
)

// Only bounded known color roles enter the hash-authorized preview CSS. No
// authored style, URL, selector or script bridge is accepted by this API.
type previewTheme struct {
	Mode   string            `json:"mode"`
	Colors map[string]string `json:"colors,omitempty"`
}

var previewColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
var previewColorRoles = []string{"color-paper", "color-ink", "color-muted", "color-line", "color-accent", "surface", "border", "accent-hover", "accent-tint", "on-accent", "focus", "draft", "draft-tint", "approved", "approved-tint", "warning", "warning-tint", "error", "error-tint"}

func (theme *previewTheme) styles() (string, string, error) {
	if theme == nil {
		return "light", "", nil
	}
	if theme.Mode != "light" && theme.Mode != "dark" {
		return "", "", fmt.Errorf("preview theme mode must be light or dark")
	}
	if len(theme.Colors) == 0 {
		return theme.Mode, "", nil
	}
	if len(theme.Colors) != len(previewColorRoles) {
		return "", "", fmt.Errorf("preview theme requires all 19 color roles")
	}
	var css strings.Builder
	css.WriteString(":root{")
	for _, role := range previewColorRoles {
		value := theme.Colors[role]
		if !previewColor.MatchString(value) {
			return "", "", fmt.Errorf("preview theme %s must be a six-digit hex color", role)
		}
		fmt.Fprintf(&css, "--%s:%s;", role, value)
	}
	css.WriteString("}")
	return theme.Mode, css.String(), nil
}
