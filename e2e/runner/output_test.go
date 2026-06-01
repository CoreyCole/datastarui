package runner

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/coreycole/datastarui/e2e/appconfig"
	"github.com/coreycole/datastarui/e2e/artifacts"
	"github.com/coreycole/datastarui/e2e/goldens"
)

func TestNewRunIDIncludesNanosecondsAndRandomSuffix(t *testing.T) {
	id, err := NewRunID(time.Date(2026, 6, 1, 10, 2, 3, 456789123, time.FixedZone("PDT", -7*60*60)), bytes.NewReader([]byte{0xde, 0xad, 0xbe, 0xef}))
	if err != nil {
		t.Fatal(err)
	}
	want := "20260601T170203.456789123Z-deadbeef"
	if id != want {
		t.Fatalf("id = %q, want %q", id, want)
	}
}

func TestDiscoverJobArtifactsFindsScreenshotsAndHTML(t *testing.T) {
	runDir := t.TempDir()
	writeFile(t, filepath.Join(runDir, "jobs/select/select-component/select-component-opens-options/desktop-full/page.png"), "png")
	writeFile(t, filepath.Join(runDir, "jobs/select/select-component/select-component-opens-options/desktop-full/page.html"), "html")
	writeFile(t, filepath.Join(runDir, "jobs/select/stdout.log"), "log")

	entries, err := DiscoverJobArtifacts(runDir, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []artifacts.ArtifactEntry{
		{FeatureSlug: "select-component", ScenarioSlug: "select-component-opens-options", Viewport: "desktop-full", Label: "page", Kind: artifacts.ArtifactKindHTML, Path: "jobs/select/select-component/select-component-opens-options/desktop-full/page.html"},
		{FeatureSlug: "select-component", ScenarioSlug: "select-component-opens-options", Viewport: "desktop-full", Label: "page", Kind: artifacts.ArtifactKindScreenshot, Path: "jobs/select/select-component/select-component-opens-options/desktop-full/page.png"},
	}
	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("entries = %#v, want %#v", entries, want)
	}
}

func TestWriteRunOutputsWritesManifestSummaryIndex(t *testing.T) {
	runDir := t.TempDir()
	shot := filepath.Join(runDir, "jobs/select/select-component/select-component-opens-options/desktop-full/page.png")
	page := filepath.Join(runDir, "jobs/select/select-component/select-component-opens-options/desktop-full/page.html")
	writeFile(t, shot, "png")
	writeFile(t, page, "html")

	plan := RunPlan{
		ID:           "run-1",
		Dir:          runDir,
		Config:       appconfig.Config{App: "datastarui", ConfigPath: "datastarui-e2e.yml"},
		BaseRef:      "main",
		ChangedFiles: []ChangedFile{{Status: "M", Path: "components/select/select.templ"}},
		Server:       ServerPlan{Mode: "managed", BaseURL: "http://127.0.0.1:1234", LogPath: filepath.Join(runDir, "server.log")},
		StartedAt:    time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
	}
	results := []JobResult{{
		Job:        E2EJob{ID: "select", Package: "./components/select", Component: "select", ArtifactsDir: filepath.Join(runDir, "jobs/select")},
		Status:     "passed",
		Duration:   time.Second,
		StdoutPath: filepath.Join(runDir, "jobs/select/stdout.log"),
		StderrPath: filepath.Join(runDir, "jobs/select/stderr.log"),
	}}

	manifest, err := WriteRunOutputs(t.Context(), plan, results)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"manifest.json", "summary.json", "index.html"} {
		if _, err := os.Stat(filepath.Join(runDir, name)); err != nil {
			t.Fatalf("%s missing: %v", name, err)
		}
	}
	loaded, err := goldens.LoadManifest(runDir)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.ID != manifest.ID || len(loaded.Screenshots) != 1 || loaded.Screenshots[0] != shot {
		t.Fatalf("loaded manifest = %#v", loaded)
	}
	data, err := os.ReadFile(filepath.Join(runDir, "summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	var summary RunSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		t.Fatal(err)
	}
	if summary.Status != "passed" || summary.PassedCount != 1 || summary.IndexPath == "" {
		t.Fatalf("summary = %#v", summary)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
