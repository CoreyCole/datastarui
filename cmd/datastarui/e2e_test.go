package main

import (
	"context"
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
		appconfig.Config{Server: appconfig.ServerConfig{Command: "just build-local", SkipWhenBaseURLSet: true}},
	)
	if len(cmd) != 0 {
		t.Fatalf("server command = %#v, want nil", cmd)
	}
}
