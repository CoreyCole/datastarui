package runtime

import (
	"path/filepath"
	"testing"

	"github.com/coreycole/datastarui/e2e/appconfig"
)

func TestLoadConfigBuildsRuntimeConfig(t *testing.T) {
	root := t.TempDir()
	cfg := appconfig.Config{
		App:          "datastarui",
		BaseURL:      "http://localhost:4242/",
		ArtifactsDir: ".e2e-runs",
		RootDir:      root,
		Selectors:    map[string]string{"select.trigger": "[data-slot='select-trigger']"},
		Pages:        map[string]string{"SelectComponent": "/components/select"},
		Viewports:    []string{"desktop-full"},
	}

	runtimeCfg, err := LoadConfig(root, cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	if runtimeCfg.AppName != "datastarui" {
		t.Fatalf("AppName = %q", runtimeCfg.AppName)
	}
	if runtimeCfg.BaseURL != "http://localhost:4242" {
		t.Fatalf("BaseURL = %q", runtimeCfg.BaseURL)
	}
	if runtimeCfg.ArtifactsDir != filepath.Join(root, ".e2e-runs") {
		t.Fatalf("ArtifactsDir = %q", runtimeCfg.ArtifactsDir)
	}
	if _, err := runtimeCfg.Selectors.Resolve("select.trigger"); err != nil {
		t.Fatal(err)
	}
	if runtimeCfg.App == nil {
		t.Fatal("expected default app runtime")
	}
}

func TestLoadConfigEnvOverridesBaseURLArtifactsAndViewports(t *testing.T) {
	t.Setenv("E2E_BASE_URL", "http://127.0.0.1:4242")
	t.Setenv("E2E_ARTIFACTS_DIR", "/tmp/e2e-artifacts")
	t.Setenv("E2E_VIEWPORTS", "mobile,desktop-half")
	cfg := appconfig.Config{
		App:          "datastarui",
		BaseURL:      "http://localhost:4242",
		ArtifactsDir: ".e2e-runs",
		RootDir:      t.TempDir(),
		Viewports:    []string{"desktop-full"},
	}

	runtimeCfg, err := LoadConfig(".", cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	if runtimeCfg.BaseURL != "http://127.0.0.1:4242" {
		t.Fatalf("BaseURL = %q", runtimeCfg.BaseURL)
	}
	if runtimeCfg.ArtifactsDir != "/tmp/e2e-artifacts" {
		t.Fatalf("ArtifactsDir = %q", runtimeCfg.ArtifactsDir)
	}
	if got := runtimeCfg.Viewports; len(got) != 2 || got[0] != ViewportMobile || got[1] != ViewportDesktopHalf {
		t.Fatalf("Viewports = %#v", got)
	}
}

func TestPagePathResolvesAliasAndLiteralPath(t *testing.T) {
	cfg := Config{Pages: map[string]string{"SelectComponent": "/components/select"}}
	if got, err := cfg.PagePath("SelectComponent"); err != nil || got != "/components/select" {
		t.Fatalf("alias PagePath() = %q, %v", got, err)
	}
	if got, err := cfg.PagePath("/literal"); err != nil || got != "/literal" {
		t.Fatalf("literal PagePath() = %q, %v", got, err)
	}
}

func TestPagePathRejectsUnknownAlias(t *testing.T) {
	cfg := Config{Pages: map[string]string{}}
	if _, err := cfg.PagePath("Missing"); err == nil {
		t.Fatal("expected unknown alias error")
	}
}
