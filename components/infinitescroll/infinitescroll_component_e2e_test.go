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
