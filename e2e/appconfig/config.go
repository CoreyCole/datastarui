package appconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultConfigFile = "datastarui-e2e.yml"

type Config struct {
	App          string       `yaml:"app"`
	BaseURL      string       `yaml:"base_url"`
	RunPackage   string       `yaml:"run_package"`
	ArtifactsDir string       `yaml:"artifacts_dir"`
	Server       ServerConfig `yaml:"server"`
	Viewports    []string     `yaml:"viewports"`
	ConfigPath   string       `yaml:"-"`
	RootDir      string       `yaml:"-"`
}

type ServerConfig struct {
	Command            string `yaml:"command"`
	SkipWhenBaseURLSet bool   `yaml:"skip_when_base_url_set"`
}

func Find(cwd string) (string, bool, error) {
	if explicit := strings.TrimSpace(os.Getenv("E2E_CONFIG")); explicit != "" {
		abs, err := filepath.Abs(explicit)
		return abs, true, err
	}

	cur, err := filepath.Abs(cwd)
	if err != nil {
		return "", false, err
	}
	for {
		candidate := filepath.Join(cur, DefaultConfigFile)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true, nil
		} else if !os.IsNotExist(err) {
			return "", false, err
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return "", false, nil
		}
		cur = parent
	}
}

func Load(path, cwd string) (Config, error) {
	if strings.TrimSpace(path) == "" {
		found, ok, err := Find(cwd)
		if err != nil {
			return Config{}, err
		}
		if !ok {
			return Config{}, fmt.Errorf("%s not found from %s", DefaultConfigFile, cwd)
		}
		path = found
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(abs)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	cfg.ConfigPath = abs
	cfg.RootDir = filepath.Dir(abs)
	cfg.applyDefaults()
	return cfg, cfg.Validate()
}

func (c *Config) applyDefaults() {
	if c.App == "" {
		c.App = "app"
	}
	if c.ArtifactsDir == "" {
		c.ArtifactsDir = ".e2e-runs"
	}
	if c.RunPackage == "" {
		c.RunPackage = "./tests/e2e"
	}
	if len(c.Viewports) == 0 {
		c.Viewports = []string{"desktop-full"}
	}
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.App) == "" {
		return fmt.Errorf("app is required")
	}
	if strings.TrimSpace(c.BaseURL) == "" {
		return fmt.Errorf("base_url is required")
	}
	if strings.TrimSpace(c.RunPackage) == "" {
		return fmt.Errorf("run_package is required")
	}
	if strings.TrimSpace(c.ArtifactsDir) == "" {
		return fmt.Errorf("artifacts_dir is required")
	}
	return nil
}

func (c Config) ResolvePath(path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}
	root := c.RootDir
	if root == "" {
		root = "."
	}
	return filepath.Clean(filepath.Join(root, path))
}
