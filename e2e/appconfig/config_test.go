package appconfig

import (
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestLoadParsesManagedServerConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultConfigFile)
	data := []byte(`app: datastarui
base_url: http://localhost:4242
run_package: ./components/...
artifacts_dir: .e2e-runs
server:
  command: just build-local
  managed_command: ./datastarui
  skip_when_base_url_set: true
  readiness_path: /components/select
  readiness_timeout: 45s
  port_env: PORT
viewports:
  - desktop-full
`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path, dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.ManagedCommand != "./datastarui" || cfg.Server.PortEnv != "PORT" {
		t.Fatalf("server config = %#v", cfg.Server)
	}
	if got := cfg.ReadinessURL("http://127.0.0.1:9999"); got != "http://127.0.0.1:9999/components/select" {
		t.Fatalf("readiness URL = %q", got)
	}
	if got := cfg.ReadinessTimeout(); got != 45*time.Second {
		t.Fatalf("timeout = %s", got)
	}
}

func TestReadinessHelpersDefaultAndNormalize(t *testing.T) {
	cfg := Config{Server: ServerConfig{ReadinessPath: "components/select", ReadinessTimeout: "bogus"}}
	if got := cfg.ReadinessURL("http://127.0.0.1:9999/"); got != "http://127.0.0.1:9999/components/select" {
		t.Fatalf("readiness URL = %q", got)
	}
	if got := cfg.ReadinessTimeout(); got != 30*time.Second {
		t.Fatalf("default timeout = %s", got)
	}
}
