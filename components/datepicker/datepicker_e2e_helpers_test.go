package datepicker_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/runtime"
	"github.com/coreycole/datastarui/e2e/spec"
)

func SingleDateInput() spec.Locator { return spec.CSS("#single_date") }
func RangeStartInput() spec.Locator { return spec.CSS("#range_date_start") }
func RangeEndInput() spec.Locator   { return spec.CSS("#range_date_end") }
func SingleDatePopover() spec.Locator {
	return spec.CSS("[data-datepicker-id='single_date'] [data-slot='datepicker-popover']")
}
func RangeDatePopover() spec.Locator {
	return spec.CSS("[data-datepicker-id='range_date'] [data-slot='datepicker-popover']")
}

func OpenSingleDate() spec.Step {
	return spec.Click(spec.CSS("[data-datepicker-id='single_date'] button[data-on-click*='open']"))
}
func OpenRangeDate() spec.Step {
	return spec.Click(spec.CSS("[data-datepicker-id='range_date'] button[data-on-click*='open']"))
}

func FillDateInput(locator spec.Locator, value string) spec.Step {
	return spec.Custom("fill date input", func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		resolved, err := locator.Resolve(ctx)
		if err != nil {
			t.Fatal(err)
		}
		input := resolved.First()
		if err := input.Fill(""); err != nil {
			t.Fatal(err)
		}
		if err := input.Fill(value); err != nil {
			t.Fatal(err)
		}
		ctx.Page.WaitForTimeout(400)
	})
}
