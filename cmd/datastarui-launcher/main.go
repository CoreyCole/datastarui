package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"syscall"
)

func main() {
	if err := maybeReexecManaged(); err != nil {
		fmt.Fprintf(os.Stderr, "datastarui launcher error: %v\n", err)
		os.Exit(1)
	}
}

func maybeReexecManaged() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve launcher executable: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(exePath); err == nil {
		exePath = resolved
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(exePath), ".."))
	fingerprint, err := runtimeFingerprint(repoRoot)
	if err != nil {
		return err
	}
	binaryPath := filepath.Join(repoRoot, "bin", "datastarui-runtime-"+fingerprint[:12])
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		if err := buildManagedRuntime(repoRoot, binaryPath); err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("stat managed runtime %q: %w", binaryPath, err)
	}
	return syscall.Exec(binaryPath, append([]string{binaryPath}, os.Args[1:]...), os.Environ())
}

func runtimeFingerprint(repoRoot string) (string, error) {
	files := []string{"go.mod", "go.sum"}
	for _, root := range []string{"cmd/datastarui", "e2e"} {
		if err := filepath.WalkDir(filepath.Join(repoRoot, root), func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || filepath.Ext(path) != ".go" {
				return nil
			}
			rel, err := filepath.Rel(repoRoot, path)
			if err != nil {
				return err
			}
			files = append(files, filepath.ToSlash(rel))
			return nil
		}); err != nil {
			return "", fmt.Errorf("walk %s: %w", root, err)
		}
	}
	sort.Strings(files)
	h := sha256.New()
	fmt.Fprintf(h, "goos=%s\ngoarch=%s\n", runtime.GOOS, runtime.GOARCH)
	for _, rel := range files {
		path := filepath.Join(repoRoot, filepath.FromSlash(rel))
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return "", fmt.Errorf("stat %s: %w", rel, err)
		}
		fmt.Fprintf(h, "file:%s\n", rel)
		f, err := os.Open(path)
		if err != nil {
			return "", fmt.Errorf("open %s: %w", rel, err)
		}
		if _, err := io.Copy(h, f); err != nil {
			_ = f.Close()
			return "", fmt.Errorf("hash %s: %w", rel, err)
		}
		if err := f.Close(); err != nil {
			return "", fmt.Errorf("close %s: %w", rel, err)
		}
		fmt.Fprintln(h)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func buildManagedRuntime(repoRoot, binaryPath string) error {
	if err := os.MkdirAll(filepath.Dir(binaryPath), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(binaryPath), filepath.Base(binaryPath)+".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	if err := tmp.Close(); err != nil {
		return err
	}
	defer os.Remove(tmpPath)
	cmd := exec.Command("go", "build", "-o", tmpPath, "./cmd/datastarui")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build managed runtime: %w", err)
	}
	if err := os.Chmod(tmpPath, 0o755); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, binaryPath); err != nil {
		return err
	}
	return nil
}
