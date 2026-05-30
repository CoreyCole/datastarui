package button

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestLinkButtonRendersAnchor(t *testing.T) {
	var body bytes.Buffer
	child := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, "Docs")
		return err
	})
	err := LinkButton(LinkButtonArgs{Href: "/docs", Target: "_blank", Rel: "noreferrer", Variant: "outline"}).Render(templ.WithChildren(t.Context(), child), &body)
	if err != nil {
		t.Fatal(err)
	}
	html := body.String()
	for _, want := range []string{`<a`, `href="/docs"`, `target="_blank"`, `rel="noreferrer"`, `Docs`} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %q: %s", want, html)
		}
	}
	if strings.Contains(html, "<button") {
		t.Fatalf("LinkButton rendered button: %s", html)
	}
}
