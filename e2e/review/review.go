package review

import (
	"context"
	"fmt"
)

type ReviewInput struct {
	RunPath      string
	BaselineRun  string
	BaselineRef  string
	WorkspaceRun string
	WorkspaceRef string
	PlanDir      string
}

type ReviewOutcome string

const (
	ReviewPass             ReviewOutcome = "pass"
	ReviewVisualDiff       ReviewOutcome = "visual-diff"
	ReviewMissingBaseline  ReviewOutcome = "missing-baseline"
	ReviewNeedsHumanReview ReviewOutcome = "needs-human-review"
	ReviewError            ReviewOutcome = "error"
)

type ScreenshotComparison struct {
	Story     string `json:"story,omitempty"`
	Scenario  string `json:"scenario,omitempty"`
	Viewport  string `json:"viewport,omitempty"`
	Outcome   string `json:"outcome"`
	Rationale string `json:"rationale,omitempty"`
}

type ReviewResult struct {
	Outcome      ReviewOutcome          `json:"outcome"`
	Summary      string                 `json:"summary"`
	IndexPath    string                 `json:"indexPath,omitempty"`
	MarkdownPath string                 `json:"markdownPath,omitempty"`
	JSONPath     string                 `json:"jsonPath,omitempty"`
	Comparisons  []ScreenshotComparison `json:"comparisons,omitempty"`
}

func Run(ctx context.Context, input ReviewInput) (ReviewResult, error) {
	if input.RunPath == "" {
		return ReviewResult{Outcome: ReviewError}, fmt.Errorf("--run is required")
	}
	result := ReviewResult{
		Outcome: ReviewNeedsHumanReview,
		Summary: "No baseline comparator configured; inspect run artifacts and approve or compare with a baseline run.",
		Comparisons: []ScreenshotComparison{{
			Outcome:   string(ReviewNeedsHumanReview),
			Rationale: "generic DatastarUI review recorded inputs; visual diff engine not configured",
		}},
	}
	if input.PlanDir != "" {
		return WritePlanArtifactsForRun(input.PlanDir, input.RunPath, result)
	}
	return result, nil
}
