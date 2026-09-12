package toast_test

import (
	"github.com/coreycole/datastarui/e2e/spec"
)

type ToastDemo struct {
	id string
}

func ToastPage() spec.Page            { return spec.Path("/components/toast") }
func DefaultToast() ToastDemo         { return ToastDemo{id: "default_toast"} }
func SuccessToast() ToastDemo         { return ToastDemo{id: "success_toast"} }
func ErrorToast() ToastDemo           { return ToastDemo{id: "error_toast"} }
func PersistentToast() ToastDemo      { return ToastDemo{id: "persistent_toast"} }

func (t ToastDemo) Trigger() spec.Locator {
	// Find the button that triggers this toast by looking for the data-on:click attribute
	return spec.CSS("button[data-on\\:click*='$" + t.id + ".open']")
}

func (t ToastDemo) Content() spec.Locator {
	// The actual toast content div with the ID
	return spec.CSS("#" + t.id)
}

func (t ToastDemo) Container() spec.Locator {
	// The outer div with data-signals and data-show
	return spec.CSS("div[data-signals*='\"" + t.id + "\"'][data-show='$" + t.id + ".open']")
}

func (t ToastDemo) CloseButton() spec.Locator {
	// The close button inside the toast
	return spec.CSS("#" + t.id + " button[data-on\\:click*='false']")
}

func (t ToastDemo) Show() spec.Step {
	return spec.Click(t.Trigger())
}

func (t ToastDemo) Close() spec.Step {
	return spec.Click(t.CloseButton())
}

func (t ToastDemo) Visible() spec.Expectation {
	return spec.ExpectStep(spec.Visible(t.Content()))
}

func (t ToastDemo) Hidden() spec.Expectation {
	return spec.ExpectStep(spec.Hidden(t.Content()))
}
