package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type FileArtifactSink struct {
	Dir string
}

func NewFileArtifactSink(dir string) (*FileArtifactSink, error) {
	if dir == "" {
		return &FileArtifactSink{}, nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &FileArtifactSink{Dir: dir}, nil
}

func (s *FileArtifactSink) Capture(label string, page playwright.Page) error {
	if s == nil || s.Dir == "" {
		return nil
	}
	safeLabel := safeArtifactName(label)
	if safeLabel == "" {
		safeLabel = "page"
	}
	pngPath := filepath.Join(s.Dir, safeLabel+".png")
	png, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(pngPath),
		FullPage: playwright.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("screenshot %s: %w", label, err)
	}
	if len(png) == 0 {
		if _, err := os.Stat(pngPath); err != nil {
			return fmt.Errorf("screenshot %s produced no bytes and no file: %w", label, err)
		}
	}
	html, err := page.Content()
	if err != nil {
		return fmt.Errorf("html snapshot %s: %w", label, err)
	}
	if err := os.WriteFile(filepath.Join(s.Dir, safeLabel+".html"), []byte(html), 0o644); err != nil {
		return fmt.Errorf("write html snapshot %s: %w", label, err)
	}
	return nil
}

func safeArtifactName(label string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(strings.TrimSpace(label)) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
