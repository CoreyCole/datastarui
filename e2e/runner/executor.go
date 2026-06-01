package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/coreycole/datastarui/e2e/appconfig"
)

type CommandRunner interface {
	Run(ctx context.Context, dir string, env []string, stdoutPath, stderrPath string, args ...string) error
}

type GoCommandRunner struct{}

func (GoCommandRunner) Run(ctx context.Context, dir string, env []string, stdoutPath, stderrPath string, args ...string) error {
	if err := os.MkdirAll(filepath.Dir(stdoutPath), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(stderrPath), 0o755); err != nil {
		return err
	}
	stdout, err := os.Create(stdoutPath)
	if err != nil {
		return err
	}
	defer stdout.Close()
	stderr, err := os.Create(stderrPath)
	if err != nil {
		return err
	}
	defer stderr.Close()

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = io.MultiWriter(os.Stdout, stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderr)
	return cmd.Run()
}

func RunJobs(ctx context.Context, plan RunPlan, workers int, commandRunner CommandRunner) ([]JobResult, error) {
	if commandRunner == nil {
		commandRunner = GoCommandRunner{}
	}
	jobs := AssignJobArtifactDirs(filepath.Join(plan.Dir, "jobs"), plan.Jobs)
	if workers <= 0 {
		workers = DefaultJobs()
	}
	if workers > len(jobs) {
		workers = len(jobs)
	}
	if workers < 1 {
		workers = 1
	}

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	type workItem struct {
		index int
		job   E2EJob
	}
	work := make(chan workItem)
	results := make([]JobResult, len(jobs))
	var errs []error
	var mu sync.Mutex
	var wg sync.WaitGroup

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range work {
				result := runOneJob(runCtx, plan, item.job, commandRunner)
				mu.Lock()
				results[item.index] = result
				if result.Error != "" {
					errs = append(errs, fmt.Errorf("%s: %s", item.job.ID, result.Error))
					cancel()
				}
				mu.Unlock()
			}
		}()
	}

sendLoop:
	for i, job := range jobs {
		select {
		case <-runCtx.Done():
			break sendLoop
		case work <- workItem{index: i, job: job}:
		}
	}
	close(work)
	wg.Wait()

	completed := make([]JobResult, 0, len(results))
	for _, result := range results {
		if result.Job.ID != "" {
			completed = append(completed, result)
		}
	}
	return completed, errors.Join(errs...)
}

func runOneJob(ctx context.Context, plan RunPlan, job E2EJob, commandRunner CommandRunner) JobResult {
	started := time.Now()
	stdoutPath := filepath.Join(job.ArtifactsDir, "stdout.log")
	stderrPath := filepath.Join(job.ArtifactsDir, "stderr.log")
	result := JobResult{
		Job:        job,
		Status:     "passed",
		StdoutPath: stdoutPath,
		StderrPath: stderrPath,
	}
	args := BuildGoTestArgs(job)
	err := commandRunner.Run(ctx, plan.Config.RootDir, JobEnv(plan.Config, job, plan.Server.BaseURL), stdoutPath, stderrPath, args...)
	result.Duration = time.Since(started)
	if err != nil {
		result.Status = "failed"
		result.Error = err.Error()
	}
	return result
}

func BuildGoTestArgs(job E2EJob) []string {
	args := []string{"test", job.Package}
	if strings.TrimSpace(job.RunPattern) != "" {
		args = append(args, "-run", job.RunPattern)
	}
	return args
}

func JobEnv(cfg appconfig.Config, job E2EJob, baseURL string) []string {
	return []string{
		"E2E_CONFIG=" + cfg.ConfigPath,
		"E2E_BASE_URL=" + strings.TrimRight(baseURL, "/"),
		"E2E_ARTIFACTS_DIR=" + job.ArtifactsDir,
		"E2E_CAPTURE_SUCCESS=1",
		"E2E_VIEWPORTS=" + strings.Join(cfg.Viewports, ","),
		"E2E_RUN_BROWSER=1",
	}
}

func AssignJobArtifactDirs(runDir string, jobs []E2EJob) []E2EJob {
	assigned := make([]E2EJob, len(jobs))
	seen := map[string]int{}
	for i, job := range jobs {
		assigned[i] = job
		id := SafeJobID(job)
		seen[id]++
		if seen[id] > 1 {
			id = fmt.Sprintf("%s-%d", id, seen[id])
		}
		assigned[i].ArtifactsDir = filepath.Join(runDir, id)
	}
	return assigned
}

func SafeJobID(job E2EJob) string {
	id := strings.TrimSpace(job.ID)
	if id == "" {
		id = strings.Trim(strings.ReplaceAll(job.Package, "/", "-"), ".-")
	}
	id = strings.ToLower(id)
	var b strings.Builder
	lastDash := false
	for _, r := range id {
		ok := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if ok {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	cleaned := strings.Trim(b.String(), "-")
	if cleaned == "" {
		return "job"
	}
	return cleaned
}

func DefaultJobs() int {
	jobs := runtime.NumCPU() / 2
	if jobs < 1 {
		return 1
	}
	if jobs > 4 {
		return 4
	}
	return jobs
}
