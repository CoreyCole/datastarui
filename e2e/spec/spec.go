package spec

import (
	"regexp"
	"strings"
	"testing"

	"github.com/coreycole/datastarui/e2e/runtime"
)

type FeatureSpec struct {
	Name      string
	Slug      string
	Scenarios []ScenarioSpec
	Options   FeatureOptions
}

type FeatureOptions struct {
	Tags []string
}

type ScenarioSpec struct {
	FeatureSlug string
	Name        string
	Slug        string
	GivenSteps  []Step
	WhenSteps   []Step
	ThenSteps   []Step
	Viewports   []runtime.ViewportClass
}

type Step interface {
	Name() string
	Run(testing.TB, *runtime.Context)
}

type StepFunc struct {
	Label string
	Fn    func(testing.TB, *runtime.Context)
}

type Actor interface{ AuthStep() Step }
type Fixture interface{ SetupStep() Step }
type Page interface{ VisitStep() Step }
type Expectation interface{ CheckStep() Step }

type StoryBuilder struct {
	t        testing.TB
	scenario ScenarioSpec
	app      runtime.App
	steps    []Step
}

func (s StepFunc) Name() string { return s.Label }
func (s StepFunc) Run(t testing.TB, ctx *runtime.Context) {
	t.Helper()
	s.Fn(t, ctx)
}

type FeatureBuilder struct {
	spec FeatureSpec
}

type ScenarioBuilder struct {
	feature  *FeatureBuilder
	scenario ScenarioSpec
}

type FeatureOption func(*FeatureSpec)
type ScenarioOption func(*ScenarioSpec)

func Story(t testing.TB, name string) *StoryBuilder {
	t.Helper()
	return &StoryBuilder{t: t, scenario: ScenarioSpec{FeatureSlug: slugify(name), Name: name, Slug: slugify(name)}}
}

func (b *StoryBuilder) App(app runtime.App) *StoryBuilder { b.app = app; return b }
func (b *StoryBuilder) As(actor Actor) *StoryBuilder {
	b.steps = append(b.steps, actor.AuthStep())
	return b
}
func (b *StoryBuilder) With(fixture Fixture) *StoryBuilder {
	b.steps = append(b.steps, fixture.SetupStep())
	return b
}
func (b *StoryBuilder) Visit(page Page) *StoryBuilder {
	b.steps = append(b.steps, page.VisitStep())
	return b
}
func (b *StoryBuilder) Expect(expectation Expectation) *StoryBuilder {
	b.steps = append(b.steps, expectation.CheckStep())
	return b
}
func (b *StoryBuilder) Do(step Step) *StoryBuilder { b.steps = append(b.steps, step); return b }
func (b *StoryBuilder) Viewport(viewport runtime.ViewportClass) *StoryBuilder {
	return b.Viewports(viewport)
}
func (b *StoryBuilder) Viewports(viewports ...runtime.ViewportClass) *StoryBuilder {
	b.scenario.Viewports = append([]runtime.ViewportClass{}, viewports...)
	return b
}

func (b *StoryBuilder) orderedSteps() []Step { return append([]Step{}, b.steps...) }

func (b *StoryBuilder) Run() {
	b.t.Helper()
	cfg := runtime.Config{App: b.app}
	run := func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		for _, step := range b.steps {
			ctx.Memory["datastarui.current_step"] = step.Name()
			t.Logf("E2E step: %s", step.Name())
			step.Run(t, ctx)
		}
		delete(ctx.Memory, "datastarui.current_step")
	}
	if len(b.scenario.Viewports) == 0 {
		if err := runtime.RunScenarioWithConfig(b.t, &cfg, b.scenario.FeatureSlug, b.scenario.Slug, "", run); err != nil {
			b.t.Fatal(err)
		}
		return
	}
	for _, viewport := range b.scenario.Viewports {
		if err := runtime.RunScenarioWithConfig(b.t, &cfg, b.scenario.FeatureSlug, b.scenario.Slug, viewport, run); err != nil {
			b.t.Fatal(err)
		}
	}
}

func Feature(name string, opts ...FeatureOption) *FeatureBuilder {
	spec := FeatureSpec{Name: name, Slug: slugify(name)}
	for _, opt := range opts {
		opt(&spec)
	}
	return &FeatureBuilder{spec: spec}
}

func (f *FeatureBuilder) Scenario(name string, opts ...ScenarioOption) *ScenarioBuilder {
	scenario := ScenarioSpec{FeatureSlug: f.spec.Slug, Name: name, Slug: slugify(name)}
	for _, opt := range opts {
		opt(&scenario)
	}
	f.spec.Scenarios = append(f.spec.Scenarios, scenario)
	return &ScenarioBuilder{feature: f, scenario: scenario}
}

func (s *ScenarioBuilder) Given(steps ...Step) *ScenarioBuilder {
	s.scenario.GivenSteps = append(s.scenario.GivenSteps, steps...)
	return s
}

func (s *ScenarioBuilder) When(steps ...Step) *ScenarioBuilder {
	s.scenario.WhenSteps = append(s.scenario.WhenSteps, steps...)
	return s
}

func (s *ScenarioBuilder) Then(steps ...Step) *ScenarioBuilder {
	s.scenario.ThenSteps = append(s.scenario.ThenSteps, steps...)
	return s
}

func (s *ScenarioBuilder) Viewport(viewport runtime.ViewportClass) *ScenarioBuilder {
	return s.Viewports(viewport)
}

func (s *ScenarioBuilder) Viewports(viewports ...runtime.ViewportClass) *ScenarioBuilder {
	s.scenario.Viewports = append([]runtime.ViewportClass{}, viewports...)
	return s
}

func (s *ScenarioBuilder) Run(t *testing.T) {
	t.Helper()
	if err := s.RunTB(t); err != nil {
		t.Fatal(err)
	}
}

func (s *ScenarioBuilder) RunTB(t testing.TB) error {
	t.Helper()
	if len(s.scenario.Viewports) == 0 {
		return runtime.RunScenarioWithConfig(t, nil, s.scenario.FeatureSlug, s.scenario.Slug, "", s.runSteps)
	}
	for _, viewport := range s.scenario.Viewports {
		if err := runtime.RunScenarioWithConfig(t, nil, s.scenario.FeatureSlug, s.scenario.Slug, viewport, s.runSteps); err != nil {
			return err
		}
	}
	return nil
}

func (s *ScenarioBuilder) orderedSteps() []Step {
	steps := make([]Step, 0, len(s.scenario.GivenSteps)+len(s.scenario.WhenSteps)+len(s.scenario.ThenSteps))
	steps = append(steps, s.scenario.GivenSteps...)
	steps = append(steps, s.scenario.WhenSteps...)
	steps = append(steps, s.scenario.ThenSteps...)
	return steps
}

func (s *ScenarioBuilder) runSteps(t testing.TB, ctx *runtime.Context) {
	t.Helper()
	for _, step := range s.orderedSteps() {
		ctx.Memory["datastarui.current_step"] = step.Name()
		t.Logf("E2E step: %s", step.Name())
		step.Run(t, ctx)
	}
	delete(ctx.Memory, "datastarui.current_step")
}

func Custom(label string, fn func(testing.TB, *runtime.Context)) Step {
	return StepFunc{Label: label, Fn: fn}
}

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = nonSlugChars.ReplaceAllString(value, "-")
	return strings.Trim(value, "-")
}
