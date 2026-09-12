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

	spec.Story(t, "infinite scroll loads 2 pages then exhausts").
		Visit(InfiniteScrollPage()).
		// Verify initial state
		Expect(HostExists(host)).
		Expect(SentinelExists(sentinel)).
		Expect(SentinelHasIntersectAttribute("diff_viewer-sentinel-below")).
		// Count initial items (≥1)
		Expect(InitialItemCount(itemsContainer, 1)).
		// Scroll to trigger first load (cursor=0)
		Expect(ScrollToSentinel(host, sentinel)).
		// Wait for first page to load and verify (≥2 items)
		Expect(ItemCountIncreased(itemsContainer, 2)).
		Expect(HostStillStable(host)).
		Expect(SentinelExists(sentinel)). // Sentinel should be reminted
		// Scroll to trigger second load (cursor=1, last content page)
		Expect(ScrollToSentinel(host, sentinel)).
		// Wait for second page to load and verify (≥3 items)
		Expect(ItemCountIncreased(itemsContainer, 3)).
		Expect(HostStillStable(host)).
		// Verify sentinel is gone (exhausted on last content page)
		Expect(SentinelGone(sentinel)).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}
