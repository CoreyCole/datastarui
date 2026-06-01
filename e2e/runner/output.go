package runner

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreycole/datastarui/e2e/artifacts"
)

type RunSummary struct {
	ID           string       `json:"id"`
	Status       string       `json:"status"`
	RunDir       string       `json:"runDir"`
	ManifestPath string       `json:"manifestPath"`
	SummaryPath  string       `json:"summaryPath"`
	IndexPath    string       `json:"indexPath"`
	ServerLog    string       `json:"serverLog,omitempty"`
	JobCount     int          `json:"jobCount"`
	PassedCount  int          `json:"passedCount"`
	FailedCount  int          `json:"failedCount"`
	Jobs         []JobSummary `json:"jobs"`
}

type JobSummary struct {
	ID        string `json:"id"`
	Package   string `json:"package"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
	Artifacts string `json:"artifactsDir,omitempty"`
}

func NewRunID(now time.Time, random io.Reader) (string, error) {
	var suffix [4]byte
	if _, err := io.ReadFull(random, suffix[:]); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s", now.UTC().Format("20060102T150405.000000000Z"), hex.EncodeToString(suffix[:])), nil
}

func WriteRunOutputs(ctx context.Context, plan RunPlan, results []JobResult) (artifacts.RunManifest, error) {
	select {
	case <-ctx.Done():
		return artifacts.RunManifest{}, ctx.Err()
	default:
	}

	entries, err := DiscoverJobArtifacts(plan.Dir, results)
	if err != nil {
		return artifacts.RunManifest{}, err
	}
	manifest := BuildManifest(plan, results, entries)
	if err := WriteManifest(filepath.Join(plan.Dir, "manifest.json"), manifest); err != nil {
		return artifacts.RunManifest{}, err
	}
	if err := WriteSummary(filepath.Join(plan.Dir, "summary.json"), plan, results); err != nil {
		return artifacts.RunManifest{}, err
	}
	if err := WriteIndex(plan.Dir, manifest); err != nil {
		return artifacts.RunManifest{}, err
	}
	return manifest, nil
}

func BuildManifest(plan RunPlan, results []JobResult, entries []artifacts.ArtifactEntry) artifacts.RunManifest {
	manifest := artifacts.RunManifest{
		ID:           plan.ID,
		App:          plan.Config.App,
		StartedAt:    plan.StartedAt,
		CompletedAt:  time.Now(),
		BaseURL:      plan.Server.BaseURL,
		ConfigPath:   plan.Config.ConfigPath,
		ArtifactsDir: plan.Dir,
		BaseRef:      plan.BaseRef,
		ServerMode:   plan.Server.Mode,
		ServerLog:    relPath(plan.Dir, plan.Server.LogPath),
		ChangedFiles: changedFilePaths(plan.ChangedFiles),
		Jobs:         buildArtifactJobSummaries(plan.Dir, results),
		Artifacts:    entries,
	}
	for _, entry := range entries {
		absPath := entry.Path
		if !filepath.IsAbs(absPath) {
			absPath = filepath.Join(plan.Dir, filepath.FromSlash(entry.Path))
		}
		switch entry.Kind {
		case artifacts.ArtifactKindScreenshot:
			manifest.Screenshots = append(manifest.Screenshots, absPath)
		case artifacts.ArtifactKindHTML:
			manifest.HTMLSnapshots = append(manifest.HTMLSnapshots, absPath)
		case artifacts.ArtifactKindTrace:
			manifest.Traces = append(manifest.Traces, absPath)
		}
	}
	return manifest
}

func DiscoverJobArtifacts(runDir string, results []JobResult) ([]artifacts.ArtifactEntry, error) {
	var entries []artifacts.ArtifactEntry
	jobsDir := filepath.Join(runDir, "jobs")
	if _, err := os.Stat(jobsDir); os.IsNotExist(err) {
		return entries, nil
	} else if err != nil {
		return nil, err
	}
	err := filepath.WalkDir(jobsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		kind, ok := artifactKind(path)
		if !ok {
			return nil
		}
		rel, err := filepath.Rel(runDir, path)
		if err != nil {
			return err
		}
		entry := parseArtifactPath(filepath.ToSlash(rel), kind)
		entries = append(entries, entry)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func WriteSummary(path string, plan RunPlan, results []JobResult) error {
	summary := RunSummary{
		ID:           plan.ID,
		Status:       "passed",
		RunDir:       plan.Dir,
		ManifestPath: filepath.Join(plan.Dir, "manifest.json"),
		SummaryPath:  path,
		IndexPath:    filepath.Join(plan.Dir, "index.html"),
		ServerLog:    plan.Server.LogPath,
		JobCount:     len(results),
		Jobs:         make([]JobSummary, 0, len(results)),
	}
	for _, result := range results {
		if result.Status == "passed" {
			summary.PassedCount++
		} else {
			summary.FailedCount++
			summary.Status = "failed"
		}
		summary.Jobs = append(summary.Jobs, JobSummary{
			ID:        result.Job.ID,
			Package:   result.Job.Package,
			Status:    result.Status,
			Error:     result.Error,
			Artifacts: result.Job.ArtifactsDir,
		})
	}
	return writeJSON(path, summary)
}

func WriteManifest(path string, manifest artifacts.RunManifest) error {
	return writeJSON(path, manifest)
}

func WriteIndex(runDir string, manifest artifacts.RunManifest) error {
	_, err := artifacts.WriteStaticIndex(manifest, runDir, artifacts.StaticIndexOptions{})
	return err
}

func writeJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func buildArtifactJobSummaries(runDir string, results []JobResult) []artifacts.JobSummary {
	summaries := make([]artifacts.JobSummary, 0, len(results))
	for _, result := range results {
		summaries = append(summaries, artifacts.JobSummary{
			ID:         result.Job.ID,
			Package:    result.Job.Package,
			RunPattern: result.Job.RunPattern,
			Component:  result.Job.Component,
			Status:     result.Status,
			Duration:   result.Duration.String(),
			StdoutPath: relPath(runDir, result.StdoutPath),
			StderrPath: relPath(runDir, result.StderrPath),
			Error:      result.Error,
		})
	}
	return summaries
}

func changedFilePaths(files []ChangedFile) []string {
	paths := make([]string, 0, len(files))
	for _, file := range files {
		if file.OldPath != "" {
			paths = append(paths, file.OldPath+" -> "+file.Path)
			continue
		}
		paths = append(paths, file.Path)
	}
	return paths
}

func artifactKind(path string) (artifacts.ArtifactKind, bool) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
		return artifacts.ArtifactKindScreenshot, true
	case ".html":
		return artifacts.ArtifactKindHTML, true
	case ".zip", ".trace":
		return artifacts.ArtifactKindTrace, true
	default:
		return "", false
	}
}

func parseArtifactPath(rel string, kind artifacts.ArtifactKind) artifacts.ArtifactEntry {
	entry := artifacts.ArtifactEntry{Kind: kind, Path: rel}
	parts := strings.Split(rel, "/")
	if len(parts) >= 6 && parts[0] == "jobs" {
		entry.FeatureSlug = parts[len(parts)-4]
		entry.ScenarioSlug = parts[len(parts)-3]
		entry.Viewport = parts[len(parts)-2]
		entry.Label = strings.TrimSuffix(parts[len(parts)-1], filepath.Ext(parts[len(parts)-1]))
	}
	return entry
}

func relPath(base, path string) string {
	if path == "" {
		return ""
	}
	rel, err := filepath.Rel(base, path)
	if err != nil || strings.HasPrefix(rel, "..") {
		return path
	}
	return filepath.ToSlash(rel)
}
