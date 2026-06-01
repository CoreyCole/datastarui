package main

import (
	"context"
	cryptoRand "crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreycole/datastarui/e2e/appconfig"
	"github.com/coreycole/datastarui/e2e/goldens"
	"github.com/coreycole/datastarui/e2e/review"
	"github.com/coreycole/datastarui/e2e/runner"
	"github.com/coreycole/datastarui/e2e/runtime"
	e2eserver "github.com/coreycole/datastarui/e2e/server"
	"github.com/spf13/cobra"
)

type e2eRunOptions struct {
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

type e2eReviewOptions struct {
	RunPath      string
	BaselineRun  string
	BaselineRef  string
	WorkspaceRun string
	WorkspaceRef string
	PlanDir      string
}

type e2eGoldensOptions struct {
	RunPath       string
	GoldenRoot    string
	PlanDir       string
	HumanApproved bool
}

func newE2ECommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{Use: "e2e", Short: "Run Go Story E2E tests"}
	cmd.AddCommand(newE2ERunCommand(ctx), newE2EReviewCommand(ctx), newE2EGoldensCommand(ctx))
	return cmd
}

func newE2ERunCommand(ctx context.Context) *cobra.Command {
	opts := e2eRunOptions{}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run configured Go Story E2E tests",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if strings.TrimSpace(opts.Story) == "" {
					return fmt.Errorf("unexpected e2e run arguments: %s", strings.Join(args, " "))
				}
				opts.Story = strings.Join(append([]string{opts.Story}, args...), " ")
			}
			if ctx == nil {
				ctx = cmd.Context()
			}
			if ctx == nil {
				ctx = context.Background()
			}
			return runE2E(ctx, opts)
		},
	}
	cmd.Flags().StringVar(&opts.ConfigPath, "config", "", "E2E YAML config path (defaults to datastarui-e2e.yml discovery)")
	cmd.Flags().StringVar(&opts.Story, "story", "", "story slug to run")
	cmd.Flags().StringVar(&opts.Scenario, "scenario", "", "scenario slug to run")
	cmd.Flags().StringVar(&opts.Viewport, "viewport", "", "viewport class or comma-separated viewport classes to run")
	cmd.Flags().StringVar(&opts.BaseURL, "base-url", "", "base URL for browser E2E")
	cmd.Flags().StringVar(&opts.ArtifactsDir, "artifacts-dir", "", "directory for run artifacts")
	cmd.Flags().BoolVar(&opts.NoRestart, "no-restart", false, "skip configured server command before running")
	cmd.Flags().BoolVar(&opts.All, "all", false, "run all configured E2E tests instead of changed-only default")
	cmd.Flags().StringVar(&opts.BaseRef, "base-ref", "main", "git ref for changed-only E2E selection")
	cmd.Flags().IntVar(&opts.Jobs, "jobs", 0, "parallel E2E jobs; defaults conservatively from CPU")
	cmd.Flags().StringVar(&opts.ReadinessPath, "readiness-path", "", "override configured readiness path")
	cmd.Flags().DurationVar(&opts.ReadinessTimeout, "readiness-timeout", 0, "override configured readiness timeout")
	return cmd
}

