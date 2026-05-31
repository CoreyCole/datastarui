package runtime

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/playwright-community/playwright-go"

	"github.com/coreycole/datastarui/e2e/appconfig"
	"github.com/coreycole/datastarui/e2e/selectors"
)

type AppRuntime interface {
	Authenticate(context.Context, playwright.Page, Config, string) error
	Preflight(context.Context, Config) error
}

type noopAppRuntime struct{}

func NoopAppRuntime() AppRuntime { return noopAppRuntime{} }

func (noopAppRuntime) Authenticate(context.Context, playwright.Page, Config, string) error {
	return nil
}
func (noopAppRuntime) Preflight(context.Context, Config) error { return nil }

type Config struct {
	AppName      string
	RepoRoot     string
	BaseURL      string
	ArtifactsDir string
	Headless     bool
	Viewports    []ViewportClass
	Selectors    selectors.Catalog
	Pages        map[string]string
	App          AppRuntime
	AppConfig    appconfig.Config
}

func LoadConfig(cwd string, cfg appconfig.Config, app AppRuntime) (Config, error) {
	catalog, err := selectors.LoadCatalogFromConfig(cfg)
	if err != nil {
		return Config{}, err
	}
	baseURL := envOr("E2E_BASE_URL", cfg.BaseURL)
	artifactsDir := envOr("E2E_ARTIFACTS_DIR", cfg.ResolvePath(cfg.ArtifactsDir))

	viewportNames := cfg.Viewports
	if raw := strings.TrimSpace(os.Getenv("E2E_VIEWPORTS")); raw != "" {
		viewportNames = splitCSV(raw)
	}
	if _, err := ResolveViewports(viewportNames); err != nil {
		return Config{}, err
	}
	viewports := make([]ViewportClass, 0, len(viewportNames))
	for _, name := range viewportNames {
		viewports = append(viewports, ViewportClass(name))
	}

	if app == nil {
		app = NoopAppRuntime()
	}
	root := cfg.RootDir
	if root == "" {
		root = findRepoRoot(cwd)
	}
	return Config{
		AppName:      cfg.App,
		RepoRoot:     root,
		BaseURL:      strings.TrimRight(baseURL, "/"),
		ArtifactsDir: artifactsDir,
		Headless:     os.Getenv("E2E_HEADLESS") != "false",
		Viewports:    viewports,
		Selectors:    catalog,
		Pages:        cfg.Pages,
		App:          app,
		AppConfig:    cfg,
	}, nil
}

func LoadConfigFromEnv(cwd string, app AppRuntime) (Config, error) {
	path, ok, err := appconfig.Find(cwd)
	if err != nil {
		return Config{}, err
	}
	if !ok && os.Getenv("E2E_RUN_BROWSER") != "1" {
		return Config{}, fmt.Errorf("%s not found", appconfig.DefaultConfigFile)
	}
	cfg, err := appconfig.Load(path, cwd)
	if err != nil {
		return Config{}, err
	}
	return LoadConfig(cwd, cfg, app)
}

func (c Config) PagePath(keyOrPath string) (string, error) {
	if strings.HasPrefix(keyOrPath, "/") {
		return keyOrPath, nil
	}
	if path, ok := c.Pages[keyOrPath]; ok {
		return path, nil
	}
	return "", fmt.Errorf("unknown page key %q", keyOrPath)
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func findRepoRoot(cwd string) string {
	cur, err := filepath.Abs(cwd)
	if err != nil {
		return cwd
	}
	for {
		if _, err := os.Stat(filepath.Join(cur, "go.mod")); err == nil {
			return cur
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return cwd
		}
		cur = parent
	}
}
