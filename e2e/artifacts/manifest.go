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
	BaseRef        string          `json:"baseRef,omitempty"`
	ChangedFiles   []string        `json:"changedFiles,omitempty"`
	ServerMode     string          `json:"serverMode,omitempty"`
	ServerLog      string          `json:"serverLog,omitempty"`
	Jobs           []JobSummary    `json:"jobs,omitempty"`
	Artifacts      []ArtifactEntry `json:"artifacts,omitempty"`
	Screenshots    []string        `json:"screenshots,omitempty"`
	HTMLSnapshots  []string        `json:"htmlSnapshots,omitempty"`
	Traces         []string        `json:"traces,omitempty"`
}

type JobSummary struct {
	ID         string `json:"id"`
	Package    string `json:"package"`
	RunPattern string `json:"runPattern,omitempty"`
	Component  string `json:"component,omitempty"`
	Status     string `json:"status"`
	Duration   string `json:"duration,omitempty"`
	StdoutPath string `json:"stdoutPath,omitempty"`
	StderrPath string `json:"stderrPath,omitempty"`
	Error      string `json:"error,omitempty"`
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
