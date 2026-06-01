package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/coreycole/datastarui/e2e/appconfig"
)

func TestAllocatePortReturnsConnectablePort(t *testing.T) {
	port, err := AllocatePort()
	if err != nil {
		t.Fatal(err)
	}
	if port == 0 {
		t.Fatal("expected non-zero port")
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("allocated port should be reusable after release: %v", err)
	}
	listener.Close()
}

func TestProbeReadyWaitsForHTTPServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := ProbeReady(ctx, server.URL, 10*time.Millisecond); err != nil {
		t.Fatal(err)
	}
}

func TestProbeReadyTimesOut(t *testing.T) {
	port, err := AllocatePort()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err = ProbeReady(ctx, fmt.Sprintf("http://127.0.0.1:%d", port), 5*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout")
	}
	if !strings.Contains(err.Error(), "readiness probe") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestManagedServerUsesConfiguredBaseURLWhenRestartSkipped(t *testing.T) {
	cfg := appconfig.Config{BaseURL: "http://localhost:4242/"}
	srv := New(cfg, Options{NoRestart: true})
	if err := srv.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := srv.BaseURL(); got != "http://localhost:4242" {
		t.Fatalf("BaseURL = %q", got)
	}
	if srv.LogPath() != "" {
		t.Fatalf("LogPath = %q", srv.LogPath())
	}
}

func TestManagedEnvInjectsConfiguredPort(t *testing.T) {
	env := managedEnv([]string{"PORT=old", "KEEP=yes"}, "PORT", 12345)
	joined := strings.Join(env, "\n")
	if strings.Contains(joined, "PORT=old") {
		t.Fatalf("old PORT retained in %v", env)
	}
	if !strings.Contains(joined, "PORT=12345") || !strings.Contains(joined, "KEEP=yes") {
		t.Fatalf("env = %v", env)
	}
}

func TestStartCreatesServerLogAndInjectsPortEnv(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell command test is unix-only")
	}
	root := t.TempDir()
	portFile := filepath.Join(root, "port.txt")
	script := filepath.Join(root, "server.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\necho $PORT > port.txt\nsleep 30\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := appconfig.Config{
		RootDir: root,
		Server: appconfig.ServerConfig{
			ManagedCommand: script,
			PortEnv:        "PORT",
		},
	}
	srv := New(cfg, Options{RunDir: filepath.Join(root, "run")})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := srv.Start(ctx); err != nil {
		t.Fatal(err)
	}
	defer srv.Cleanup(context.Background())
	if srv.BaseURL() == "" || srv.LogPath() == "" {
		t.Fatalf("baseURL/logPath = %q/%q", srv.BaseURL(), srv.LogPath())
	}
	if _, err := os.Stat(srv.LogPath()); err != nil {
		t.Fatalf("server log missing: %v", err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(portFile)
		if err == nil && strings.TrimSpace(string(data)) != "" {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("managed command did not receive PORT env")
}

func TestKillProcessGroupHandlesNil(t *testing.T) {
	if err := KillProcessGroup(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if err := KillProcessGroup(context.Background(), &exec.Cmd{}); err != nil {
		t.Fatal(err)
	}
}
