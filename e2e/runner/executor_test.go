package runner

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coreycole/datastarui/e2e/appconfig"
)

type fakeCommandRunner struct {
	delay     time.Duration
	failJobID string

	mu        sync.Mutex
	inFlight  int
	maxFlight int
	calls     []fakeRunCall
}

type fakeRunCall struct {
	dir        string
	env        []string
	stdoutPath string
	stderrPath string
	args       []string
}

func (f *fakeCommandRunner) Run(ctx context.Context, dir string, env []string, stdoutPath, stderrPath string, args ...string) error {
	f.mu.Lock()
	f.inFlight++
	if f.inFlight > f.maxFlight {
		f.maxFlight = f.inFlight
	}
	f.calls = append(f.calls, fakeRunCall{dir: dir, env: env, stdoutPath: stdoutPath, stderrPath: stderrPath, args: append([]string(nil), args...)})
	thisCall := len(f.calls) - 1
	f.mu.Unlock()

	if f.delay > 0 {
		select {
		case <-ctx.Done():
			f.finish()
			return ctx.Err()
		case <-time.After(f.delay):
		}
	}

	f.finish()
	if f.failJobID != "" && strings.Contains(stdoutPath, f.failJobID) {
		return errors.New("boom")
	}
	_ = thisCall
	return nil
}

func (f *fakeCommandRunner) finish() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.inFlight--
}

func TestRunJobsUsesBoundedConcurrency(t *testing.T) {
	fake := &fakeCommandRunner{delay: 20 * time.Millisecond}
	plan := RunPlan{
		Dir:    t.TempDir(),
		Config: appconfig.Config{RootDir: ".", Viewports: []string{"desktop-full"}},
		Server: ServerPlan{BaseURL: "http://127.0.0.1:4242"},
		Jobs: []E2EJob{
			{ID: "one", Package: "./components/one"},
			{ID: "two", Package: "./components/two"},
			{ID: "three", Package: "./components/three"},
		},
	}
	results, err := RunJobs(context.Background(), plan, 2, fake)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("results = %d, want 3", len(results))
	}
	if fake.maxFlight != 2 {
		t.Fatalf("max in-flight = %d, want 2", fake.maxFlight)
	}
}

func TestRunJobsWritesUniqueArtifactDirs(t *testing.T) {
	fake := &fakeCommandRunner{}
	runDir := t.TempDir()
	plan := RunPlan{
		Dir:    runDir,
		Config: appconfig.Config{RootDir: ".", ConfigPath: "datastarui-e2e.yml", Viewports: []string{"desktop-full"}},
		Server: ServerPlan{BaseURL: "http://127.0.0.1:4242"},
		Jobs: []E2EJob{
			{ID: "select", Package: "./components/select"},
			{ID: "select", Package: "./components/select"},
		},
	}
	results, err := RunJobs(context.Background(), plan, 2, fake)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("results = %d, want 2", len(results))
	}
	if results[0].Job.ArtifactsDir == results[1].Job.ArtifactsDir {
		t.Fatalf("artifact dirs not unique: %q", results[0].Job.ArtifactsDir)
	}
	wantFirst := filepath.Join(runDir, "jobs", "select")
	wantSecond := filepath.Join(runDir, "jobs", "select-2")
	if results[0].Job.ArtifactsDir != wantFirst || results[1].Job.ArtifactsDir != wantSecond {
		t.Fatalf("artifact dirs = %q, %q; want %q, %q", results[0].Job.ArtifactsDir, results[1].Job.ArtifactsDir, wantFirst, wantSecond)
	}
	if !envHas(fake.calls[0].env, "E2E_ARTIFACTS_DIR="+wantFirst) || !envHas(fake.calls[1].env, "E2E_ARTIFACTS_DIR="+wantSecond) {
		t.Fatalf("envs missing per-job artifact dirs: %#v", fake.calls)
	}
}

func TestRunJobsRecordsFailureAndCancelsPending(t *testing.T) {
	fake := &fakeCommandRunner{failJobID: "bad"}
	plan := RunPlan{
		Dir:    t.TempDir(),
		Config: appconfig.Config{RootDir: ".", Viewports: []string{"desktop-full"}},
		Server: ServerPlan{BaseURL: "http://127.0.0.1:4242"},
		Jobs: []E2EJob{
			{ID: "bad", Package: "./components/bad"},
			{ID: "later", Package: "./components/later"},
		},
	}
	results, err := RunJobs(context.Background(), plan, 1, fake)
	if err == nil {
		t.Fatal("expected aggregate error")
	}
	var failed bool
	for _, result := range results {
		if result.Job.ID == "bad" && result.Status == "failed" && result.Error != "" {
			failed = true
		}
	}
	if !failed {
		t.Fatalf("missing failed job result: %#v", results)
	}
}

func TestBuildGoTestArgs(t *testing.T) {
	tests := []struct {
		name string
		job  E2EJob
		want []string
	}{
		{name: "package only", job: E2EJob{Package: "./components/select"}, want: []string{"test", "./components/select"}},
		{name: "run pattern", job: E2EJob{Package: "./components/select", RunPattern: "SelectComponent"}, want: []string{"test", "./components/select", "-run", "SelectComponent"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildGoTestArgs(tt.job); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("args = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestDefaultJobs(t *testing.T) {
	got := DefaultJobs()
	if got < 1 || got > 4 {
		t.Fatalf("DefaultJobs = %d, want 1..4", got)
	}
}

func envHas(env []string, want string) bool {
	for _, value := range env {
		if value == want {
			return true
		}
	}
	return false
}
