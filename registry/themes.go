package registry

import (
	_ "embed"
	"encoding/json"
)

//go:embed themes.json
var themesJSON []byte

// Theme represents a color theme
type Theme struct {
	Name        string  `json:"name"`
	Title       string  `json:"title"`
	IsBaseColor bool    `json:"isBaseColor"`
	CSSVars     CSSVars `json:"cssVars"`
}

// CSSVars holds light and dark mode CSS variables
type CSSVars struct {
	Light map[string]string `json:"light"`
	Dark  map[string]string `json:"dark"`
}

type themesData struct {
	Themes []Theme `json:"themes"`
}

var (
	themes     []Theme
	themeMap   map[string]*Theme
	baseColors []Theme
)

func init() {
	var data themesData
	if err := json.Unmarshal(themesJSON, &data); err != nil {
		panic("failed to parse themes.json: " + err.Error())
	}

	themes = data.Themes
	themeMap = make(map[string]*Theme, len(themes))

	for i := range themes {
		themeMap[themes[i].Name] = &themes[i]
		if themes[i].IsBaseColor {
			baseColors = append(baseColors, themes[i])
		}
	}
}

// GetTheme returns a theme by name
func GetTheme(name string) *Theme {
	return themeMap[name]
}

// GetBaseColor returns a base color by name
func GetBaseColor(name string) *Theme {
	t := themeMap[name]
	if t != nil && t.IsBaseColor {
		return t
	}
	return nil
}

// AllThemes returns all themes
func AllThemes() []Theme { return themes }

// BaseColors returns the 4 base color themes
func BaseColors() []Theme { return baseColors }

// AccentThemes returns themes available for a given base color
// Includes the matching base color + all non-base themes
func AccentThemes(baseColorName string) []Theme {
	var result []Theme
	for _, t := range themes {
		// Include the matching base color + all non-base themes
		if t.Name == baseColorName || !t.IsBaseColor {
			result = append(result, t)
		}
	}
	return result
}
