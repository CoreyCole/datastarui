package registry

import (
	_ "embed"
	"encoding/json"
)

//go:embed styles.json
var stylesJSON []byte

// Style represents a component styling preset
type Style struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type stylesData struct {
	Styles []Style `json:"styles"`
}

var (
	styles   []Style
	styleMap map[string]*Style
)

func init() {
	var data stylesData
	if err := json.Unmarshal(stylesJSON, &data); err != nil {
		panic("failed to parse styles.json: " + err.Error())
	}

	styles = data.Styles
	styleMap = make(map[string]*Style, len(styles))
	for i := range styles {
		styleMap[styles[i].Name] = &styles[i]
	}
}

// GetStyle returns a style by name
func GetStyle(name string) *Style { return styleMap[name] }

// AllStyles returns all available styles
func AllStyles() []Style { return styles }
