package appconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigFileIsGenericDatastarUIName(t *testing.T) {
	if DefaultConfigFile != "datastarui-e2e.yml" {
		t.Fatalf("DefaultConfigFile = %q", DefaultConfigFile)
	}
}

func TestFindWalksUpToDatastarUIE2EConfig(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(root, DefaultConfigFile)
	if err := os.WriteFile(configPath, []byte("app: test\nbase_url: http://localhost:4242\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	found, ok, err := Find(nested)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected config to be found")
	}
	if found != configPath {
		t.Fatalf("Find() = %q, want %q", found, configPath)
	}
}

func TestLoadAppliesDefaultsAndRootDir(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, DefaultConfigFile)
	if err := os.WriteFile(configPath, []byte("app: datastarui\nbase_url: http://localhost:4242\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath, root)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.RunPackage != "./tests/e2e" {
		t.Fatalf("RunPackage = %q", cfg.RunPackage)
	}
	if cfg.RootDir != root {
		t.Fatalf("RootDir = %q, want %q", cfg.RootDir, root)
	}
	if len(cfg.Viewports) != 1 || cfg.Viewports[0] != "desktop-full" {
		t.Fatalf("Viewports = %#v", cfg.Viewports)
	}
}

func TestLoadRejectsMissingBaseURL(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, DefaultConfigFile)
	if err := os.WriteFile(configPath, []byte("app: datastarui\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath, root)
	if err == nil {
		t.Fatal("expected missing base_url error")
	}
}

func TestResolvePathUsesConfigRoot(t *testing.T) {
	cfg := Config{RootDir: "/tmp/project"}
	if got := cfg.ResolvePath("artifacts"); got != "/tmp/project/artifacts" {
		t.Fatalf("ResolvePath() = %q", got)
	}
	if got := cfg.ResolvePath("/already/absolute"); got != "/already/absolute" {
		t.Fatalf("absolute ResolvePath() = %q", got)
	}
}
