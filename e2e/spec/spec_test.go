package spec

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/runtime"
)

func TestBuilderPreservesGivenWhenThenOrder(t *testing.T) {
	builder := Feature("Story API").
		Scenario("runs steps").
		Given(namedStep("given one"), namedStep("given two")).
		When(namedStep("when one")).
		Then(namedStep("then one"))

	steps := builder.orderedSteps()
	names := make([]string, 0, len(steps))
	for _, step := range steps {
		names = append(names, step.Name())
	}
	want := []string{"given one", "given two", "when one", "then one"}
	if len(names) != len(want) {
		t.Fatalf("names = %v, want %v", names, want)
	}
	for i := range names {
		if names[i] != want[i] {
			t.Fatalf("names = %v, want %v", names, want)
		}
	}
}

func TestBuilderSlugifiesFeatureAndScenarioNames(t *testing.T) {
	builder := Feature("Select Component!").Scenario("Opens options & checks console")
	if got, want := builder.scenario.FeatureSlug, "select-component"; got != want {
		t.Fatalf("feature slug = %q, want %q", got, want)
	}
	if got, want := builder.scenario.Slug, "opens-options-checks-console"; got != want {
		t.Fatalf("scenario slug = %q, want %q", got, want)
	}
}

func TestViewportsAreStoredOnScenario(t *testing.T) {
	builder := Feature("Viewport story").Scenario("mobile").Viewports(runtime.ViewportMobile, runtime.ViewportDesktopHalf)
	if got := builder.scenario.Viewports; len(got) != 2 || got[0] != runtime.ViewportMobile || got[1] != runtime.ViewportDesktopHalf {
		t.Fatalf("viewports = %v", got)
	}
}

func TestStoryBuilderAppendsDuckTypedStepsFlat(t *testing.T) {
	builder := Story(t, "flat story").
		As(fakeActor{"as"}).
		With(fakeFixture{"with"}).
		Visit(fakePage{"visit"}).
		Expect(fakeExpectation{"expect"}).
		Do(namedStep("do"))

	steps := builder.orderedSteps()
	names := make([]string, 0, len(steps))
	for _, step := range steps {
		names = append(names, step.Name())
	}
	want := []string{"as", "with", "visit", "expect", "do"}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("names = %v, want %v", names, want)
		}
	}
}

func TestStoryBuilderViewportStoresExplicitViewport(t *testing.T) {
	builder := Story(t, "flat viewport").Viewports(runtime.ViewportMobile, runtime.ViewportDesktopHalf)
	if got := builder.scenario.Viewports; len(got) != 2 || got[0] != runtime.ViewportMobile || got[1] != runtime.ViewportDesktopHalf {
		t.Fatalf("viewports = %v", got)
	}
}

type fakeActor struct{ name string }

func (f fakeActor) AuthStep() Step { return namedStep(f.name) }

type fakeFixture struct{ name string }

func (f fakeFixture) SetupStep() Step { return namedStep(f.name) }

type fakePage struct{ name string }

func (f fakePage) VisitStep() Step { return namedStep(f.name) }

type fakeExpectation struct{ name string }

func (f fakeExpectation) CheckStep() Step { return namedStep(f.name) }

func namedStep(name string) Step {
	return Custom(name, func(testing.TB, *runtime.Context) {})
}
