package spec

import "testing"

func TestStepNamesAreReadable(t *testing.T) {
	cases := []struct {
		name string
		step Step
		want string
	}{
		{name: "click", step: Click(TestID("save")), want: `click css:[data-testid="save"], [data-test-id="save"], [data-e2e="save"], [data-slot="save"]`},
		{name: "press page", step: PressPage(KeyTab), want: "press page Tab"},
		{name: "input value", step: InputValue(CSS("#name"), "Pedro").CheckStep(), want: "input value css:#name = Pedro"},
		{name: "attribute", step: AttributeEquals(CSS("#menu"), AriaExpanded, "true").CheckStep(), want: "attribute aria-expanded equals css:#menu = true"},
		{name: "focused", step: Focused(CSS("#name")).CheckStep(), want: "focused css:#name"},
		{name: "visible", step: Visible(Text("Saved")), want: "visible text:Saved"},
		{name: "console", step: ConsoleClean(), want: "console clean"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.step.Name(); got != tc.want {
				t.Fatalf("Name() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestOpenPageStepNameIncludesKey(t *testing.T) {
	if got, want := OpenPage("SelectComponent").Name(), "open page SelectComponent"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}
