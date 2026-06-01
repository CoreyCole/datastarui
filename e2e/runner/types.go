package runner

import (
	"time"

	"github.com/coreycole/datastarui/e2e/appconfig"
)

type RunOptions struct {
	ConfigPath       string
	Story            string
	Scenario         string
	Viewport         string
	BaseURL          string
	ArtifactsDir     string
	NoRestart        bool
	All              bool
	BaseRef          string
	Jobs             int
	ReadinessPath    string
	ReadinessTimeout time.Duration
}

type ChangedFile struct {
	Status  string
	Path    string
	OldPath string
}

type RunPlan struct {
	ID           string
	Dir          string
	Config       appconfig.Config
	BaseRef      string
	ChangedFiles []ChangedFile
	Jobs         []E2EJob
	Server       ServerPlan
	StartedAt    time.Time
}

type ServerPlan struct {
	Mode         string
	BaseURL      string
	ReadinessURL string
	LogPath      string
}

type E2EJob struct {
	ID           string
	Package      string
	RunPattern   string
	Component    string
	ArtifactsDir string
	Reason       string
}
