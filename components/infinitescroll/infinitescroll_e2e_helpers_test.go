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
	return spec.ExpectStep(spec.Custom("host exists with stable ID: "+demo.hostID, func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		locator := ctx.Page.Locator("#" + demo.hostID)
		if err := locator.WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(5000),
		}); err != nil {
			t.Errorf("host element #%s not found or not visible: %v", demo.hostID, err)
		}
	}))
}

func SentinelHasIntersectAttribute(sentinelID string) spec.Expectation {
	return spec.ExpectStep(spec.Custom("sentinel has data-on:intersect attribute: "+sentinelID, func(t testing.TB, ctx *runtime.Context) {
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
	}))
}

type ItemsContainer struct {
	itemsID string
}

func DiffViewerItems() ItemsContainer {
	return ItemsContainer{itemsID: "diff_viewer-items"}
}

type Sentinel struct {
	sentinelID string
}

func DiffViewerSentinel() Sentinel {
	return Sentinel{sentinelID: "diff_viewer-sentinel-below"}
}

func SentinelExists(sentinel Sentinel) spec.Expectation {
	return spec.ExpectStep(spec.Custom("sentinel exists: "+sentinel.sentinelID, func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		locator := ctx.Page.Locator("#" + sentinel.sentinelID)
		if err := locator.WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateAttached,
			Timeout: playwright.Float(5000),
		}); err != nil {
			t.Errorf("sentinel #%s not found: %v", sentinel.sentinelID, err)
		}
	}))
}

func InitialItemCount(items ItemsContainer, expectedMin int) spec.Expectation {
	return spec.ExpectStep(spec.Custom("initial item count >= "+string(rune(expectedMin+'0')), func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		locator := ctx.Page.Locator("#" + items.itemsID + " > div")
		count, err := locator.Count()
		if err != nil {
			t.Errorf("failed to count items in #%s: %v", items.itemsID, err)
			return
		}
		if count < expectedMin {
			t.Errorf("expected at least %d items, got %d", expectedMin, count)
			return
		}
		t.Logf("✓ initial item count: %d", count)
	}))
}

func ScrollToSentinel(demo InfiniteScrollDemo, sentinel Sentinel) spec.Expectation {
	return spec.ExpectStep(spec.Custom("scroll to sentinel: "+sentinel.sentinelID, func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		sentinelLocator := ctx.Page.Locator("#" + sentinel.sentinelID)
		
		// Scroll the sentinel into view
		if err := sentinelLocator.ScrollIntoViewIfNeeded(); err != nil {
			t.Errorf("failed to scroll sentinel into view: %v", err)
			return
		}
		
		// Wait a bit for intersection observer to trigger
		ctx.Page.WaitForTimeout(500)
		
		// Wait for loading state or new content
		ctx.Page.WaitForTimeout(1000)
		
		t.Logf("✓ scrolled to sentinel %s", sentinel.sentinelID)
	}))
}

func ItemCountIncreased(items ItemsContainer, expectedMin int) spec.Expectation {
	return spec.ExpectStep(spec.Custom("item count increased to >= "+string(rune(expectedMin+'0')), func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		locator := ctx.Page.Locator("#" + items.itemsID + " > div")
		
		// Wait for items to increase
		ctx.Page.WaitForTimeout(1000)
		
		count, err := locator.Count()
		if err != nil {
			t.Errorf("failed to count items in #%s: %v", items.itemsID, err)
			return
		}
		if count < expectedMin {
			t.Errorf("expected at least %d items after load, got %d", expectedMin, count)
			return
		}
		t.Logf("✓ item count after load: %d", count)
	}))
}

func HostStillStable(demo InfiniteScrollDemo) spec.Expectation {
	return spec.ExpectStep(spec.Custom("host still stable: "+demo.hostID, func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		locator := ctx.Page.Locator("#" + demo.hostID)
		
		// Verify host still exists
		count, err := locator.Count()
		if err != nil || count != 1 {
			t.Errorf("host #%s is not stable (count: %d, err: %v)", demo.hostID, count, err)
			return
		}
		
		t.Logf("✓ host %s remains stable", demo.hostID)
	}))
}

func SentinelGone(sentinel Sentinel) spec.Expectation {
	return spec.ExpectStep(spec.Custom("sentinel gone (exhausted): "+sentinel.sentinelID, func(t testing.TB, ctx *runtime.Context) {
		t.Helper()
		locator := ctx.Page.Locator("#" + sentinel.sentinelID)
		
		// Wait for sentinel to be removed
		ctx.Page.WaitForTimeout(1000)
		
		count, err := locator.Count()
		if err != nil {
			t.Errorf("failed to check sentinel count: %v", err)
			return
		}
		if count > 0 {
			t.Errorf("expected sentinel #%s to be gone (exhausted), but it still exists", sentinel.sentinelID)
			return
		}
		
		t.Logf("✓ sentinel %s is gone (exhausted)", sentinel.sentinelID)
	}))
}
