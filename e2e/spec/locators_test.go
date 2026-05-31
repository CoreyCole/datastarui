package spec

import (
	"fmt"
	"strings"
	"testing"

	"github.com/coreycole/datastarui/e2e/runtime"
)

func TestTestIDBuildsDatastarFriendlySelector(t *testing.T) {
	label := fmt.Sprint(TestID("select-trigger"))
	for _, want := range []string{"data-testid", "data-test-id", "data-e2e", "data-slot", "select-trigger"} {
		if !strings.Contains(label, want) {
			t.Fatalf("TestID label %q missing %q", label, want)
		}
	}
}

func TestSelectorAliasUnknownKeyReturnsError(t *testing.T) {
	locator := SelectorAlias("missing")
	_, err := locator.Resolve(&runtime.Context{})
	if err == nil {
		t.Fatal("expected unknown selector error")
	}
}
