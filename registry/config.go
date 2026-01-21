package registry

import (
	"fmt"
	"sort"
	"strings"
)

// DesignSystemConfig holds all customization options
type DesignSystemConfig struct {
	BaseColor   string `json:"baseColor"`   // neutral, stone, zinc, gray
	Theme       string `json:"theme"`       // neutral, blue, green, red, etc.
	Radius      string `json:"radius"`      // default, none, small, medium, large
	Font        string `json:"font"`        // inter, geist, etc.
	Style       string `json:"style"`       // vega, nova, maia, lyra, mira
	MenuColor   string `json:"menuColor"`   // default, inverted
	MenuAccent  string `json:"menuAccent"`  // subtle, bold
	IconLibrary string `json:"iconLibrary"` // lucide, phosphor, hugeicons, tabler
}

// DefaultConfig returns the default configuration
func DefaultConfig() DesignSystemConfig {
	return DesignSystemConfig{
		BaseColor:   "neutral",
		Theme:       "neutral",
		Radius:      "default",
		Font:        "geist",
		Style:       "vega",
		MenuColor:   "default",
		MenuAccent:  "subtle",
		IconLibrary: "lucide",
	}
}

// Validate checks if the config has valid values and returns a corrected copy
func (c DesignSystemConfig) Validate() DesignSystemConfig {
	if GetBaseColor(c.BaseColor) == nil {
		c.BaseColor = "neutral"
	}
	if GetTheme(c.Theme) == nil {
		c.Theme = "neutral"
	}
	if GetRadius(c.Radius) == nil {
		c.Radius = "default"
	}
	if GetFont(c.Font) == nil {
		c.Font = "geist"
	}
	if GetStyle(c.Style) == nil {
		c.Style = "vega"
	}
	if GetMenuColor(c.MenuColor) == nil {
		c.MenuColor = "default"
	}
	if GetMenuAccent(c.MenuAccent) == nil {
		c.MenuAccent = "subtle"
	}
	if GetIconLibrary(c.IconLibrary) == nil {
		c.IconLibrary = "lucide"
	}
	return c
}

// BuildThemeCSS generates CSS variable declarations for light and dark modes
func BuildThemeCSS(config DesignSystemConfig) (lightCSS, darkCSS string) {
	baseColor := GetBaseColor(config.BaseColor)
	theme := GetTheme(config.Theme)
	radius := GetRadius(config.Radius)

	if baseColor == nil || theme == nil {
		return "", ""
	}

	// Merge base color + theme (theme overrides base)
	lightVars := mergeMaps(baseColor.CSSVars.Light, theme.CSSVars.Light)
	darkVars := mergeMaps(baseColor.CSSVars.Dark, theme.CSSVars.Dark)

	// Apply radius override
	if radius != nil && radius.Value != "" {
		lightVars["radius"] = radius.Value
		darkVars["radius"] = radius.Value
	}

	// Apply menu color (inverted swaps sidebar background/foreground)
	if config.MenuColor == "inverted" {
		// Light mode: dark sidebar
		if bg, ok := lightVars["sidebar"]; ok {
			if fg, ok := lightVars["sidebar-foreground"]; ok {
				lightVars["sidebar"] = fg
				lightVars["sidebar-foreground"] = bg
			}
		}
		// Dark mode: light sidebar (inverse of dark theme)
		if bg, ok := darkVars["sidebar"]; ok {
			if fg, ok := darkVars["sidebar-foreground"]; ok {
				darkVars["sidebar"] = fg
				darkVars["sidebar-foreground"] = bg
			}
		}
	}

	// Apply menu accent (bold uses primary for accent)
	if config.MenuAccent == "bold" {
		if primary, ok := lightVars["primary"]; ok {
			lightVars["accent"] = primary
			lightVars["sidebar-accent"] = primary
		}
		if primary, ok := darkVars["primary"]; ok {
			darkVars["accent"] = primary
			darkVars["sidebar-accent"] = primary
		}
	}

	// Build CSS strings
	lightCSS = buildCSSBlock(":root", lightVars)
	darkCSS = buildCSSBlock(".dark", darkVars)

	return lightCSS, darkCSS
}

func mergeMaps(base, overlay map[string]string) map[string]string {
	result := make(map[string]string, len(base)+len(overlay))
	for k, v := range base {
		result[k] = v
	}
	for k, v := range overlay {
		result[k] = v
	}
	return result
}

func buildCSSBlock(selector string, vars map[string]string) string {
	if len(vars) == 0 {
		return ""
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString(selector)
	sb.WriteString(" {\n")
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("  --%s: %s;\n", k, vars[k]))
	}
	sb.WriteString("}")
	return sb.String()
}
