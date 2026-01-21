package create

import (
	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"

	"github.com/coreycole/datastarui/registry"
	"github.com/coreycole/datastarui/theme"
)

// HandleCreatePage renders the create page with theme customization
func HandleCreatePage(c echo.Context) error {
	var config registry.DesignSystemConfig
	var isPreview bool

	// Check if URL has any params (preview mode)
	if hasAnyParam(c) {
		// Preview mode - use URL params
		config = registry.DesignSystemConfig{
			BaseColor:   c.QueryParam("baseColor"),
			Theme:       c.QueryParam("theme"),
			Radius:      c.QueryParam("radius"),
			Font:        c.QueryParam("font"),
			Style:       c.QueryParam("style"),
			MenuColor:   c.QueryParam("menuColor"),
			MenuAccent:  c.QueryParam("menuAccent"),
			IconLibrary: c.QueryParam("iconLibrary"),
		}.Validate()
		isPreview = true
	} else {
		// No params - use saved theme from disk
		config = theme.Get()
		isPreview = false
	}

	return CreatePage(config, isPreview).Render(c.Request().Context(), c.Response().Writer)
}

// HandleSaveTheme persists the preview config to disk
func HandleSaveTheme(c echo.Context) error {
	config := registry.DesignSystemConfig{
		BaseColor:   c.QueryParam("baseColor"),
		Theme:       c.QueryParam("theme"),
		Radius:      c.QueryParam("radius"),
		Font:        c.QueryParam("font"),
		Style:       c.QueryParam("style"),
		MenuColor:   c.QueryParam("menuColor"),
		MenuAccent:  c.QueryParam("menuAccent"),
		IconLibrary: c.QueryParam("iconLibrary"),
	}.Validate()

	// Save to disk (theme.json + static/css/theme.css)
	if err := theme.Update(config); err != nil {
		return err
	}

	// Patch the save button to show success state
	sse := datastar.NewSSE(c.Response().Writer, c.Request())
	return sse.PatchElementTempl(SaveButtonSaved())
}

func hasAnyParam(c echo.Context) bool {
	params := []string{"baseColor", "theme", "radius", "font", "style", "menuColor", "menuAccent", "iconLibrary"}
	for _, p := range params {
		if c.QueryParam(p) != "" {
			return true
		}
	}
	return false
}
