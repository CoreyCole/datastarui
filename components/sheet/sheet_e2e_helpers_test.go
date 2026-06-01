package sheet_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/runtime"
	"github.com/coreycole/datastarui/e2e/spec"
)

type CloseMethod int

type SheetMode int

type SheetSide string

const (
	CloseByEscape CloseMethod = iota
	CloseByBackdrop
	CloseByButton

	SheetModal SheetMode = iota
	SheetNonModal

	SheetRight  SheetSide = "right"
	SheetLeft   SheetSide = "left"
	SheetBottom SheetSide = "bottom"
)

type SheetDemo struct {
	id   string
	mode SheetMode
	side SheetSide
}

func SheetPage() spec.Page { return spec.Path("/components/sheet") }
func RightModalSheet() SheetDemo {
	return SheetDemo{id: "right_modal_sheet", mode: SheetModal, side: SheetRight}
}
func NonModalSheet() SheetDemo {
	return SheetDemo{id: "non_modal_sheet", mode: SheetNonModal, side: SheetBottom}
}

func (s SheetDemo) Trigger() spec.Locator {
	switch s.id {
	case "right_modal_sheet":
		return spec.Role("button", "Open Right Sheet")
	case "non_modal_sheet":
		return spec.Role("button", "Open Format Panel")
	default:
		return spec.CSS("[data-on\\:click*='$" + s.id + ".open']")
	}
}
func (s SheetDemo) Content() spec.Locator { return spec.CSS("#" + s.id) }
func (s SheetDemo) Backdrop() spec.Locator {
	return spec.CSS("[data-signals*='\"" + s.id + "\"'] > div[data-show='$" + s.id + ".open']")
}
func (s SheetDemo) CloseButton() spec.Locator {
	return spec.CSS("#" + s.id + " [data-slot='sheet-content'] button[data-on\\:click]")
}
func (s SheetDemo) Open() spec.Step { return spec.Click(s.Trigger()) }
func (s SheetDemo) Close(method CloseMethod) spec.Step {
	switch method {
	case CloseByEscape:
		return spec.PressPage(spec.KeyEscape)
	case CloseByBackdrop:
		return spec.Custom("click backdrop for "+s.id, func(t testing.TB, ctx *runtime.Context) {
			t.Helper()
			if err := ctx.Page.Mouse().Click(10, 10); err != nil {
				t.Fatal(err)
			}
		})
	case CloseByButton:
		return spec.Click(s.CloseButton())
	default:
		return spec.PressPage(spec.KeyEscape)
	}
}
func (s SheetDemo) Opened() spec.Expectation { return spec.ExpectStep(spec.Visible(s.Content())) }
func (s SheetDemo) Closed() spec.Expectation { return spec.ExpectStep(spec.Hidden(s.Content())) }
func (s SheetDemo) Modal() bool              { return s.mode == SheetModal }
func (s SheetDemo) Side() SheetSide          { return s.side }
