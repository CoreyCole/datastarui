package copy

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestAddCopiesRewritesLocksAndDoctorClean(t *testing.T) {
	t.Parallel()

	source := makeSource(t)
	target := t.TempDir()

	result, err := Add(context.Background(), Options{
		SourceRoot:   source,
		TargetRoot:   target,
		TargetModule: "example.com/app",
		Components:   []string{"button"},
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if !slices.Contains(result.CopiedFiles, "components/button/button.templ") {
		t.Fatalf("CopiedFiles = %v, want button templ", result.CopiedFiles)
	}
	copied := readFile(t, filepath.Join(target, "components/button/button.templ"))
	if !strings.Contains(copied, "example.com/app/pkg/datastarui/utils") {
		t.Fatalf("copied import not rewritten: %s", copied)
	}
	if _, err := os.Stat(filepath.Join(target, "datastarui.lock.json")); err != nil {
		t.Fatalf("lock missing: %v", err)
	}
	if err := Doctor(context.Background(), Options{TargetRoot: target}); err != nil {
		t.Fatalf("Doctor() error = %v", err)
	}
}

func TestDiffReportsModifiedAndUpdateRefusesDrift(t *testing.T) {
	t.Parallel()

	source := makeSource(t)
	target := t.TempDir()
	opts := Options{SourceRoot: source, TargetRoot: target, TargetModule: "example.com/app", Components: []string{"button"}}
	if _, err := Add(context.Background(), opts); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	buttonPath := filepath.Join(target, "components/button/button.templ")
	if err := os.WriteFile(buttonPath, []byte(readFile(t, buttonPath)+"\n// local edit\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := Diff(context.Background(), Options{TargetRoot: target})
	if err != nil {
		t.Fatalf("Diff() error = %v", err)
	}
	if !hasDrift(result.Drift, "components/button/button.templ", "modified") {
		t.Fatalf("Diff() drift = %v, want modified button", result.Drift)
	}
	_, err = Update(context.Background(), opts)
	if err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("Update() error = %v, want refusing to overwrite", err)
	}
}

func TestDiffReportsUnmanagedFiles(t *testing.T) {
	t.Parallel()

	source := makeSource(t)
	target := t.TempDir()
	if _, err := Add(context.Background(), Options{SourceRoot: source, TargetRoot: target, TargetModule: "example.com/app", Components: []string{"button"}}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(target, "components/button/local.go"), []byte("package button\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := Diff(context.Background(), Options{TargetRoot: target})
	if err != nil {
		t.Fatalf("Diff() error = %v", err)
	}
	if !hasDrift(result.Drift, "components/button/local.go", "unmanaged") {
		t.Fatalf("Diff() drift = %v, want unmanaged local.go", result.Drift)
	}
}

func makeSource(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "components/button/args.go"), "package button\n")
	writeFile(t, filepath.Join(root, "components/button/button.templ"), "package button\n\nimport \"github.com/coreycole/datastarui/utils\"\n\ntempl Button() { <button></button> }\n")
	writeFile(t, filepath.Join(root, "components/button/variants.go"), "package button\n")
	return root
}

func writeFile(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func hasDrift(entries []DriftEntry, path, status string) bool {
	return slices.ContainsFunc(entries, func(entry DriftEntry) bool {
		return entry.Path == path && entry.Status == status
	})
}
