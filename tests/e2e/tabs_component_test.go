package e2e_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/spec"
)

func TestTabsComponent(t *testing.T) {
	spec.Feature("Tabs component").
		Scenario("renders tab triggers and panels").
		Given(spec.OpenPage("TabsComponent")).
		Then(
			spec.Visible(spec.SelectorAlias("tabs.first_trigger")),
			spec.Visible(spec.SelectorAlias("tabs.content")),
			spec.ConsoleClean(),
		).
		Run(t)
}
