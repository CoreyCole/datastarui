package runtime

import (
	"testing"

	"github.com/playwright-community/playwright-go"
)

type Context struct {
	Config     Config
	Playwright *playwright.Playwright
	Browser    playwright.Browser
	Page       playwright.Page
	Console    *ConsoleMonitor
	Artifacts  ArtifactSink
	Memory     map[string]string
}

type ScenarioFunc func(t testing.TB, ctx *Context)

type ArtifactSink interface {
	Capture(label string, page playwright.Page) error
}
