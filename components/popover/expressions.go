package popover

import (
	"fmt"
)

// PopoverHandler builds expressions for popover behavior
type PopoverHandler struct {
	popoverID string
}

// NewPopoverHandler creates a new popover expression handler
func NewPopoverHandler(popoverID string) *PopoverHandler {
	return &PopoverHandler{
		popoverID: popoverID,
	}
}

// BuildToggleHandler creates a popover toggle expression
func (p *PopoverHandler) BuildToggleHandler() string {
	return fmt.Sprintf("document.getElementById('%s').togglePopover()", p.popoverID)
}

// BuildAnchorStyle creates anchor positioning style
func (p *PopoverHandler) BuildAnchorStyle(anchorName string) string {
	if anchorName == "" {
		return ""
	}
	return fmt.Sprintf("anchor-name: --%s", anchorName)
}

// BuildPositionAnchorStyle creates position anchor style
func (p *PopoverHandler) BuildPositionAnchorStyle(anchorName string) string {
	if anchorName == "" {
		return ""
	}
	return fmt.Sprintf("position-anchor: --%s", anchorName)
}