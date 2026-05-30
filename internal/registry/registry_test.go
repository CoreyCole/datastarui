package registry

import (
	"errors"
	"slices"
	"testing"
)

func TestResolveIncludesDependencies(t *testing.T) {
	t.Parallel()

	manifests, err := Default().Resolve([]string{"select"})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	var names []string
	for _, manifest := range manifests {
		names = append(names, manifest.Name)
	}
	if !slices.Contains(names, "utils") || !slices.Contains(names, "select") {
		t.Fatalf("Resolve(select) names = %v, want utils and select", names)
	}
}

func TestResolveUnknownComponent(t *testing.T) {
	t.Parallel()

	_, err := Default().Resolve([]string{"nope"})
	var unknown UnknownComponentError
	if !errors.As(err, &unknown) {
		t.Fatalf("Resolve(unknown) error = %T %v, want UnknownComponentError", err, err)
	}
}

func TestDefaultRegistryContainsOutlineComponents(t *testing.T) {
	t.Parallel()

	registry := Default()
	for _, name := range []string{"avatar", "breadcrumb", "button", "card", "checkbox", "dialog", "dropdown", "form", "input", "label", "select", "sheet", "tabs", "textarea", "tooltip", "utils", "tailwind"} {
		if _, ok := registry.Components[name]; !ok {
			t.Fatalf("Default registry missing %q", name)
		}
	}
}

func TestComponentManifestIncludesCompletePublicSource(t *testing.T) {
	t.Parallel()

	breadcrumb := Default().Components["breadcrumb"]
	if !slices.Contains(breadcrumb.Files, "components/breadcrumb/fromitems.templ") {
		t.Fatalf("breadcrumb files = %v, want fromitems.templ", breadcrumb.Files)
	}
}
