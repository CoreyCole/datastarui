package registry

import (
	_ "embed"
	"encoding/json"
)

//go:embed presets.json
var presetsJSON []byte

// Preset represents a curated combination of style, icons, and font
type Preset struct {
	Name        string             `json:"name"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Config      DesignSystemConfig `json:"config"`
}

type presetsData struct {
	Presets []Preset `json:"presets"`
}

var (
	presets   []Preset
	presetMap map[string]*Preset
)

func init() {
	var data presetsData
	if err := json.Unmarshal(presetsJSON, &data); err != nil {
		panic("failed to parse presets.json: " + err.Error())
	}

	presets = data.Presets
	presetMap = make(map[string]*Preset, len(presets))
	for i := range presets {
		presetMap[presets[i].Name] = &presets[i]
	}
}

// GetPreset returns a preset by name
func GetPreset(name string) *Preset { return presetMap[name] }

// AllPresets returns all available presets
func AllPresets() []Preset { return presets }

// MatchPreset finds a preset that matches all config values
// Returns nil if no preset matches (config is "Custom")
func MatchPreset(config DesignSystemConfig) *Preset {
	for i := range presets {
		p := &presets[i]
		if p.Config.Style == config.Style &&
			p.Config.IconLibrary == config.IconLibrary &&
			p.Config.Font == config.Font &&
			p.Config.BaseColor == config.BaseColor &&
			p.Config.Theme == config.Theme &&
			p.Config.MenuAccent == config.MenuAccent &&
			p.Config.MenuColor == config.MenuColor &&
			p.Config.Radius == config.Radius {
			return p
		}
	}
	return nil
}

// IsCustom returns true if the config doesn't match any preset
func IsCustom(config DesignSystemConfig) bool {
	return MatchPreset(config) == nil
}
