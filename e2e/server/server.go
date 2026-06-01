package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/coreycole/datastarui/e2e/appconfig"
)

const defaultPortEnv = "PORT"

type Options struct {
	BaseURL   string
	RunDir    string
	NoRestart bool
}

type ManagedServer struct {
	cfg  appconfig.Config
	opts Options

	cmd     *exec.Cmd
	baseURL string
	logPath string
	done    chan error

	cleanupOnce sync.Once
}

func AllocatePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func New(cfg appconfig.Config, opts Options) *ManagedServer {
	return &ManagedServer{
		cfg:     cfg,
		opts:    opts,
		baseURL: strings.TrimRight(opts.BaseURL, "/"),
		done:    make(chan error, 1),
	}
}

func (s *ManagedServer) BaseURL() string { return strings.TrimRight(s.baseURL, "/") }
func (s *ManagedServer) LogPath() string { return s.logPath }

func (s *ManagedServer) Start(ctx context.Context) error {
	if s.opts.NoRestart || strings.TrimSpace(s.cfg.Server.ManagedCommand) == "" {
		if s.baseURL == "" {
			s.baseURL = strings.TrimRight(s.cfg.BaseURL, "/")
		}
		return nil
	}

	port, err := AllocatePort()
	if err != nil {
		return fmt.Errorf("allocate port: %w", err)
	}
	s.baseURL = fmt.Sprintf("http://127.0.0.1:%d", port)
	s.logPath = filepath.Join(s.opts.RunDir, "server.log")
	if err := os.MkdirAll(filepath.Dir(s.logPath), 0o755); err != nil {
		return err
	}
	logFile, err := os.Create(s.logPath)
	if err != nil {
		return err
	}

	command := strings.Fields(s.cfg.Server.ManagedCommand)
	if len(command) == 0 {
		logFile.Close()
		return fmt.Errorf("managed_command is empty")
	}
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Dir = s.cfg.RootDir
	cmd.Env = managedEnv(os.Environ(), envName(s.cfg.Server.PortEnv), port)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	setProcessGroup(cmd)

	if err := cmd.Start(); err != nil {
		logFile.Close()
		return fmt.Errorf("start managed server %q: %w", s.cfg.Server.ManagedCommand, err)
	}
	s.cmd = cmd
	go func() {
		err := cmd.Wait()
		logFile.Close()
		s.done <- err
	}()
	return nil
}

func (s *ManagedServer) WaitReady(ctx context.Context) error {
	readyCtx, cancel := context.WithTimeout(ctx, s.cfg.ReadinessTimeout())
	defer cancel()
	return ProbeReady(readyCtx, s.cfg.ReadinessURL(s.BaseURL()), 250*time.Millisecond)
}

func (s *ManagedServer) Cleanup(ctx context.Context) error {
	var err error
	s.cleanupOnce.Do(func() {
		if s.cmd == nil || s.cmd.Process == nil {
			return
		}
		err = KillProcessGroup(ctx, s.cmd)
		select {
		case <-s.done:
		case <-time.After(5 * time.Second):
		}
	})
	return err
}

func ProbeReady(ctx context.Context, readinessURL string, interval time.Duration) error {
	if interval <= 0 {
		interval = 250 * time.Millisecond
	}
	client := &http.Client{Timeout: 2 * time.Second}
	var lastErr error
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, readinessURL, nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err == nil && resp.Body != nil {
			resp.Body.Close()
		}
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 500 {
			return nil
		}
		if err != nil {
			lastErr = err
		} else {
			lastErr = fmt.Errorf("status %s", resp.Status)
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("readiness probe %s timed out: %w; last error: %v", readinessURL, ctx.Err(), lastErr)
		case <-time.After(interval):
		}
	}
}

func KillProcessGroup(ctx context.Context, cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	if runtime.GOOS == "windows" {
		return cmd.Process.Kill()
	}
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil && !errors.Is(err, syscall.ESRCH) {
		return err
	}
	select {
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
		return nil
	}
}

func setProcessGroup(cmd *exec.Cmd) {
	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}
}

func envName(name string) string {
	if strings.TrimSpace(name) == "" {
		return defaultPortEnv
	}
	return strings.TrimSpace(name)
}

func managedEnv(base []string, name string, port int) []string {
	prefix := name + "="
	env := make([]string, 0, len(base)+1)
	for _, entry := range base {
		if !strings.HasPrefix(entry, prefix) {
			env = append(env, entry)
		}
	}
	return append(env, fmt.Sprintf("%s=%d", name, port))
}
