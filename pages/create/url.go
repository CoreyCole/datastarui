package create

import (
	"fmt"
	"math/rand"

	"github.com/coreycole/datastarui/registry"
)

// buildURL creates a URL with one param changed
func buildURL(config registry.DesignSystemConfig, key, value string) string {
	c := config
	switch key {
	case "baseColor":
		c.BaseColor = value
	case "theme":
		c.Theme = value
	case "radius":
		c.Radius = value
	case "font":
		c.Font = value
	case "style":
		c.Style = value
	case "menuColor":
		c.MenuColor = value
	case "menuAccent":
		c.MenuAccent = value
	case "iconLibrary":
		c.IconLibrary = value
	}
	return buildConfigURL(c)
}

// buildConfigURL creates a URL with all config params
func buildConfigURL(config registry.DesignSystemConfig) string {
	return fmt.Sprintf("/create?baseColor=%s&theme=%s&radius=%s&font=%s&style=%s&menuColor=%s&menuAccent=%s&iconLibrary=%s",
		config.BaseColor, config.Theme, config.Radius, config.Font, config.Style, config.MenuColor, config.MenuAccent, config.IconLibrary)
}

// buildRandomConfigURL generates a URL with random valid config values
func buildRandomConfigURL() string {
	themes := registry.AllThemes()
	baseColors := registry.BaseColors()
	fonts := registry.AllFonts()
	styles := registry.AllStyles()
	radii := registry.AllRadii()
	menuColors := registry.AllMenuColors()
	menuAccents := registry.AllMenuAccents()
	iconLibs := registry.AllIconLibraries()

	// Filter to only accent themes (non-base colors)
	var accentThemes []registry.Theme
	for _, t := range themes {
		if !t.IsBaseColor {
			accentThemes = append(accentThemes, t)
		}
	}

	config := registry.DesignSystemConfig{
		BaseColor:   baseColors[rand.Intn(len(baseColors))].Name,
		Theme:       accentThemes[rand.Intn(len(accentThemes))].Name,
		Radius:      radii[rand.Intn(len(radii))].Name,
		Font:        fonts[rand.Intn(len(fonts))].Value,
		Style:       styles[rand.Intn(len(styles))].Name,
		MenuColor:   menuColors[rand.Intn(len(menuColors))].Name,
		MenuAccent:  menuAccents[rand.Intn(len(menuAccents))].Name,
		IconLibrary: iconLibs[rand.Intn(len(iconLibs))].Name,
	}

	return buildConfigURL(config)
}

// getClickAction returns the data-on:click action that updates URL and fetches content
func getClickAction(url string) string {
	return fmt.Sprintf("history.replaceState(null, '', '%s'); @get('%s')", url, url)
}
