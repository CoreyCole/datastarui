package toast

import (
	"fmt"
	"github.com/coreycole/datastarui/utils"
)

// ToastHandler provides clean expression builders for toast interactions
type ToastHandler struct {
	signals *utils.SignalManager
}

// NewToastHandler creates a toast handler
func NewToastHandler(signals *utils.SignalManager) *ToastHandler {
	return &ToastHandler{signals: signals}
}

// buildDismissHandler creates a handler that removes a toast from the array
func (h *ToastHandler) buildDismissHandler(toastID string) string {
	return fmt.Sprintf("%s = %s.filter(t => t.id !== '%s')",
		h.signals.Signal("toasts"),
		h.signals.Signal("toasts"),
		toastID,
	)
}

// buildAddToastHandler creates a handler that adds a new toast
func (h *ToastHandler) buildAddToastHandler() string {
	// This will be called with a toast object passed as a parameter
	return fmt.Sprintf("%s = [...%s, toast]",
		h.signals.Signal("toasts"),
		h.signals.Signal("toasts"),
	)
}

// buildAutoRemoveHandler creates a handler that auto-removes a toast after a duration
func (h *ToastHandler) buildAutoRemoveHandler(toastID string, duration int) string {
	if duration <= 0 {
		return ""
	}

	dismissExpr := h.buildDismissHandler(toastID)
	return fmt.Sprintf("setTimeout(() => { %s }, %d)", dismissExpr, duration)
}
