package e2e_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/spec"
)

func TestSelectComponent(t *testing.T) {
	spec.Story(t, "select component opens options").
		Visit(spec.Path("/components/select")).
		Do(spec.Click(spec.CSS("[data-select-id='slice_theme_select'] [data-slot='select-trigger']"))).
		Expect(spec.ExpectStep(spec.Visible(spec.CSS("[data-select-id='slice_theme_select'] [data-slot='select-content']")))).
		Expect(spec.ExpectStep(spec.Visible(spec.CSS("[data-select-id='slice_theme_select'] [data-select-item]")))).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}
