package registry

import (
	_ "embed"
	"encoding/json"
)

//go:embed fonts.json
var fontsJSON []byte

// Font represents a font option
type Font struct {
	Name       string `json:"name"`
	Value      string `json:"value"`
	Type       string `json:"type"` // "sans" or "mono"
	File       string `json:"file"`
	FontFamily string `json:"fontFamily"`
}

type fontsData struct {
	Fonts []Font `json:"fonts"`
}

var (
	fonts   []Font
	fontMap map[string]*Font
)

func init() {
	var data fontsData
	if err := json.Unmarshal(fontsJSON, &data); err != nil {
		panic("failed to parse fonts.json: " + err.Error())
	}

	fonts = data.Fonts
	fontMap = make(map[string]*Font, len(fonts))
	for i := range fonts {
		fontMap[fonts[i].Value] = &fonts[i]
	}
}

// GetFont returns a font by its value (e.g., "inter", "geist")
func GetFont(value string) *Font { return fontMap[value] }

// AllFonts returns all available fonts
func AllFonts() []Font { return fonts }
