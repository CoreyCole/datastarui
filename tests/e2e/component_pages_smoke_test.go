package e2e_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/spec"
)

func TestComponentPagesRender(t *testing.T) {
	pages := []struct {
		name    string
		path    string
		heading string
	}{
		{name: "breadcrumb", path: "/components/breadcrumb", heading: "Breadcrumb"},
		{name: "button", path: "/components/button", heading: "Button"},
		{name: "calendar", path: "/components/calendar", heading: "Calendar"},
		{name: "card", path: "/components/card", heading: "Card"},
		{name: "checkbox", path: "/components/checkbox", heading: "Checkbox"},
		{name: "datepicker", path: "/components/datepicker", heading: "Datepicker"},
		{name: "dialog", path: "/components/dialog", heading: "Dialog"},
		{name: "dropdown", path: "/components/dropdown", heading: "Dropdown"},
		{name: "form", path: "/components/form", heading: "Form"},
		{name: "popover", path: "/components/popover", heading: "Popover"},
		{name: "select", path: "/components/select", heading: "Select"},
		{name: "sheet", path: "/components/sheet", heading: "Sheet"},
		{name: "sidebar", path: "/components/sidebar", heading: "Sidebar"},
		{name: "tabs", path: "/components/tabs", heading: "Tabs"},
		{name: "tooltip", path: "/components/tooltip", heading: "Tooltip"},
	}

	for _, page := range pages {
		page := page
		t.Run(page.name, func(t *testing.T) {
			spec.Story(t, page.name+" component page renders").
				Visit(spec.Path(page.path)).
				Expect(spec.ExpectStep(spec.Visible(spec.Role("heading", page.heading)))).
				Expect(spec.ExpectStep(spec.ConsoleClean())).
				Run()
		})
	}
}
