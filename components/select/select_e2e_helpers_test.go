package selectcomponent_test

import (
	"fmt"
	"testing"

	"github.com/coreycole/datastarui/e2e/runtime"
	"github.com/coreycole/datastarui/e2e/spec"
)

type Theme string

const (
	ThemeLight  Theme = "light"
	ThemeDark   Theme = "dark"
	ThemeSystem Theme = "system"
)

type SelectControl struct{ id string }

func SliceThemeSelect() SelectControl  { return SelectControl{id: "slice_theme_select"} }
func FormThemeSelect() SelectControl   { return SelectControl{id: "form_theme_select"} }
func TabCategorySelect() SelectControl { return SelectControl{id: "tab_category_select"} }
func TabPrioritySelect() SelectControl { return SelectControl{id: "tab_priority_select"} }

func TabNameInput() spec.Locator  { return spec.CSS("#name_input") }
func TabEmailInput() spec.Locator { return spec.CSS("#email_input") }
func SubmitButton() spec.Locator  { return spec.CSS("#tab_test_form button[type='submit']") }

func (s SelectControl) Root() spec.Locator {
	return spec.CSS(fmt.Sprintf("[data-select-id='%s']", s.id))
}
func (s SelectControl) Trigger() spec.Locator {
	return spec.CSS(fmt.Sprintf("[data-select-id='%s'] [data-slot='select-trigger']", s.id))
}
func (s SelectControl) Content() spec.Locator {
	return spec.CSS(fmt.Sprintf("[data-select-id='%s'] [data-slot='select-content']", s.id))
}
func (s SelectControl) Option(value string) spec.Locator {
	return spec.CSS(fmt.Sprintf("[data-select-id='%s'] [data-select-item][data-value='%s']", s.id, value))
}
func (s SelectControl) Value() spec.Locator {
	return spec.CSS(fmt.Sprintf("[data-select-id='%s'] [data-slot='select-value']", s.id))
}
func (s SelectControl) HiddenInput() spec.Locator {
	return spec.CSS(fmt.Sprintf("[data-select-id='%s'] input[type='hidden']", s.id))
}
func (s SelectControl) Open() spec.Step { return spec.Click(s.Trigger()) }
func (s SelectControl) Choose(value string) spec.Step {
	return spec.Custom("choose "+value+" in "+s.id, func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		s.Open().Run(t, ctx)
		spec.Click(s.Option(value)).Run(t, ctx)
	})
}
func (s SelectControl) Opened() spec.Expectation { return spec.ExpectStep(spec.Visible(s.Content())) }
func (s SelectControl) Closed() spec.Expectation { return spec.ExpectStep(spec.Hidden(s.Content())) }
func (s SelectControl) Expanded(want bool) spec.Expectation {
	if want {
		return spec.AttributeEquals(s.Trigger(), spec.AriaExpanded, "true")
	}
	return spec.AttributeEquals(s.Trigger(), spec.AriaExpanded, "false")
}
func (s SelectControl) Focused() spec.Expectation { return spec.Focused(s.Trigger()) }
func (s SelectControl) ValueEquals(label string) spec.Expectation {
	return spec.TextContains(s.Value(), label)
}
