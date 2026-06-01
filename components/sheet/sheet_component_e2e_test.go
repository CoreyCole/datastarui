package sheet_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/spec"
)

func TestModalSheetClosesWithEscape(t *testing.T) {
	sheet := RightModalSheet()

	spec.Story(t, "modal sheet closes with escape").
		Visit(SheetPage()).
		Do(sheet.Open()).
		Expect(sheet.Opened()).
		Do(sheet.Close(CloseByEscape)).
		Expect(sheet.Closed()).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestModalSheetClosesWithBackdrop(t *testing.T) {
	sheet := RightModalSheet()

	spec.Story(t, "modal sheet closes with backdrop").
		Visit(SheetPage()).
		Do(sheet.Open()).
		Expect(sheet.Opened()).
		Do(sheet.Close(CloseByBackdrop)).
		Expect(sheet.Closed()).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestNonModalSheetAllowsMainInteraction(t *testing.T) {
	sheet := NonModalSheet()

	spec.Story(t, "non modal sheet allows main interaction").
		Visit(SheetPage()).
		Do(sheet.Open()).
		Expect(sheet.Opened()).
		Expect(spec.ExpectStep(spec.Visible(spec.Role("heading", "Sheet")))).
		Do(sheet.Close(CloseByEscape)).
		Expect(sheet.Closed()).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}
