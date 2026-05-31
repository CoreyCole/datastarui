package review

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWritePlanArtifactsUsesDatastarUIToolName(t *testing.T) {
	result, err := WritePlanArtifacts(t.TempDir(), ReviewResult{Outcome: ReviewNeedsHumanReview, Summary: "inspect"})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(result.MarkdownPath)
	if err != nil {
		t.Fatal(err)
	}
	body := string(data)
	if !strings.Contains(body, "tool: datastarui e2e review") {
		t.Fatalf("markdown missing tool name: %s", body)
	}
	if strings.Contains(body, "QRSPI") || strings.Contains(body, "Vamos") {
		t.Fatalf("markdown should be neutral: %s", body)
	}
	if _, err := os.Stat(result.JSONPath); err != nil {
		t.Fatal(err)
	}
}

func TestWritePlanArtifactsForRunWritesSelfContainedIndex(t *testing.T) {
	runDir := t.TempDir()
	artifactDir := filepath.Join(runDir, "story", "scenario", "desktop-full")
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(artifactDir, "page.png"), []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(artifactDir, "page.html"), []byte("<main>snapshot</main>"), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := WritePlanArtifactsForRun(t.TempDir(), runDir, ReviewResult{Outcome: ReviewNeedsHumanReview, Summary: "inspect"})
	if err != nil {
		t.Fatal(err)
	}
	if result.IndexPath == "" {
		t.Fatal("IndexPath is empty")
	}
	data, err := os.ReadFile(result.IndexPath)
	if err != nil {
		t.Fatal(err)
	}
	body := string(data)
	for _, want := range []string{"DatastarUI E2E visual review", "data:image/png;base64", "story/scenario/desktop-full/page", "&lt;main&gt;snapshot&lt;/main&gt;"} {
		if !strings.Contains(body, want) {
			t.Fatalf("index missing %q: %s", want, body)
		}
	}
}
