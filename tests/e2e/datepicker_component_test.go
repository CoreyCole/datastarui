package e2e_test

import (
	"testing"
	"time"

	"github.com/coreycole/datastarui/e2e/runtime"
	"github.com/coreycole/datastarui/e2e/spec"
)

func TestDatePickerComponent(t *testing.T) {
	spec.Story(t, "datepicker component opens popover").
		Visit(spec.Path("/components/datepicker")).
		Do(spec.Click(spec.CSS("[data-datepicker-id='single_date'] button[data-on-click*='open']"))).
		Expect(spec.ExpectStep(spec.Visible(spec.CSS("[data-datepicker-id='single_date'] [data-slot='datepicker-popover']")))).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestDatePickerInputFormatting(t *testing.T) {
	spec.Story(t, "datepicker single input auto-formats compact dates").
		Visit(spec.Path("/components/datepicker")).
		Do(fillDateInput("#single_date", "12252024")).
		Expect(inputValue("#single_date", "12/25/2024")).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestDatePickerYearCompletion(t *testing.T) {
	spec.Story(t, "datepicker single input completes two digit year on blur").
		Visit(spec.Path("/components/datepicker")).
		Do(fillDateInput("#single_date", "12/25/24")).
		Do(spec.Click(spec.Role("heading", "Datepicker"))).
		Expect(inputValue("#single_date", "12/25/2024")).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestDatePickerRangeInputs(t *testing.T) {
	spec.Story(t, "datepicker range inputs format independently").
		Visit(spec.Path("/components/datepicker")).
		Do(spec.Click(spec.CSS("[data-datepicker-id='range_date'] button[data-on-click*='open']"))).
		Expect(spec.ExpectStep(spec.Visible(spec.CSS("[data-datepicker-id='range_date'] [data-slot='datepicker-popover']")))).
		Do(fillDateInput("#range_date_start", "01012025")).
		Do(fillDateInput("#range_date_end", "01152025")).
		Expect(inputValue("#range_date_start", "01/01/2025")).
		Expect(inputValue("#range_date_end", "01/15/2025")).
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

func inputValue(selector, want string) spec.Expectation {
	return spec.ExpectStep(spec.Custom("input value "+selector+" = "+want, func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		assertEventually(t, func() (bool, string) {
			got, err := ctx.Page.Locator(selector).First().InputValue()
			if err != nil {
				return false, err.Error()
			}
			return got == want, got
		})
	}))
}

func assertEventually(t testing.TB, check func() (bool, string)) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	last := ""
	for time.Now().Before(deadline) {
		ok, got := check()
		if ok {
			return
		}
		last = got
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("condition not met before timeout; last value %q", last)
}
