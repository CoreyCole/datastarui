package selectors

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/coreycole/datastarui/e2e/appconfig"
)

func TestLoadCatalogFromMapYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "selectors.yaml")
	if err := os.WriteFile(path, []byte("select.trigger: \"[data-slot='select-trigger']\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	catalog, err := LoadCatalog(path)
	if err != nil {
		t.Fatal(err)
	}
	entry, err := catalog.Resolve("select.trigger")
	if err != nil {
		t.Fatal(err)
	}
	if entry.CSS != "[data-slot='select-trigger']" {
		t.Fatalf("CSS = %q", entry.CSS)
	}
}

func TestLoadCatalogFromEntryListYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "selectors.yaml")
	data := []byte("- key: select.content\n  css: \"[data-slot='select-content']\"\n  description: Select content\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	catalog, err := LoadCatalog(path)
	if err != nil {
		t.Fatal(err)
	}
	entry, err := catalog.Resolve("select.content")
	if err != nil {
		t.Fatal(err)
	}
	if entry.Description != "Select content" {
		t.Fatalf("Description = %q", entry.Description)
	}
}

func TestLoadCatalogFromConfigMergesInlineAndFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "selectors.yaml")
	if err := os.WriteFile(path, []byte("shared: .from-file\nfile.only: .file\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := appconfig.Config{
		RootDir:      root,
		Selectors:    map[string]string{"inline.only": ".inline", "shared": ".inline-shared"},
		SelectorFile: "selectors.yaml",
	}

	catalog, err := LoadCatalogFromConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if got, _ := catalog.Resolve("inline.only"); got.CSS != ".inline" {
		t.Fatalf("inline CSS = %q", got.CSS)
	}
	if got, _ := catalog.Resolve("file.only"); got.CSS != ".file" {
		t.Fatalf("file CSS = %q", got.CSS)
	}
	if got, _ := catalog.Resolve("shared"); got.CSS != ".from-file" {
		t.Fatalf("file should override inline, got %q", got.CSS)
	}
}

func TestResolveUnknownKeyReturnsError(t *testing.T) {
	catalog := NewCatalog(nil)
	if _, err := catalog.Resolve("missing"); err == nil {
		t.Fatal("expected unknown selector error")
	}
}
