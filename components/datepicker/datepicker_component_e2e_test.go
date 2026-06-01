package datepicker_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/runtime"
	"github.com/coreycole/datastarui/e2e/spec"
)

func DatepickerPage() spec.Page { return spec.Path("/components/datepicker") }

func TestDatePickerComponent(t *testing.T) {
	spec.Story(t, "datepicker component opens popover").
		Visit(DatepickerPage()).
		Do(spec.Click(spec.CSS("[data-datepicker-id='single_date'] button[data-on-click*='open']"))).
		Expect(spec.ExpectStep(spec.Visible(spec.CSS("[data-datepicker-id='single_date'] [data-slot='datepicker-popover']")))).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestDatePickerInputFormatting(t *testing.T) {
	spec.Story(t, "datepicker single input auto-formats compact dates").
		Visit(DatepickerPage()).
		Do(fillDateInput("#single_date", "12252024")).
		Expect(spec.InputValue(spec.CSS("#single_date"), "12/25/2024")).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestDatePickerYearCompletion(t *testing.T) {
	spec.Story(t, "datepicker single input completes two digit year on blur").
		Visit(DatepickerPage()).
		Do(fillDateInput("#single_date", "12/25/24")).
		Do(spec.Click(spec.Role("heading", "Datepicker"))).
		Expect(spec.InputValue(spec.CSS("#single_date"), "12/25/2024")).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestDatePickerRangeInputs(t *testing.T) {
	spec.Story(t, "datepicker range inputs format independently").
		Visit(DatepickerPage()).
		Do(spec.Click(spec.CSS("[data-datepicker-id='range_date'] button[data-on-click*='open']"))).
		Expect(spec.ExpectStep(spec.Visible(spec.CSS("[data-datepicker-id='range_date'] [data-slot='datepicker-popover']")))).
		Do(fillDateInput("#range_date_start", "01012025")).
		Do(fillDateInput("#range_date_end", "01152025")).
		Expect(spec.InputValue(spec.CSS("#range_date_start"), "01/01/2025")).
		Expect(spec.InputValue(spec.CSS("#range_date_end"), "01/15/2025")).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func fillDateInput(selector, value string) spec.Step {
	return spec.Custom("fill date input "+selector, func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		locator := ctx.Page.Locator(selector).First()
		if err := locator.Fill(""); err != nil {
			t.Fatal(err)
		}
		if err := locator.Fill(value); err != nil {
			t.Fatal(err)
		}
		ctx.Page.WaitForTimeout(400)
	})
}
