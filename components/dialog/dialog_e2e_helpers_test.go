package dialog_test

import (
	"testing"

	"github.com/coreycole/datastarui/e2e/runtime"
	"github.com/coreycole/datastarui/e2e/spec"
)

type CloseMethod int

const (
	CloseByEscape CloseMethod = iota
	CloseByBackdrop
	CloseByButton
)

type DialogDemo struct{ id string }

func DialogPage() spec.Page   { return spec.Path("/components/dialog") }
func BasicDialog() DialogDemo { return DialogDemo{id: "modal_demo"} }

func (d DialogDemo) Trigger() spec.Locator { return spec.Role("button", "Open Modal Dialog") }
func (d DialogDemo) Content() spec.Locator { return spec.CSS("#" + d.id) }
func (d DialogDemo) Backdrop() spec.Locator {
	return spec.CSS("[data-signals*='\"" + d.id + "\"'] > div[data-show='$" + d.id + ".open']")
}
func (d DialogDemo) Status() spec.Locator {
	return spec.CSS("span[data-text=\"$" + d.id + ".open ? 'Open' : 'Closed'\"]")
}
func (d DialogDemo) CloseButton() spec.Locator {
	return spec.CSS("#" + d.id + " [data-slot='dialog-footer'] button")
}
func (d DialogDemo) Open() spec.Step { return spec.Click(d.Trigger()) }
func (d DialogDemo) Close(method CloseMethod) spec.Step {
	switch method {
	case CloseByEscape:
		return spec.PressPage(spec.KeyEscape)
	case CloseByBackdrop:
		return spec.Custom("click backdrop for "+d.id, func(t testing.TB, ctx *runtime.Context) {
			t.Helper()
			if err := ctx.Page.Mouse().Click(10, 10); err != nil {
				t.Fatal(err)
			}
		})
	case CloseByButton:
		return spec.Click(d.CloseButton())
	default:
		return spec.PressPage(spec.KeyEscape)
	}
}
func (d DialogDemo) Opened() spec.Expectation { return spec.TextEquals(d.Status(), "Open") }
func (d DialogDemo) Closed() spec.Expectation { return spec.TextEquals(d.Status(), "Closed") }
