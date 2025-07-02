package dropdown

import (
	"fmt"

	"github.com/coreycole/datastarui/utils"
)

// DropdownHandler creates handlers for dropdown components
type DropdownHandler struct {
	signals *utils.SignalManager
}

// NewDropdownHandler creates a dropdown handler
func newDropdownHandler(signals *utils.SignalManager) *DropdownHandler {
	return &DropdownHandler{
		signals: signals,
	}
}

// BuildClickOutsideHandler creates a click outside handler for closing dropdown
func (d *DropdownHandler) buildClickOutsideHandler() string {
	return d.signals.ConditionalAction(d.signals.Signal("open"), "open", "false")
}

// BuildEscapeHandler creates an escape key handler for closing dropdown
func (d *DropdownHandler) buildEscapeHandler() string {
	condition := fmt.Sprintf("evt.key === 'Escape' && %s", d.signals.Signal("open"))
	return d.signals.ConditionalAction(condition, "open", "false")
}

// CreateSideClasses generates positioning classes for dropdown sides
func createSideClasses(side string, offset int) string {
	if offset == 0 {
		offset = 4 // Default offset like shadcn/ui
	}

	switch side {
	case "top":
		return fmt.Sprintf("bottom-full mb-%d", offset)
	case "bottom":
		return fmt.Sprintf("top-full mt-%d", offset)
	case "left":
		return fmt.Sprintf("right-full mr-%d", offset)
	case "right":
		return fmt.Sprintf("left-full ml-%d", offset)
	default:
		return fmt.Sprintf("top-full mt-%d", offset)
	}
}
