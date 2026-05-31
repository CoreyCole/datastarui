package goldens

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coreycole/datastarui/e2e/artifacts"
	"github.com/coreycole/datastarui/e2e/review"
)

type GoldenInput struct {
	RunPath       string
	GoldenRoot    string
	PlanDir       string
	HumanApproved bool
}

func Compare(ctx context.Context, input GoldenInput) (review.ReviewResult, error) {
	if input.RunPath == "" {
		return review.ReviewResult{Outcome: review.ReviewError}, fmt.Errorf("--run is required")
	}
	result := review.ReviewResult{Outcome: review.ReviewNeedsHumanReview, Summary: "golden comparison is not configured; inspect run artifacts"}
	if input.PlanDir != "" {
		return review.WritePlanArtifacts(input.PlanDir, result)
	}
	return result, nil
}

func Accept(ctx context.Context, input GoldenInput) error {
	if !input.HumanApproved {
		return fmt.Errorf("goldens accept requires --human-approved")
	}
	manifest, err := LoadManifest(input.RunPath)
	if err != nil {
		return err
	}
	return copyScreenshots(manifest, input.GoldenRoot)
}

func LoadManifest(runDir string) (artifacts.RunManifest, error) {
	for _, name := range []string{"manifest.json", "run.json"} {
		path := filepath.Join(runDir, name)
		data, err := os.ReadFile(path)
		if err == nil {
			var manifest artifacts.RunManifest
			if err := json.Unmarshal(data, &manifest); err != nil {
				return artifacts.RunManifest{}, fmt.Errorf("read %s: %w", path, err)
			}
			return manifest, nil
		}
		if !os.IsNotExist(err) {
			return artifacts.RunManifest{}, err
		}
	}
	return artifacts.RunManifest{}, fmt.Errorf("run manifest missing in %s", runDir)
}

func copyScreenshots(run artifacts.RunManifest, root string) error {
	if root == "" {
		return fmt.Errorf("golden root is empty")
	}
	for _, src := range run.Screenshots {
		dst := filepath.Join(root, filepath.Base(src))
		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return err
		}
	}
	return nil
}
