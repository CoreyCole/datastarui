package tabs_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/spec"
)

func TabsPage() spec.Page { return spec.Path("/components/tabs") }

func TestTabsComponent(t *testing.T) {
	spec.Story(t, "tabs component renders tab triggers and panels").
		Visit(TabsPage()).
		Expect(spec.ExpectStep(spec.Visible(spec.CSS("[data-slot='tabs-trigger'][data-value], [data-slot='tabs-list'] button[role='tab']")))).
		Expect(spec.ExpectStep(spec.Visible(spec.CSS("[data-slot='tabs-content'][data-value]")))).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}
