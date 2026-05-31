package runtime

import (
	"fmt"
	"strings"
	"sync"

	"github.com/playwright-community/playwright-go"
)

type ConsoleEntry struct {
	Type      string
	Text      string
	URL       string
	Line      int
	Column    int
	PageError bool
}

type ConsoleMonitor struct {
	mu      sync.Mutex
	entries []ConsoleEntry
}

func NewConsoleMonitor(page playwright.Page) *ConsoleMonitor {
	m := &ConsoleMonitor{}
	page.OnConsole(func(message playwright.ConsoleMessage) {
		entry := ConsoleEntry{Type: message.Type(), Text: message.Text()}
		if loc := message.Location(); loc != nil {
			entry.URL = loc.URL
			entry.Line = loc.LineNumber + 1
			entry.Column = loc.ColumnNumber + 1
		}
		m.add(entry)
	})
	page.OnPageError(func(err error) {
		m.add(ConsoleEntry{Type: "pageerror", Text: err.Error(), PageError: true})
	})
	return m
}

func (m *ConsoleMonitor) Problems() []ConsoleEntry {
	if m == nil {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	problems := make([]ConsoleEntry, 0)
	for _, entry := range m.entries {
		switch entry.Type {
		case "error", "warning", "pageerror":
			problems = append(problems, entry)
		}
	}
	return problems
}

func (m *ConsoleMonitor) add(entry ConsoleEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, entry)
}

func FormatConsoleProblems(entries []ConsoleEntry) string {
	if len(entries) == 0 {
		return ""
	}
	lines := make([]string, 0, len(entries))
	for _, entry := range entries {
		location := ""
		if entry.URL != "" {
			location = fmt.Sprintf(" (%s:%d:%d)", entry.URL, entry.Line, entry.Column)
		}
		lines = append(lines, fmt.Sprintf("[%s] %s%s", entry.Type, entry.Text, location))
	}
	return strings.Join(lines, "\n")
}
