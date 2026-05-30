package copy

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/coreycole/datastarui/internal/registry"
)

// Options configures a copied-source operation.
type Options struct {
	SourceRoot   string
	TargetRoot   string
	TargetModule string
	Components   []string
	Update       bool
	Force        bool // reserved for explicit future overwrite/merge policy; default protects drift
}

// Result describes a copied-source operation result.
type Result struct {
	CopiedFiles []string
	LockPath    string
	Drift       []DriftEntry
}

// DriftEntry describes target drift from datastarui.lock.json.
type DriftEntry struct {
	Path   string
	Status string // modified | missing | unmanaged
}

// LockFile records copied-source provenance.
type LockFile struct {
	Source     string                     `json:"source"`
	Commit     string                     `json:"commit"`
	ModulePath string                     `json:"modulePath"`
	Components map[string]LockedComponent `json:"components"`
	Files      map[string]LockedFile      `json:"files"`
}

// LockedComponent records copied files for one component.
type LockedComponent struct {
	Version string   `json:"version"`
	Files   []string `json:"files"`
}

// LockedFile records source path and copied hash for one target file.
type LockedFile struct {
	SourcePath string `json:"sourcePath"`
	Hash       string `json:"hash"`
}

// Add copies selected DatastarUI source into a consumer target.
func Add(ctx context.Context, opts Options) (Result, error) {
	opts.Update = false
	return copyComponents(ctx, opts)
}

// Update refreshes selected DatastarUI source, refusing drift by default.
func Update(ctx context.Context, opts Options) (Result, error) {
	opts.Update = true
	return copyComponents(ctx, opts)
}

// Diff reports drift from datastarui.lock.json.
func Diff(_ context.Context, opts Options) (Result, error) {
	lock, lockPath, err := readLock(opts.TargetRoot)
	if err != nil {
		return Result{}, err
	}
	var drift []DriftEntry
	for rel, locked := range lock.Files {
		gotHash, err := hashFile(filepath.Join(opts.TargetRoot, rel))
		if errors.Is(err, os.ErrNotExist) {
			drift = append(drift, DriftEntry{Path: rel, Status: "missing"})
			continue
		}
		if err != nil {
			return Result{}, err
		}
		if gotHash != locked.Hash {
			drift = append(drift, DriftEntry{Path: rel, Status: "modified"})
		}
	}
	unmanaged, err := unmanagedFiles(opts.TargetRoot, lock)
	if err != nil {
		return Result{}, err
	}
	drift = append(drift, unmanaged...)
	sortDrift(drift)
	return Result{LockPath: lockPath, Drift: drift}, nil
}

// Doctor fails when the target has drift.
func Doctor(ctx context.Context, opts Options) error {
	result, err := Diff(ctx, opts)
	if err != nil {
		return err
	}
	if len(result.Drift) > 0 {
		return fmt.Errorf("datastarui copied source drift: %d files", len(result.Drift))
	}
	return nil
}

func copyComponents(ctx context.Context, opts Options) (Result, error) {
	if strings.TrimSpace(opts.SourceRoot) == "" || strings.TrimSpace(opts.TargetRoot) == "" {
		return Result{}, fmt.Errorf("source and target are required")
	}
	if strings.TrimSpace(opts.TargetModule) == "" {
		return Result{}, fmt.Errorf("target module is required")
	}
	if opts.Update && !opts.Force {
		current, err := Diff(ctx, opts)
		if err == nil && len(current.Drift) > 0 {
			return current, fmt.Errorf("refusing to overwrite modified datastarui copied source: %d drift entries", len(current.Drift))
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return Result{}, err
		}
	}
	components := opts.Components
	if len(components) == 0 {
		components = []string{"avatar", "breadcrumb", "button", "card", "checkbox", "dialog", "dropdown", "form", "input", "label", "select", "sheet", "tabs", "textarea", "tooltip", "utils", "tailwind"}
	}
	resolved, err := registry.Default().Resolve(components)
	if err != nil {
		return Result{}, err
	}
	commit := sourceCommit(ctx, opts.SourceRoot)
	lock := LockFile{
		Source:     filepath.Clean(opts.SourceRoot),
		Commit:     commit,
		ModulePath: opts.TargetModule,
		Components: map[string]LockedComponent{},
		Files:      map[string]LockedFile{},
	}
	var copied []string
	for _, component := range resolved {
		var componentFiles []string
		for _, srcRel := range component.Files {
			targetRel := rewriteTargetPath(srcRel)
			if err := copyOne(opts.SourceRoot, opts.TargetRoot, srcRel, targetRel, opts.TargetModule); err != nil {
				return Result{}, err
			}
			hash, err := hashFile(filepath.Join(opts.TargetRoot, targetRel))
			if err != nil {
				return Result{}, err
			}
			lock.Files[targetRel] = LockedFile{SourcePath: srcRel, Hash: hash}
			componentFiles = append(componentFiles, targetRel)
			copied = append(copied, targetRel)
		}
		sort.Strings(componentFiles)
		lock.Components[component.Name] = LockedComponent{Version: commit, Files: componentFiles}
	}
	if err := writeLock(opts.TargetRoot, lock); err != nil {
		return Result{}, err
	}
	sort.Strings(copied)
	return Result{CopiedFiles: copied, LockPath: filepath.Join(opts.TargetRoot, "datastarui.lock.json")}, nil
}

func rewriteTargetPath(srcRel string) string { return filepath.ToSlash(srcRel) }

func copyOne(sourceRoot, targetRoot, srcRel, targetRel, targetModule string) error {
	data, err := os.ReadFile(filepath.Join(sourceRoot, srcRel))
	if err != nil {
		return err
	}
	data = []byte(strings.ReplaceAll(string(data), "github.com/coreycole/datastarui", targetModule+"/pkg/datastarui"))
	outPath := filepath.Join(targetRoot, targetRel)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outPath, data, 0o644)
}

func writeLock(targetRoot string, lock LockFile) error {
	data, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(targetRoot, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(targetRoot, "datastarui.lock.json"), data, 0o644)
}

func readLock(targetRoot string) (LockFile, string, error) {
	path := filepath.Join(targetRoot, "datastarui.lock.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return LockFile{}, path, err
	}
	var lock LockFile
	if err := json.Unmarshal(data, &lock); err != nil {
		return LockFile{}, path, err
	}
	return lock, path, nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func sourceCommit(ctx context.Context, sourceRoot string) string {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
	cmd.Dir = sourceRoot
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func unmanagedFiles(targetRoot string, lock LockFile) ([]DriftEntry, error) {
	locked := map[string]bool{"datastarui.lock.json": true, "AGENTS.md": true}
	for rel := range lock.Files {
		locked[filepath.ToSlash(rel)] = true
	}
	var drift []DriftEntry
	err := filepath.WalkDir(targetRoot, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(targetRoot, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if locked[rel] || isGeneratedTempl(rel) {
			return nil
		}
		drift = append(drift, DriftEntry{Path: rel, Status: "unmanaged"})
		return nil
	})
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return drift, err
}

func isGeneratedTempl(path string) bool {
	return strings.HasSuffix(path, "_templ.go")
}

func sortDrift(drift []DriftEntry) {
	sort.Slice(drift, func(i, j int) bool {
		if drift[i].Path == drift[j].Path {
			return drift[i].Status < drift[j].Status
		}
		return drift[i].Path < drift[j].Path
	})
}
