package runtime

import "testing"

func TestScenarioViewportsUsesExplicitViewport(t *testing.T) {
	viewports, err := scenarioViewports(Config{Viewports: []ViewportClass{ViewportMobile}}, ViewportDesktopFull)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := viewports, []ViewportClass{ViewportDesktopFull}; len(got) != 1 || got[0] != want[0] {
		t.Fatalf("viewports = %v, want %v", got, want)
	}
}

func TestScenarioViewportsUsesConfiguredViewports(t *testing.T) {
	viewports, err := scenarioViewports(Config{Viewports: []ViewportClass{ViewportMobile, ViewportDesktopHalf}}, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(viewports) != 2 || viewports[0] != ViewportMobile || viewports[1] != ViewportDesktopHalf {
		t.Fatalf("viewports = %v", viewports)
	}
}

func TestScenarioViewportsRejectsUnknownViewport(t *testing.T) {
	if _, err := scenarioViewports(Config{}, ViewportClass("watch")); err == nil {
		t.Fatal("expected unknown viewport error")
	}
}
