package dialog_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/spec"
)

func TestDialogClosesWithEscape(t *testing.T) {
	dialog := BasicDialog()

	spec.Story(t, "dialog closes with escape").
		Visit(DialogPage()).
		Do(dialog.Open()).
		Expect(dialog.Opened()).
		Do(dialog.Close(CloseByEscape)).
		Expect(dialog.Closed()).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestDialogClosesWithBackdrop(t *testing.T) {
	dialog := BasicDialog()

	spec.Story(t, "dialog closes with backdrop").
		Visit(DialogPage()).
		Do(dialog.Open()).
		Expect(dialog.Opened()).
		Do(dialog.Close(CloseByBackdrop)).
		Expect(dialog.Closed()).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}