func newE2EReviewCommand(ctx context.Context) *cobra.Command {
	opts := e2eReviewOptions{}
	cmd := &cobra.Command{
		Use:   "review",
		Short: "Write neutral E2E visual review artifacts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if ctx == nil {
				ctx = cmd.Context()
			}
			result, err := review.Run(ctx, review.ReviewInput{
				RunPath:      opts.RunPath,
				BaselineRun:  opts.BaselineRun,
				BaselineRef:  opts.BaselineRef,
				WorkspaceRun: opts.WorkspaceRun,
				WorkspaceRef: opts.WorkspaceRef,
				PlanDir:      opts.PlanDir,
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "E2E review outcome: %s\n", result.Outcome)
			if result.MarkdownPath != "" {
				fmt.Fprintf(os.Stdout, "E2E review markdown: %s\n", result.MarkdownPath)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.RunPath, "run", "", "E2E run directory or manifest path")
	cmd.Flags().StringVar(&opts.BaselineRun, "baseline-run", "", "baseline run directory or manifest path")
	cmd.Flags().StringVar(&opts.BaselineRef, "baseline-ref", "", "baseline git ref")
	cmd.Flags().StringVar(&opts.WorkspaceRun, "workspace-run", "", "workspace run directory or manifest path")
	cmd.Flags().StringVar(&opts.WorkspaceRef, "workspace-ref", "", "workspace git ref")
	cmd.Flags().StringVar(&opts.PlanDir, "plan-dir", "", "optional plan directory for Markdown/JSON outputs")
	return cmd
}

func newE2EGoldensCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{Use: "goldens", Short: "Compare or accept E2E golden screenshots"}
	cmd.AddCommand(newE2EGoldensCompareCommand(ctx), newE2EGoldensAcceptCommand(ctx))
	return cmd
}

func newE2EGoldensCompareCommand(ctx context.Context) *cobra.Command {
	opts := e2eGoldensOptions{}
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare a run against golden screenshots",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if ctx == nil {
				ctx = cmd.Context()
			}
			result, err := goldens.Compare(ctx, goldens.GoldenInput{RunPath: opts.RunPath, GoldenRoot: opts.GoldenRoot, PlanDir: opts.PlanDir})
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "E2E goldens outcome: %s\n", result.Outcome)
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.RunPath, "run", "", "E2E run directory")
	cmd.Flags().StringVar(&opts.GoldenRoot, "golden-root", "", "golden screenshot root")
	cmd.Flags().StringVar(&opts.PlanDir, "plan-dir", "", "optional plan directory for Markdown/JSON outputs")
	return cmd
}

func newE2EGoldensAcceptCommand(ctx context.Context) *cobra.Command {
	opts := e2eGoldensOptions{}
	cmd := &cobra.Command{
		Use:   "accept",
		Short: "Accept run screenshots as goldens after human approval",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if ctx == nil {
				ctx = cmd.Context()
			}
			return goldens.Accept(ctx, goldens.GoldenInput{RunPath: opts.RunPath, GoldenRoot: opts.GoldenRoot, HumanApproved: opts.HumanApproved})
		},
	}
	cmd.Flags().StringVar(&opts.RunPath, "run", "", "E2E run directory")
	cmd.Flags().StringVar(&opts.GoldenRoot, "golden-root", "", "golden screenshot root")
	cmd.Flags().BoolVar(&opts.HumanApproved, "human-approved", false, "confirm a human approved accepting goldens")
	return cmd
}

