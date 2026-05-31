package runtime

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/playwright-community/playwright-go"

	"github.com/coreycole/datastarui/e2e/appconfig"
)

type App struct {
	Name         string
	Authenticate func(context.Context, playwright.Page, Config, string) error
	Preflight    func(context.Context, Config) error
}

func (a App) authenticate(ctx context.Context, page playwright.Page, cfg Config, email string) error {
	if a.Authenticate == nil {
		return nil
	}
	return a.Authenticate(ctx, page, cfg, email)
}

func (a App) preflight(ctx context.Context, cfg Config) error {
	if a.Preflight == nil {
		return nil
	}
	return a.Preflight(ctx, cfg)
}

type AppRuntime interface {
	Authenticate(context.Context, playwright.Page, Config, string) error
	Preflight(context.Context, Config) error
}

type noopAppRuntime struct{}

var defaultAppRuntime App = App{Name: "app"}

func NoopAppRuntime() AppRuntime { return noopAppRuntime{} }

func SetDefaultAppRuntime(app AppRuntime) {
	if app == nil {
		defaultAppRuntime = App{Name: "app"}
		return
	}
	defaultAppRuntime = App{
		Authenticate: app.Authenticate,
		Preflight:    app.Preflight,
	}
}

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
	App          App
	AppConfig    appconfig.Config
}

func LoadConfig(cwd string, cfg appconfig.Config, app App) (Config, error) {
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

	if isZeroApp(app) {
		app = defaultAppRuntime
		app.Name = cfg.App
	} else if app.Name == "" {
		app.Name = cfg.App
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
		App:          app,
		AppConfig:    cfg,
	}, nil
}

func LoadConfigFromEnv(cwd string, app App) (Config, error) {
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
	return "", fmt.Errorf("unknown page key %q; use a literal path or app-specific page object", keyOrPath)
}

func isZeroApp(app App) bool {
	return app.Name == "" && app.Authenticate == nil && app.Preflight == nil
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
