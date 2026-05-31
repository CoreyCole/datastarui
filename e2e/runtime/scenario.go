package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
)

const scenarioTimeout = 5 * time.Minute

type testRunner interface {
	Run(name string, fn func(t *testing.T)) bool
}

func RunScenario(t *testing.T, featureSlug, scenarioSlug string, fn ScenarioFunc) {
	t.Helper()
	if err := RunScenarioWithConfig(t, nil, featureSlug, scenarioSlug, "", fn); err != nil {
		t.Fatal(err)
	}
}

func RunScenarioWithViewport(
	t *testing.T,
	featureSlug, scenarioSlug string,
	viewport ViewportClass,
	fn ScenarioFunc,
) {
	t.Helper()
	if err := RunScenarioWithConfig(t, nil, featureSlug, scenarioSlug, viewport, fn); err != nil {
		t.Fatal(err)
	}
}

func RunScenarioWithConfig(
	t testing.TB,
	cfg *Config,
	featureSlug, scenarioSlug string,
	viewport ViewportClass,
	fn ScenarioFunc,
) error {
	t.Helper()
	resolved, err := resolveScenarioConfig(t, cfg)
	if err != nil {
		return err
	}
	viewports, err := scenarioViewports(resolved, viewport)
	if err != nil {
		return err
	}
	for _, vp := range viewports {
		vp := vp
		if runner, ok := t.(testRunner); ok {
			runner.Run(string(vp), func(t *testing.T) {
				t.Helper()
				if err := runScenarioOnce(t, resolved, featureSlug, scenarioSlug, vp, fn); err != nil {
					t.Fatal(err)
				}
			})
			continue
		}
		if err := runScenarioOnce(t, resolved, featureSlug, scenarioSlug, vp, fn); err != nil {
			return err
		}
	}
	return nil
}

func resolveScenarioConfig(t testing.TB, cfg *Config) (Config, error) {
	t.Helper()
	if cfg != nil {
		return *cfg, nil
	}
	if os.Getenv("E2E_BASE_URL") == "" && os.Getenv("E2E_RUN_BROWSER") != "1" {
		t.Skip("E2E_BASE_URL is required for browser E2E")
	}
	return LoadConfigFromEnv(".", nil)
}

func scenarioViewports(cfg Config, explicit ViewportClass) ([]ViewportClass, error) {
	if explicit != "" {
		if _, err := ResolveViewports([]string{string(explicit)}); err != nil {
			return nil, err
		}
		return []ViewportClass{explicit}, nil
	}
	if len(cfg.Viewports) == 0 {
		return []ViewportClass{ViewportDesktopFull}, nil
	}
	names := make([]string, 0, len(cfg.Viewports))
	for _, viewport := range cfg.Viewports {
		names = append(names, string(viewport))
	}
	if _, err := ResolveViewports(names); err != nil {
		return nil, err
	}
	return append([]ViewportClass{}, cfg.Viewports...), nil
}

func newPageOptionsForViewport(viewport Viewport) playwright.BrowserNewPageOptions {
	return playwright.BrowserNewPageOptions{
		Viewport: &playwright.Size{
			Width:  viewport.Width,
			Height: viewport.Height,
		},
		ExtraHttpHeaders: map[string]string{
			"X-E2E-Viewport-Class": string(viewport.Class),
		},
	}
}

func runScenarioOnce(
	t testing.TB,
	cfg Config,
	featureSlug, scenarioSlug string,
	viewport ViewportClass,
	fn ScenarioFunc,
) error {
	t.Helper()

	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("start playwright: %w", err)
	}
	defer func() {
		if err := pw.Stop(); err != nil {
			t.Fatalf("stop playwright: %v", err)
		}
	}()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: playwright.Bool(cfg.Headless)})
	if err != nil {
		return fmt.Errorf("launch chromium: %w", err)
	}
	defer func() {
		if err := browser.Close(); err != nil {
			t.Fatalf("close browser: %v", err)
		}
	}()

	resolvedViewports, err := ResolveViewports([]string{string(viewport)})
	if err != nil {
		return err
	}
	page, err := browser.NewPage(newPageOptionsForViewport(resolvedViewports[0]))
	if err != nil {
		return fmt.Errorf("new page: %w", err)
	}

	artifactDir := cfg.ArtifactsDir
	if artifactDir != "" {
		artifactDir = filepath.Join(artifactDir, featureSlug, scenarioSlug, string(viewport))
	}
	artifactSink, err := NewFileArtifactSink(artifactDir)
	if err != nil {
		return fmt.Errorf("artifact sink: %w", err)
	}
	ctx := &Context{
		Config:     cfg,
		Playwright: pw,
		Browser:    browser,
		Page:       page,
		Console:    NewConsoleMonitor(page),
		Artifacts:  artifactSink,
		Memory:     map[string]string{},
	}
	defer func() {
		if t.Failed() || os.Getenv("E2E_CAPTURE_SUCCESS") == "1" {
			_ = ctx.Artifacts.Capture("page", ctx.Page)
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		fn(t, ctx)
	}()
	select {
	case <-done:
	case <-time.After(scenarioTimeout):
		t.Fatalf("scenario %s/%s timed out", featureSlug, scenarioSlug)
	}
	return nil
}
