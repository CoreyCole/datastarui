package e2e_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/spec"
)

func TestDatePickerComponent(t *testing.T) {
	spec.Feature("DatePicker component").
		Scenario("opens popover").
		Given(spec.OpenPage("DatePickerComponent")).
		When(spec.Click(spec.SelectorAlias("datepicker.single.trigger"))).
		Then(
			spec.Visible(spec.SelectorAlias("datepicker.single.popover")),
			spec.ConsoleClean(),
		).
		Run(t)
}
