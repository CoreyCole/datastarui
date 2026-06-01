package dropdown_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/coreycole/datastarui/e2e/runtime"
	"github.com/coreycole/datastarui/e2e/spec"
)

type Theme string

const (
	ThemeMonokai Theme = "monokai"
	ThemeDracula Theme = "dracula"
	ThemeGitHub  Theme = "github"
	ThemeNord    Theme = "nord"
	ThemeNext    Theme = "next"
)

type DropdownDemo struct{ id string }
type SelectControl struct{ id string }

func DropdownPage() spec.Page              { return spec.Path("/components/dropdown") }
func SelectOverflowDropdown() DropdownDemo { return DropdownDemo{id: "select_dropdown"} }
func NumberedDropdown(index int) DropdownDemo {
	return DropdownDemo{id: fmt.Sprintf("test-dropdown-%d", index)}
}
func UserProfileDropdown() DropdownDemo       { return DropdownDemo{id: "user_profile_test"} }
func ManualCompositionDropdown() DropdownDemo { return DropdownDemo{id: "manual_composition_test"} }
func DropdownSelect() SelectControl           { return SelectControl{id: "dropdown_select"} }
func NestedThemeSelect() SelectControl        { return SelectControl{id: "nested_theme_select"} }
func ManualSyntaxSelect() SelectControl       { return SelectControl{id: "manual_syntax_select"} }

func (d DropdownDemo) Root() spec.Locator {
	return spec.CSS(fmt.Sprintf("[data-slot='dropdown-menu']:has([data-slot='dropdown-menu-trigger'][data-on-click*='%s'])", d.signalID()))
}
func (d DropdownDemo) Trigger() spec.Locator {
	return spec.CSS(fmt.Sprintf("[data-slot='dropdown-menu-trigger'][data-on-click*='%s']", d.signalID()))
}
func (d DropdownDemo) Content() spec.Locator {
	return spec.CSS(fmt.Sprintf("[data-slot='dropdown-menu-content'][data-show='$%s.open']", d.signalID()))
}
func (d DropdownDemo) signalID() string         { return strings.ReplaceAll(d.id, "-", "_") }
func (d DropdownDemo) Open() spec.Step          { return spec.Click(d.Trigger()) }
func (d DropdownDemo) Opened() spec.Expectation { return spec.ExpectStep(spec.Visible(d.Content())) }
func (d DropdownDemo) Closed() spec.Expectation { return spec.ExpectStep(spec.Hidden(d.Content())) }

func (s SelectControl) Root() spec.Locator {
	return spec.CSS(fmt.Sprintf("[data-select-id='%s']", s.id))
}
func (s SelectControl) Trigger() spec.Locator {
	return spec.CSS(fmt.Sprintf("[data-select-id='%s'] [data-slot='select-trigger']", s.id))
}
func (s SelectControl) Content() spec.Locator {
	return spec.CSS(fmt.Sprintf("[data-select-id='%s'] [data-slot='select-content']", s.id))
}
func (s SelectControl) Option(value Theme) spec.Locator {
	return spec.CSS(fmt.Sprintf("[data-select-id='%s'] [data-select-item][data-value='%s']", s.id, value))
}
func (s SelectControl) Value() spec.Locator {
	return spec.CSS(fmt.Sprintf("[data-select-id='%s'] [data-slot='select-value']", s.id))
}
func (s SelectControl) Open() spec.Step { return spec.Click(s.Trigger()) }
func (s SelectControl) Choose(value Theme) spec.Step {
	return spec.Custom("choose "+string(value)+" in "+s.id, func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		s.Open().Run(t, ctx)
		spec.Click(s.Option(value)).Run(t, ctx)
	})
}
func (s SelectControl) Opened() spec.Expectation { return spec.ExpectStep(spec.Visible(s.Content())) }
func (s SelectControl) ValueEquals(value Theme) spec.Expectation {
	return spec.TextContains(s.Value(), labelForTheme(value))
}
func (s SelectControl) Expanded(want bool) spec.Expectation {
	if want {
		return spec.AttributeEquals(s.Trigger(), spec.AriaExpanded, "true")
	}
	return spec.AttributeEquals(s.Trigger(), spec.AriaExpanded, "false")
}

func labelForTheme(value Theme) string {
	switch value {
	case ThemeMonokai:
		return "Monokai"
	case ThemeDracula:
		return "Dracula"
	case ThemeGitHub:
		return "GitHub"
	case ThemeNord:
		return "Nord"
	case ThemeNext:
		return "Next.js"
	default:
		return string(value)
	}
}