func runE2E(ctx context.Context, opts e2eRunOptions) error {
	cfg, err := loadE2EConfig(opts)
	if err != nil {
		return err
	}
	applyE2ERunOverrides(&cfg, opts)

	runID, err := runner.NewRunID(time.Now(), cryptoRand.Reader)
	if err != nil {
		return err
	}
	artifactRoot := cfg.ResolvePath(cfg.ArtifactsDir)
	if opts.ArtifactsDir != "" {
		artifactRoot = opts.ArtifactsDir
	}
	runDir := filepath.Join(artifactRoot, runID)
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		return err
	}
	if err := setE2EEnvironment(cfg, runDir, strings.Join(cfg.Viewports, ",")); err != nil {
		return err
	}

	if command := buildE2EServerCommand(opts, cfg); len(command) > 0 {
		setup := exec.CommandContext(ctx, command[0], command[1:]...)
		setup.Stdout = os.Stdout
		setup.Stderr = os.Stderr
		setup.Dir = cfg.RootDir
		if err := setup.Run(); err != nil {
			return fmt.Errorf("%s: %w", strings.Join(command, " "), err)
		}
	}

	changed, err := changedFilesForRun(ctx, cfg, opts)
	if err != nil {
		return err
	}
	jobs, err := runner.PlanJobs(cfg, toRunnerOptions(opts), changed)
	if err != nil {
		return err
	}
	jobs = runner.AssignJobArtifactDirs(filepath.Join(runDir, "jobs"), jobs)
	for _, job := range jobs {
		if err := ensureE2ETestsExist(ctx, runner.BuildGoTestArgs(job), cfg.RootDir); err != nil {
			return err
		}
	}

	managed := e2eserver.New(cfg, e2eserver.Options{BaseURL: cfg.BaseURL, RunDir: runDir, NoRestart: opts.NoRestart})
	if err := managed.Start(ctx); err != nil {
		return err
	}
	defer func() { _ = managed.Cleanup(context.Background()) }()
	if err := managed.WaitReady(ctx); err != nil {
		return err
	}

	serverMode := "managed"
	if opts.NoRestart || strings.TrimSpace(cfg.Server.ManagedCommand) == "" {
		serverMode = "external"
	}
	plan := runner.RunPlan{
		ID:           runID,
		Dir:          runDir,
		Config:       cfg,
		BaseRef:      opts.BaseRef,
		ChangedFiles: changed,
		Jobs:         jobs,
		StartedAt:    time.Now(),
		Server: runner.ServerPlan{
			Mode:         serverMode,
			BaseURL:      managed.BaseURL(),
			ReadinessURL: cfg.ReadinessURL(managed.BaseURL()),
			LogPath:      managed.LogPath(),
		},
	}
	results, runErr := runner.RunJobs(ctx, plan, opts.Jobs, runner.GoCommandRunner{})
	_, outputErr := runner.WriteRunOutputs(ctx, plan, results)
	printE2ERunArtifacts(runDir, managed.LogPath())
	if outputErr != nil {
		return outputErr
	}
	if runErr != nil {
		return runErr
	}
	return nil
}

func loadE2EConfig(opts e2eRunOptions) (appconfig.Config, error) {
	path := opts.ConfigPath
	if path == "" {
		path = os.Getenv("E2E_CONFIG")
	}
	return appconfig.Load(path, ".")
}

func applyE2ERunOverrides(cfg *appconfig.Config, opts e2eRunOptions) {
	if opts.BaseURL != "" {
		cfg.BaseURL = opts.BaseURL
	}
	if opts.ArtifactsDir != "" {
		cfg.ArtifactsDir = opts.ArtifactsDir
	}
	if opts.ReadinessPath != "" {
		cfg.Server.ReadinessPath = opts.ReadinessPath
	}
	if opts.ReadinessTimeout > 0 {
		cfg.Server.ReadinessTimeout = opts.ReadinessTimeout.String()
	}
	cfg.Viewports = selectedE2EViewports(opts)
}

func changedFilesForRun(ctx context.Context, cfg appconfig.Config, opts e2eRunOptions) ([]runner.ChangedFile, error) {
	if opts.All || strings.TrimSpace(opts.Story) != "" || strings.TrimSpace(opts.Scenario) != "" {
		return nil, nil
	}
	return runner.ChangedFiles(ctx, cfg.RootDir, opts.BaseRef)
}

func toRunnerOptions(opts e2eRunOptions) runner.RunOptions {
	return runner.RunOptions{
		ConfigPath:       opts.ConfigPath,
		Story:            opts.Story,
		Scenario:         opts.Scenario,
		Viewport:         opts.Viewport,
		BaseURL:          opts.BaseURL,
		ArtifactsDir:     opts.ArtifactsDir,
		NoRestart:        opts.NoRestart,
		All:              opts.All,
		BaseRef:          opts.BaseRef,
		Jobs:             opts.Jobs,
		ReadinessPath:    opts.ReadinessPath,
		ReadinessTimeout: opts.ReadinessTimeout,
	}
}

