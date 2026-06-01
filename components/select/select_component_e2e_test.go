package selectcomponent_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/spec"
)

func SelectPage() spec.Page { return spec.Path("/components/select") }

func TestSelectComponent(t *testing.T) {
	selectControl := SliceThemeSelect()

	spec.Story(t, "select component opens options").
		Visit(SelectPage()).
		Do(selectControl.Open()).
		Expect(selectControl.Opened()).
		Expect(spec.ExpectStep(spec.Visible(selectControl.Option(string(ThemeDark))))).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestSelectTabNavigation(t *testing.T) {
	spec.Story(t, "select tab navigation follows form order").
		Visit(SelectPage()).
		Do(spec.Click(TabNameInput())).
		Expect(spec.Focused(TabNameInput())).
		Do(spec.PressPage(spec.KeyTab)).
		Expect(TabCategorySelect().Focused()).
		Do(spec.PressPage(spec.KeyTab)).
		Expect(TabPrioritySelect().Focused()).
		Do(spec.PressPage(spec.KeyTab)).
		Expect(spec.Focused(TabEmailInput())).
		Do(spec.PressPage(spec.KeyTab)).
		Expect(spec.Focused(SubmitButton())).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}
