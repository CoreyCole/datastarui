package dialog

import (
	"fmt"
	"github.com/coreycole/datastarui/utils"
)

// DialogHandler creates handlers for Dialog component functionality
type DialogHandler struct {
	signals *utils.SignalManager
}

// NewDialogHandler creates a dialog handler
func NewDialogHandler(signals *utils.SignalManager) *DialogHandler {
	return &DialogHandler{
		signals: signals,
	}
}

// BuildBackdropClickHandler creates the backdrop click handler for closing dialog
func (d *DialogHandler) BuildBackdropClickHandler() string {
	return d.signals.ConditionalAction("evt.target === evt.currentTarget", "open", "false")
}

// BuildEscapeHandler creates an escape key handler for closing dialog
func (d *DialogHandler) BuildEscapeHandler() string {
	condition := fmt.Sprintf("evt.key === 'Escape' && %s", d.signals.Signal("open"))
	return d.signals.ConditionalAction(condition, "open", "false")
}

// BuildCloseHandler creates a close handler with optional return value
func (d *DialogHandler) BuildCloseHandler(returnValue string) string {
	expr := utils.NewExpression().Statement(d.signals.Set("open", "false"))

	if returnValue != "" {
		expr.Statement(d.signals.SetString("returnValue", returnValue))
	}

	return expr.Build()
}
