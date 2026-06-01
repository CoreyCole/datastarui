package runner

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

func ChangedFiles(ctx context.Context, repoRoot, baseRef string) ([]ChangedFile, error) {
	if strings.TrimSpace(baseRef) == "" {
		baseRef = "main"
	}
	mergeBase, err := gitOutput(ctx, repoRoot, "merge-base", "HEAD", baseRef)
	if err != nil {
		return nil, fmt.Errorf("find merge-base with %s: %w", baseRef, err)
	}
	committed, err := gitNameStatus(ctx, repoRoot, strings.TrimSpace(mergeBase)+"...HEAD")
	if err != nil {
		return nil, err
	}
	unstaged, err := gitNameStatus(ctx, repoRoot)
	if err != nil {
		return nil, err
	}
	staged, err := gitNameStatus(ctx, repoRoot, "--cached")
	if err != nil {
		return nil, err
	}
	untracked, err := gitUntracked(ctx, repoRoot)
	if err != nil {
		return nil, err
	}
	return dedupeChangedFiles(append(append(append(committed, staged...), unstaged...), untracked...)), nil
}

func gitNameStatus(ctx context.Context, repoRoot string, args ...string) ([]ChangedFile, error) {
	gitArgs := append([]string{"diff", "--name-status"}, args...)
	out, err := gitOutput(ctx, repoRoot, gitArgs...)
	if err != nil {
		return nil, fmt.Errorf("git %s: %w", strings.Join(gitArgs, " "), err)
	}
	return parseNameStatus(out), nil
}

func gitUntracked(ctx context.Context, repoRoot string) ([]ChangedFile, error) {
	out, err := gitOutput(ctx, repoRoot, "ls-files", "--others", "--exclude-standard")
	if err != nil {
		return nil, fmt.Errorf("git ls-files --others --exclude-standard: %w", err)
	}
	return parseUntracked(out), nil
}

func gitOutput(ctx context.Context, repoRoot string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func parseNameStatus(out string) []ChangedFile {
	var files []ChangedFile
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 2 {
			continue
		}
		status := fields[0]
		if strings.HasPrefix(status, "R") || strings.HasPrefix(status, "C") {
			if len(fields) >= 3 {
				files = append(files, ChangedFile{Status: status[:1], OldPath: fields[1], Path: fields[2]})
			}
			continue
		}
		files = append(files, ChangedFile{Status: status[:1], Path: fields[1]})
	}
	return files
}

func parseUntracked(out string) []ChangedFile {
	var files []ChangedFile
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		path := strings.TrimSpace(scanner.Text())
		if path != "" {
			files = append(files, ChangedFile{Status: "A", Path: path})
		}
	}
	return files
}

func dedupeChangedFiles(files []ChangedFile) []ChangedFile {
	seen := map[string]bool{}
	out := make([]ChangedFile, 0, len(files))
	for _, file := range files {
		key := file.Status + "\x00" + file.OldPath + "\x00" + file.Path
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, file)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path == out[j].Path {
			return out[i].OldPath < out[j].OldPath
		}
		return out[i].Path < out[j].Path
	})
	return out
}
