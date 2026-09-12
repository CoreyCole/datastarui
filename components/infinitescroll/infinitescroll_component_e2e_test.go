package infinitescroll_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/spec"
)

func TestInfiniteScrollSentinelHasIntersectAttribute(t *testing.T) {
	spec.Story(t, "infinite scroll sentinel has data-on:intersect attribute").
		Visit(InfiniteScrollPage()).
		Expect(SentinelHasIntersectAttribute("diff_viewer-sentinel-below")).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestInfiniteScrollHostStable(t *testing.T) {
	host := DiffViewerHost()

	spec.Story(t, "infinite scroll host remains stable").
		Visit(InfiniteScrollPage()).
		Expect(HostExists(host)).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestInfiniteScrollLoadCycle(t *testing.T) {
	host := DiffViewerHost()
	itemsContainer := DiffViewerItems()
	sentinel := DiffViewerSentinel()

	spec.Story(t, "infinite scroll loads multiple pages then exhausts").
		Visit(InfiniteScrollPage()).
		// Verify initial state
		Expect(HostExists(host)).
		Expect(SentinelExists(sentinel)).
		Expect(SentinelHasIntersectAttribute("diff_viewer-sentinel-below")).
		// Count initial items
		Expect(InitialItemCount(itemsContainer, 1)). // 1 file card initially
		// Scroll to trigger first load
		Expect(ScrollToSentinel(host, sentinel)).
		// Wait for first page to load and verify
		Expect(ItemCountIncreased(itemsContainer, 2)). // Should have 2 cards now
		Expect(HostStillStable(host)).
		Expect(SentinelExists(sentinel)). // Sentinel should be reminted
		// Scroll to trigger second load
		Expect(ScrollToSentinel(host, sentinel)).
		// Wait for second page to load and verify
		Expect(ItemCountIncreased(itemsContainer, 3)). // Should have 3+ cards now
		Expect(HostStillStable(host)).
		Expect(SentinelExists(sentinel)). // Sentinel should be reminted again
		// Scroll to trigger exhaustion
		Expect(ScrollToSentinel(host, sentinel)).
		// Verify sentinel is gone (exhausted)
		Expect(SentinelGone(sentinel)).
		Expect(HostStillStable(host)).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}
