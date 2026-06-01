package dropdown_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/spec"
)

func TestDropdownSelectOptionsNotClipped(t *testing.T) {
	dropdown := SelectOverflowDropdown()
	selectControl := DropdownSelect()

	spec.Story(t, "dropdown select options not clipped").
		Visit(DropdownPage()).
		Do(dropdown.Open()).
		Expect(dropdown.Opened()).
		Do(selectControl.Open()).
		Expect(selectControl.Opened()).
		Expect(spec.NotClippedBy(selectControl.Content(), dropdown.Content())).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestMultipleDropdownsOpenIndependently(t *testing.T) {
	first := NumberedDropdown(1)
	second := NumberedDropdown(2)
	third := NumberedDropdown(3)

	spec.Story(t, "multiple dropdowns open independently").
		Visit(DropdownPage()).
		Do(first.Open()).
		Expect(first.Opened()).
		Expect(second.Closed()).
		Expect(third.Closed()).
		Do(second.Open()).
		Expect(first.Closed()).
		Expect(second.Opened()).
		Expect(third.Closed()).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestNestedSelectDoesNotCloseDropdown(t *testing.T) {
	dropdown := UserProfileDropdown()
	selectControl := NestedThemeSelect()

	spec.Story(t, "nested select does not close dropdown").
		Visit(DropdownPage()).
		Do(dropdown.Open()).
		Do(selectControl.Choose(ThemeDracula)).
		Expect(dropdown.Opened()).
		Expect(selectControl.ValueEquals(ThemeDracula)).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestManualCompositionSelectDoesNotCloseDropdown(t *testing.T) {
	dropdown := ManualCompositionDropdown()
	selectControl := ManualSyntaxSelect()

	spec.Story(t, "manual composition select does not close dropdown").
		Visit(DropdownPage()).
		Do(dropdown.Open()).
		Do(selectControl.Choose(ThemeNord)).
		Expect(dropdown.Opened()).
		Expect(selectControl.ValueEquals(ThemeNord)).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}
