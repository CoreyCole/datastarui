package registry

import (
	"fmt"
	"slices"
)

// ComponentManifest describes the DatastarUI source files copied for a component.
type ComponentManifest struct {
	Name         string
	Files        []string
	Dependencies []string
	CSS          []string
}

// Registry contains the source-copy manifests for DatastarUI components.
type Registry struct {
	Components map[string]ComponentManifest
}

// Default returns the built-in registry used by the copy CLI.
func Default() Registry {
	components := map[string]ComponentManifest{
		"avatar":     component("avatar", []string{"args.go", "avatar.templ"}),
		"breadcrumb": component("breadcrumb", []string{"args.go", "breadcrumb.templ", "fromitems.templ", "variants.go"}),
		"button":     component("button", []string{"args.go", "button.templ", "variants.go"}),
		"card":       component("card", []string{"args.go", "card.templ", "variants.go"}),
		"checkbox":   component("checkbox", []string{"args.go", "checkbox.templ", "variants.go"}),
		"dialog":     component("dialog", []string{"args.go", "dialog.templ", "expressions.go", "variants.go"}, "utils"),
		"dropdown":   component("dropdown", []string{"args.go", "dropdown.templ", "expressions.go", "variants.go"}, "utils"),
		"form":       component("form", []string{"args.go", "form.templ", "variants.go"}),
		"input":      component("input", []string{"args.go", "input.templ", "variants.go"}),
		"label":      component("label", []string{"args.go", "label.templ", "variants.go"}),
		"select":     component("select", []string{"args.go", "expressions.go", "select.templ", "variants.go"}, "utils"),
		"sheet":      component("sheet", []string{"args.go", "expressions.go", "sheet.templ", "variants.go"}, "utils"),
		"tabs":       component("tabs", []string{"args.go", "expressions.go", "tabs.templ", "variants.go"}, "utils"),
		"textarea":   component("textarea", []string{"args.go", "textarea.templ", "variants.go"}),
		"tooltip":    component("tooltip", []string{"args.go", "expressions.go", "tooltip.templ", "variants.go"}, "utils"),
		"utils": {
			Name: "utils",
			Files: []string{
				"utils/anchor.go",
				"utils/connect_errors.go",
				"utils/context.go",
				"utils/data_class.go",
				"utils/device.go",
				"utils/expressions.go",
				"utils/signals.go",
				"utils/tailwind_merge.go",
			},
		},
		"tailwind": {
			Name:  "tailwind",
			Files: []string{"tailwind/utilities.css", "tailwind/README.md"},
		},
	}
	return Registry{Components: components}
}

func component(name string, files []string, deps ...string) ComponentManifest {
	paths := make([]string, 0, len(files))
	for _, file := range files {
		paths = append(paths, fmt.Sprintf("components/%s/%s", name, file))
	}
	return ComponentManifest{Name: name, Files: paths, Dependencies: deps}
}

// Resolve returns manifests for names and transitive dependencies.
func (r Registry) Resolve(names []string) ([]ComponentManifest, error) {
	seen := map[string]bool{}
	var out []ComponentManifest
	var visit func(string) error
	visit = func(name string) error {
		if seen[name] {
			return nil
		}
		m, ok := r.Components[name]
		if !ok {
			return UnknownComponentError{Name: name}
		}
		seen[name] = true
		for _, dep := range m.Dependencies {
			if err := visit(dep); err != nil {
				return err
			}
		}
		out = append(out, m)
		return nil
	}
	for _, name := range names {
		if err := visit(name); err != nil {
			return nil, err
		}
	}
	slices.SortFunc(out, func(a, b ComponentManifest) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})
	return out, nil
}

// UnknownComponentError reports an unknown component name.
type UnknownComponentError struct{ Name string }

func (e UnknownComponentError) Error() string { return "unknown datastarui component: " + e.Name }
