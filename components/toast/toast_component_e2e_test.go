package toast_test

import (
	"testing"
	"time"

	"github.com/coreycole/datastarui/e2e/runtime"
	"github.com/coreycole/datastarui/e2e/spec"
)

func TestToastAppearsWhenTriggered(t *testing.T) {
	toast := DefaultToast()

	spec.Story(t, "toast appears when triggered").
		Visit(ToastPage()).
		Expect(toast.Hidden()).
		Do(toast.Show()).
		Expect(toast.Visible()).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestToastCanBeManuallyDismissed(t *testing.T) {
	toast := PersistentToast()

	spec.Story(t, "toast can be manually dismissed").
		Visit(ToastPage()).
		Do(toast.Show()).
		Expect(toast.Visible()).
		Do(toast.Close()).
		Expect(toast.Hidden()).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestToastAutoDismisses(t *testing.T) {
	toast := DefaultToast()

	spec.Story(t, "toast auto-dismisses after duration").
		Visit(ToastPage()).
		Do(toast.Show()).
		Expect(toast.Visible()).
		Do(spec.Custom("wait for auto-dismiss", func(t testing.TB, ctx *runtime.Context) {
			t.Helper()
			// Default toast has 3000ms duration, wait 3.5s to be safe
			time.Sleep(3500 * time.Millisecond)
		})).
		Expect(toast.Hidden()).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestMultipleToastsCanBeShown(t *testing.T) {
	defaultToast := DefaultToast()
	successToast := SuccessToast()
	errorToast := ErrorToast()

	spec.Story(t, "multiple toasts can be shown").
		Visit(ToastPage()).
		Do(defaultToast.Show()).
		Expect(defaultToast.Visible()).
		Do(successToast.Show()).
		Expect(successToast.Visible()).
		Do(errorToast.Show()).
		Expect(errorToast.Visible()).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}

func TestSuccessToastShowsCopiedMessage(t *testing.T) {
	toast := SuccessToast()

	spec.Story(t, "success toast shows copied to clipboard message").
		Visit(ToastPage()).
		Do(toast.Show()).
		Expect(toast.Visible()).
		Expect(spec.TextContains(toast.Content(), "Copied to clipboard")).
		Expect(spec.ExpectStep(spec.ConsoleClean())).
		Run()
}
