package goldens

import (
	"context"
	"testing"
)

func TestAcceptRequiresHumanApproval(t *testing.T) {
	err := Accept(context.Background(), GoldenInput{RunPath: t.TempDir(), GoldenRoot: t.TempDir()})
	if err == nil {
		t.Fatal("expected approval error")
	}
}