func printE2ERunArtifacts(runDir, serverLog string) {
	fmt.Fprintf(os.Stdout, "E2E run: %s\n", runDir)
	fmt.Fprintf(os.Stdout, "E2E manifest: %s\n", filepath.Join(runDir, "manifest.json"))
	fmt.Fprintf(os.Stdout, "E2E summary: %s\n", filepath.Join(runDir, "summary.json"))
	fmt.Fprintf(os.Stdout, "E2E index: %s\n", filepath.Join(runDir, "index.html"))
	if serverLog != "" {
		fmt.Fprintf(os.Stdout, "E2E server log: %s\n", serverLog)
	}
}

func setE2EEnvironment(cfg appconfig.Config, runDir, viewportEnv string) error {
	values := map[string]string{
		"E2E_CONFIG":          cfg.ConfigPath,
		"E2E_BASE_URL":        strings.TrimRight(cfg.BaseURL, "/"),
		"E2E_ARTIFACTS_DIR":   runDir,
		"E2E_CAPTURE_SUCCESS": "1",
		"E2E_VIEWPORTS":       viewportEnv,
		"E2E_RUN_BROWSER":     "1",
	}
	for key, value := range values {
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	return nil
}

func buildE2EServerCommand(opts e2eRunOptions, cfg appconfig.Config) []string {
	if opts.NoRestart {
		return nil
	}
	if cfg.Server.SkipWhenBaseURLSet && (opts.BaseURL != "" || os.Getenv("E2E_BASE_URL") != "") {
		return nil
	}
	command := strings.TrimSpace(cfg.Server.Command)
	if command == "" {
		return nil
	}
	return strings.Fields(command)
}

func selectedE2EViewportEnv(opts e2eRunOptions) string {
	return strings.Join(selectedE2EViewports(opts), ",")
}

func selectedE2EViewports(opts e2eRunOptions) []string {
	if viewport := strings.TrimSpace(opts.Viewport); viewport != "" {
		parts := strings.Split(viewport, ",")
		out := make([]string, 0, len(parts))
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				out = append(out, trimmed)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	defaults := runtime.DefaultVerifyViewports()
	parts := make([]string, 0, len(defaults))
	for _, viewport := range defaults {
		parts = append(parts, string(viewport))
	}
	return parts
}

func buildE2EGoTestArgs(opts e2eRunOptions, cfg appconfig.Config) []string {
	runPackage := strings.TrimSpace(cfg.RunPackage)
	if runPackage == "" {
		runPackage = "./tests/e2e"
	}
	args := []string{"test", runPackage}
	if opts.Story != "" || opts.Scenario != "" {
		pattern := slugToTestFragment(opts.Story)
		if opts.Scenario != "" {
			pattern += ".*" + slugToTestFragment(opts.Scenario)
		}
		args = append(args, "-run", pattern)
	}
	return args
}

func ensureE2ETestsExist(ctx context.Context, args []string, dir string) error {
	pattern := ""
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-run" {
			pattern = args[i+1]
			break
		}
	}
	if pattern == "" {
		return nil
	}
	if len(args) < 2 {
		return fmt.Errorf("go test package argument missing")
	}
	listArgs := []string{"test", args[1], "-list", pattern}
	cmd := exec.CommandContext(ctx, "go", listArgs...)
	cmd.Dir = repoRootForCommand(dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go %v: %w\n%s", listArgs, err, strings.TrimSpace(string(out)))
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "Test") {
			return nil
		}
	}
	return fmt.Errorf("no E2E tests matched -run %q", pattern)
}

func repoRootForCommand(cwd string) string {
	abs, err := filepath.Abs(cwd)
	if err != nil {
		return cwd
	}
	for {
		if _, err := os.Stat(filepath.Join(abs, "go.mod")); err == nil {
			return abs
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			return cwd
		}
		abs = parent
	}
}

func slugToTestFragment(slug string) string {
	return runner.SlugToTestFragment(slug)
}
