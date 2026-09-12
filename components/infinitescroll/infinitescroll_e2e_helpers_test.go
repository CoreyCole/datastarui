package infinitescroll_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/runtime"
	"github.com/coreycole/datastarui/e2e/spec"
	"github.com/playwright-community/playwright-go"
)

type InfiniteScrollDemo struct {
	hostID string
}

func InfiniteScrollPage() spec.Page {
	return spec.Path("/components/infinitescroll")
}

func DiffViewerHost() InfiniteScrollDemo {
	return InfiniteScrollDemo{hostID: "diff_viewer-host"}
}

func (d InfiniteScrollDemo) Host() spec.Locator {
	return spec.CSS("#" + d.hostID)
}

func HostExists(demo InfiniteScrollDemo) spec.Expectation {
	return spec.Custom("host exists with stable ID: "+demo.hostID, func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		locator := ctx.Page.Locator("#" + demo.hostID)
		if err := locator.WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(5000),
		}); err != nil {
			t.Errorf("host element #%s not found or not visible: %v", demo.hostID, err)
		}
	})
}

func SentinelHasIntersectAttribute(sentinelID string) spec.Expectation {
	return spec.Custom("sentinel has data-on:intersect attribute: "+sentinelID, func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		locator := ctx.Page.Locator("#" + sentinelID)
		
		// Wait for sentinel to exist
		if err := locator.WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateAttached,
			Timeout: playwright.Float(5000),
		}); err != nil {
			t.Errorf("sentinel element #%s not found: %v", sentinelID, err)
			return
		}

		// Check for data-on:intersect attribute
		attr, err := locator.GetAttribute("data-on:intersect")
		if err != nil {
			t.Errorf("failed to get data-on:intersect attribute from #%s: %v", sentinelID, err)
			return
		}
		if attr == "" {
			t.Errorf("sentinel #%s is missing data-on:intersect attribute", sentinelID)
			return
		}
		
		t.Logf("✓ sentinel #%s has data-on:intersect attribute: %s", sentinelID, attr)
	})
}
