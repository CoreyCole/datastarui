package dropdown

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func renderWithText(t *testing.T, component templ.Component, text string) string {
	t.Helper()
	var body bytes.Buffer
	child := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, text)
		return err
	})
	if err := component.Render(templ.WithChildren(t.Context(), child), &body); err != nil {
		t.Fatal(err)
	}
	return body.String()
}

func TestDropdownMenuLinkItemRendersAnchor(t *testing.T) {
	html := renderWithText(t, DropdownMenuLinkItem(DropdownMenuLinkItemArgs{ID: "workspace_actions", Href: "/docs", Target: "_blank", Rel: "noreferrer"}), "Docs")
	for _, want := range []string{`<a`, `role="menuitem"`, `href="/docs"`, `target="_blank"`, `rel="noreferrer"`, `Docs`} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %q: %s", want, html)
		}
	}
	if strings.Contains(html, "<span") {
		t.Fatalf("link item rendered span wrapper: %s", html)
	}
}

func TestDropdownMenuFormItemRendersFormWithSubmitButton(t *testing.T) {
	html := renderWithText(t, DropdownMenuFormItem(DropdownMenuFormItemArgs{ID: "workspace_actions", Action: "/workspaces/example/start"}), "Start")
	for _, want := range []string{`<form`, `action="/workspaces/example/start"`, `method="post"`, `<button`, `type="submit"`, `role="menuitem"`, `Start`} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %q: %s", want, html)
		}
	}
	if strings.Contains(html, `<span`) {
		t.Fatalf("form item rendered span wrapper: %s", html)
	}
}

func TestDropdownMenuCustomItemUsesFlowContainer(t *testing.T) {
	html := renderWithText(t, DropdownMenuCustomItem(DropdownMenuCustomItemArgs{ID: "workspace_actions", CloseOnClick: true}), "Custom")
	for _, want := range []string{`<div`, `data-slot="dropdown-menu-item"`, `role="menuitem"`, `data-on-click`, `Custom`} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %q: %s", want, html)
		}
	}
	if strings.Contains(html, `<span`) {
		t.Fatalf("custom item rendered span wrapper: %s", html)
	}
}

func TestDropdownMenuItemAsChildDoesNotWrapFlowContentInSpan(t *testing.T) {
	html := renderWithText(t, DropdownMenuItem(DropdownMenuItemArgs{ID: "workspace_actions", AsChild: true}), "Child")
	if strings.Contains(html, `<span`) {
		t.Fatalf("AsChild item rendered span wrapper: %s", html)
	}
	if !strings.Contains(html, `<div`) || !strings.Contains(html, `role="menuitem"`) {
		t.Fatalf("AsChild item did not render flow container menuitem: %s", html)
	}
}
