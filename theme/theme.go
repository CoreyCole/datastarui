package theme

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coreycole/datastarui/registry"
)

const (
	cssDir     = "static/css"
	filePrefix = "generated-theme."
	fileSuffix = ".css"
)

// Get reads theme.json and returns the config (or defaults if missing/invalid)
func Get() registry.DesignSystemConfig {
	data, err := os.ReadFile("theme.json")
	if err != nil {
		return registry.DefaultConfig()
	}
	var config registry.DesignSystemConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return registry.DefaultConfig()
	}
	return config
}

// GetCSSFile returns the path to the current hashed theme CSS file (e.g., "/css/generated-theme.abc123.css")
// Returns a fallback path if no hashed file exists
func GetCSSFile() string {
	pattern := filepath.Join(cssDir, filePrefix+"*"+fileSuffix)
	files, err := filepath.Glob(pattern)
	if err != nil || len(files) == 0 {
		return "/css/generated-theme.css" // fallback
	}
	// Return the web path (not filesystem path)
	return "/css/" + filepath.Base(files[0])
}

// Update writes theme.json and regenerates static/css/generated-theme.{hash}.css
func Update(config registry.DesignSystemConfig) error {
	data, _ := json.MarshalIndent(config, "", "  ")
	if err := os.WriteFile("theme.json", data, 0644); err != nil {
		return err
	}
	return regenerateCSS(config)
}

// Load initializes files on startup (creates defaults if theme.json missing)
func Load() error {
	if _, err := os.Stat("theme.json"); os.IsNotExist(err) {
		return Update(registry.DefaultConfig())
	}
	return regenerateCSS(Get())
}

func regenerateCSS(config registry.DesignSystemConfig) error {
	lightCSS, darkCSS := registry.BuildThemeCSS(config)
	font := registry.GetFont(config.Font)

	css := lightCSS + "\n\n" + darkCSS
	if font != nil {
		css += "\n\nbody { font-family: " + font.FontFamily + "; }"
	}

	// Generate content hash (first 8 chars of SHA256)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(css)))[:8]
	newFilename := filePrefix + hash + fileSuffix
	newPath := filepath.Join(cssDir, newFilename)

	// Clean up old generated-theme.*.css files
	if err := cleanupOldFiles(newFilename); err != nil {
		return err
	}

	return os.WriteFile(newPath, []byte(css), 0644)
}

func cleanupOldFiles(keepFilename string) error {
	pattern := filepath.Join(cssDir, filePrefix+"*"+fileSuffix)
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !strings.HasSuffix(f, keepFilename) {
			if err := os.Remove(f); err != nil {
				return err
			}
		}
	}
	return nil
}
