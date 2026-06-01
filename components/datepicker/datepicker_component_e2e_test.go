package datepicker_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/spec"
)

func TestDatepickerComponent(t *testing.T) {
	spec.Story(t, "datepicker component opens popover").
		Visit(DatepickerPage()).
		Do(OpenSingleDate()).
		Expect(spec.ExpectStep(spec.Visible(SingleDatePopover()))).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestDatepickerInputFormatting(t *testing.T) {
	spec.Story(t, "datepicker single input auto-formats compact dates").
		Visit(DatepickerPage()).
		Do(FillDateInput(SingleDateInput(), "12252024")).
		Expect(spec.InputValue(SingleDateInput(), "12/25/2024")).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestDatepickerYearCompletion(t *testing.T) {
	spec.Story(t, "datepicker single input completes two digit year on blur").
		Visit(DatepickerPage()).
		Do(FillDateInput(SingleDateInput(), "12/25/24")).
		Do(spec.Click(spec.Role("heading", "Datepicker"))).
		Expect(spec.InputValue(SingleDateInput(), "12/25/2024")).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestDatepickerRangeInputs(t *testing.T) {
	spec.Story(t, "datepicker range inputs format independently").
		Visit(DatepickerPage()).
		Do(OpenRangeDate()).
		Expect(spec.ExpectStep(spec.Visible(RangeDatePopover()))).
		Do(FillDateInput(RangeStartInput(), "01012025")).
		Do(FillDateInput(RangeEndInput(), "01152025")).
		Expect(spec.InputValue(RangeStartInput(), "01/01/2025")).
		Expect(spec.InputValue(RangeEndInput(), "01/15/2025")).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}
