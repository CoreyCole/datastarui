package artifacts

import (
	"os"
	"strings"
	"testing"
)

func TestWriteStaticIndexIsGenericDatastarUI(t *testing.T) {
	dir := t.TempDir()
	path, err := WriteStaticIndex(RunManifest{ID: "run-1", App: "demo", BaseURL: "http://localhost:4242"}, dir, StaticIndexOptions{})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	body := string(data)
	if !strings.Contains(body, "DatastarUI E2E run run-1") {
		t.Fatalf("index missing generic title: %s", body)
	}
	if strings.Contains(body, "Vamos") {
		t.Fatalf("index should not contain Vamos branding: %s", body)
	}
}
