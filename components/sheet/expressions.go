package sheet

import (
	"fmt"
	"github.com/coreycole/datastarui/utils"
)

// SheetHandler creates handlers for Sheet component functionality
type SheetHandler struct {
	signals *utils.SignalManager
}

// NewSheetHandler creates a sheet handler
func NewSheetHandler(signals *utils.SignalManager) *SheetHandler {
	return &SheetHandler{
		signals: signals,
	}
}

// BuildBackdropClickHandler creates the backdrop click handler for closing modal sheets
func (s *SheetHandler) BuildBackdropClickHandler() string {
	// Only close if clicking directly on the backdrop (not on sheet content)
	return s.signals.ConditionalAction("evt.target === evt.currentTarget", "open", "false")
}

// BuildEscapeHandler creates an escape key handler for closing sheet
func (s *SheetHandler) BuildEscapeHandler() string {
	condition := fmt.Sprintf("evt.key === 'Escape' && %s", s.signals.Signal("open"))
	return s.signals.ConditionalAction(condition, "open", "false")
}

// BuildCloseHandler creates a close handler with optional return value
func (s *SheetHandler) BuildCloseHandler(returnValue string) string {
	expr := utils.NewExpression().Statement(s.signals.Set("open", "false"))
	
	if returnValue != "" {
		expr.Statement(s.signals.SetString("returnValue", returnValue))
	}
	
	return expr.Build()
}