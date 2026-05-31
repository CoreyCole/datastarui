package artifacts

import "time"

type ArtifactKind string

const (
	ArtifactKindScreenshot ArtifactKind = "screenshot"
	ArtifactKindHTML       ArtifactKind = "html"
	ArtifactKindTrace      ArtifactKind = "trace"
)

type ArtifactEntry struct {
	FeatureSlug  string       `json:"featureSlug"`
	ScenarioSlug string       `json:"scenarioSlug"`
	Viewport     string       `json:"viewport"`
	Label        string       `json:"label"`
	Kind         ArtifactKind `json:"kind"`
	Path         string       `json:"path"`
}

type RunManifest struct {
	ID             string          `json:"id"`
	App            string          `json:"app,omitempty"`
	StartedAt      time.Time       `json:"startedAt"`
	CompletedAt    time.Time       `json:"completedAt,omitempty"`
	BaseURL        string          `json:"baseUrl,omitempty"`
	ConfigPath     string          `json:"configPath,omitempty"`
	ArtifactsDir   string          `json:"artifactsDir,omitempty"`
	ViewportFilter string          `json:"viewportFilter,omitempty"`
	Artifacts      []ArtifactEntry `json:"artifacts,omitempty"`
	Screenshots    []string        `json:"screenshots,omitempty"`
	HTMLSnapshots  []string        `json:"htmlSnapshots,omitempty"`
	Traces         []string        `json:"traces,omitempty"`
}

type ScenarioResult struct {
	Feature    string `json:"feature"`
	Scenario   string `json:"scenario"`
	Viewport   string `json:"viewport"`
	Status     string `json:"status"`
	Screenshot string `json:"screenshot,omitempty"`
	HTML       string `json:"html,omitempty"`
	Trace      string `json:"trace,omitempty"`
	Error      string `json:"error,omitempty"`
}
