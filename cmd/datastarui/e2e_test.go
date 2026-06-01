package main

import (
	"context"
	"strings"
	"testing"

	"github.com/coreycole/datastarui/e2e/appconfig"
)

func TestRootCommandIncludesE2E(t *testing.T) {
	cmd := newRootCommand(context.Background(), nil, nil)
	for _, child := range cmd.Commands() {
		if child.Name() == "e2e" {
			return
		}
	}
	t.Fatal("missing e2e command")
}

func TestBuildE2EGoTestArgsUsesConfiguredRunPackage(t *testing.T) {
	args := buildE2EGoTestArgs(e2eRunOptions{Story: "select-component"}, appconfig.Config{RunPackage: "./custom/e2e"})
	want := []string{"test", "./custom/e2e", "-run", "SelectComponent"}
	assertStringSliceEqual(t, args, want)
}

func TestBuildE2EGoTestArgsSupportsPackagePattern(t *testing.T) {
	args := buildE2EGoTestArgs(e2eRunOptions{Story: "select-component"}, appconfig.Config{RunPackage: "./components/..."})
	want := []string{"test", "./components/...", "-run", "SelectComponent"}
	assertStringSliceEqual(t, args, want)
}

func assertStringSliceEqual(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("args len = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("args[%d] = %q, want %q (all %#v)", i, got[i], want[i], got)
		}
	}
}

func TestE2ECommandIncludesReviewAndGoldens(t *testing.T) {
	cmd := newE2ECommand(context.Background())
	seen := map[string]bool{}
	for _, child := range cmd.Commands() {
		seen[child.Name()] = true
	}
	for _, name := range []string{"run", "review", "goldens"} {
		if !seen[name] {
			t.Fatalf("missing e2e %s command", name)
		}
	}
}

func TestBuildE2EServerCommandSkipsWhenBaseURLSet(t *testing.T) {
	cmd := buildE2EServerCommand(
		e2eRunOptions{BaseURL: "http://localhost:4242"},
		appconfig.Config{Server: appconfig.ServerConfig{Command: "just build-e2e-server", SkipWhenBaseURLSet: true}},
	)
	if len(cmd) != 0 {
		t.Fatalf("server command = %#v, want nil", cmd)
	}
}

func TestBuildE2EServerCommandRunsSetupForManagedDefault(t *testing.T) {
	t.Setenv("E2E_BASE_URL", "")
	cmd := buildE2EServerCommand(
		e2eRunOptions{},
		appconfig.Config{Server: appconfig.ServerConfig{Command: "just build-e2e-server", SkipWhenBaseURLSet: true}},
	)
	assertStringSliceEqual(t, cmd, []string{"just", "build-e2e-server"})
}

func TestE2ERunCommandIncludesRunnerFlags(t *testing.T) {
	cmd := newE2ERunCommand(context.Background())
	for _, name := range []string{"all", "base-ref", "jobs", "readiness-path", "readiness-timeout"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Fatalf("missing --%s flag", name)
		}
	}
	baseRef, err := cmd.Flags().GetString("base-ref")
	if err != nil {
		t.Fatal(err)
	}
	if baseRef != "main" {
		t.Fatalf("base-ref default = %q, want main", baseRef)
	}
}

func TestE2ERunCommandRejectsUnexpectedArgsWithoutStory(t *testing.T) {
	cmd := newE2ERunCommand(context.Background())
	cmd.SetArgs([]string{"unexpected"})
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "unexpected e2e run arguments") {
		t.Fatalf("err = %v, want unexpected arguments error", err)
	}
}
