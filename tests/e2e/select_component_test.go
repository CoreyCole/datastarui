package e2e_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/spec"
)

func TestSelectComponent(t *testing.T) {
	spec.Feature("Select component").
		Scenario("opens options").
		Given(spec.OpenPage("SelectComponent")).
		When(spec.Click(spec.SelectorAlias("select.trigger"))).
		Then(
			spec.Visible(spec.SelectorAlias("select.content")),
			spec.Visible(spec.SelectorAlias("select.first_item")),
			spec.ConsoleClean(),
		).
		Run(t)
}
