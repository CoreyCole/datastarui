package review

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func WritePlanArtifacts(planDir string, result ReviewResult) (ReviewResult, error) {
	return WritePlanArtifactsForRun(planDir, "", result)
}

func WritePlanArtifactsForRun(planDir, runPath string, result ReviewResult) (ReviewResult, error) {
	dir := filepath.Join(planDir, "context", "implement", "e2e-runs", time.Now().UTC().Format("20060102T150405Z"))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return result, err
	}
	if strings.TrimSpace(runPath) != "" {
		indexPath, err := WriteReviewIndex(runPath, dir)
		if err != nil {
			return result, err
		}
		if indexPath != "" {
			result.IndexPath = indexPath
		}
	}
	markdownPath := filepath.Join(dir, "e2e-visual.md")
	jsonPath := filepath.Join(dir, "e2e-visual.json")
	result.MarkdownPath = markdownPath
	result.JSONPath = jsonPath
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return result, err
	}
	if err := os.WriteFile(jsonPath, data, 0o644); err != nil {
		return result, err
	}
	if err := WriteMarkdown(markdownPath, result); err != nil {
		return result, err
	}
	return result, nil
}

func WriteReviewIndex(runPath, outDir string) (string, error) {
	shots, err := collectScreenshots(runPath)
	if err != nil {
		return "", err
	}
	if len(shots) == 0 {
		return "", nil
	}
	var body strings.Builder
	body.WriteString(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>DatastarUI E2E visual review</title>
  <style>
    body{font-family:system-ui,sans-serif;margin:0;background:#0f172a;color:#e2e8f0}main{max-width:1200px;margin:0 auto;padding:24px}article{margin:24px 0;padding:16px;border:1px solid #334155;border-radius:12px;background:#111827}img{max-width:100%;border:1px solid #334155;border-radius:8px;background:white}code{color:#93c5fd;word-break:break-all}summary{cursor:pointer;margin-top:12px}pre{max-height:420px;overflow:auto;white-space:pre-wrap;background:#020617;padding:12px;border-radius:8px}</style>
</head>
<body><main>
  <h1>DatastarUI E2E visual review</h1>
`)
	body.WriteString("  <p>Run: <code>" + html.EscapeString(runPath) + "</code></p>\n")
	for _, shot := range shots {
		body.WriteString("  <article>\n")
		body.WriteString("    <h2>" + html.EscapeString(shot.Label) + "</h2>\n")
		body.WriteString("    <p><code>" + html.EscapeString(shot.Path) + "</code></p>\n")
		body.WriteString("    <img alt=\"" + html.EscapeString(shot.Label) + "\" src=\"data:image/png;base64," + shot.PNGBase64 + "\">\n")
		if shot.HTML != "" {
			body.WriteString("    <details><summary>HTML snapshot</summary><pre>" + html.EscapeString(shot.HTML) + "</pre></details>\n")
		}
		body.WriteString("  </article>\n")
	}
	body.WriteString("</main></body></html>\n")
	path := filepath.Join(outDir, "index.html")
	if err := os.WriteFile(path, []byte(body.String()), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

type screenshotArtifact struct {
	Label     string
	Path      string
	PNGBase64 string
	HTML      string
}

func collectScreenshots(runPath string) ([]screenshotArtifact, error) {
	var shots []screenshotArtifact
	if err := filepath.WalkDir(runPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || strings.ToLower(filepath.Ext(path)) != ".png" {
			return nil
		}
		png, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(runPath, path)
		if err != nil {
			rel = path
		}
		shot := screenshotArtifact{
			Label:     strings.TrimSuffix(filepath.ToSlash(rel), filepath.Ext(rel)),
			Path:      filepath.ToSlash(rel),
			PNGBase64: base64.StdEncoding.EncodeToString(png),
		}
		htmlPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".html"
		if data, err := os.ReadFile(htmlPath); err == nil {
			shot.HTML = string(data)
		}
		shots = append(shots, shot)
		return nil
	}); err != nil {
		return nil, err
	}
	sort.Slice(shots, func(i, j int) bool { return shots[i].Path < shots[j].Path })
	return shots, nil
}

func WriteMarkdown(path string, result ReviewResult) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	body := fmt.Sprintf(`---
date: %s
tool: datastarui e2e review
review_type: e2e_visual
status: complete
outcome: %s
---

# E2E Visual Review

## Summary

%s

## Comparisons
`, time.Now().Format(time.RFC3339), result.Outcome, result.Summary)
	for _, cmp := range result.Comparisons {
		body += fmt.Sprintf("- %s/%s/%s: %s — %s\n", cmp.Story, cmp.Scenario, cmp.Viewport, cmp.Outcome, cmp.Rationale)
	}
	return os.WriteFile(path, []byte(body), 0o644)
}
