package tooltip

import (
	"fmt"
	"github.com/coreycole/datastarui/utils"
)

const defaultSideOffset = 4

// TooltipHandler provides methods for building tooltip-related expressions
type TooltipHandler struct {
	tooltipID string
}

// NewTooltipHandler creates a new tooltip handler
func NewTooltipHandler(tooltipID string) *TooltipHandler {
	return &TooltipHandler{
		tooltipID: tooltipID,
	}
}

// BuildShowHandler creates the hover/focus handler to show tooltip
func (h *TooltipHandler) BuildShowHandler(delayMs int) string {
	return fmt.Sprintf("setTimeout(() => { document.getElementById('%s').showPopover(); }, %d)", h.tooltipID, delayMs)
}

// BuildHideHandler creates the handler to hide tooltip
func (h *TooltipHandler) BuildHideHandler() string {
	return fmt.Sprintf("document.getElementById('%s').hidePopover()", h.tooltipID)
}

// BuildInstantShowHandler creates handler to show tooltip without delay
func (h *TooltipHandler) BuildInstantShowHandler() string {
	signals := utils.Signals(h.tooltipID, TooltipSignals{})
	
	// Clear timeouts and show immediately
	return utils.NewExpression().
		Statement("clearTimeout($" + h.tooltipID + ".hideTimeout)").
		Statement("clearTimeout($" + h.tooltipID + ".showTimeout)").
		Statement("$" + h.tooltipID + ".hideTimeout = null").
		Statement("$" + h.tooltipID + ".showTimeout = null").
		Statement(signals.Set("open", "true")).
		Build()
}

// BuildDelayedHideHandler creates handler to hide tooltip with delay
func (h *TooltipHandler) BuildDelayedHideHandler(delayMs int) string {
	// Clear show timeout and set hide timeout
	return utils.NewExpression().
		Statement("clearTimeout($" + h.tooltipID + ".showTimeout)").
		Statement("$" + h.tooltipID + ".showTimeout = null").
		Statement("$" + h.tooltipID + ".hideTimeout = setTimeout(() => { $" + h.tooltipID + ".open = false; }, " + fmt.Sprintf("%d", delayMs) + ")").
		Build()
}

// BuildAnchorStyle creates anchor positioning style (same as popover)
func (h *TooltipHandler) BuildAnchorStyle(anchorName string) string {
	if anchorName == "" {
		return ""
	}
	return fmt.Sprintf("anchor-name: --%s", anchorName)
}

// BuildPositionAnchorStyle creates position anchor style (same as popover)
func (h *TooltipHandler) BuildPositionAnchorStyle(anchorName string) string {
	if anchorName == "" {
		return ""
	}
	return fmt.Sprintf("position-anchor: --%s", anchorName)
}

// BuildTouchStartHandler creates handler for touch start (mobile touch-and-hold)
func (h *TooltipHandler) BuildTouchStartHandler(touchHoldMs int) string {
	signals := utils.Signals(h.tooltipID, TooltipSignals{})
	
	if touchHoldMs == 0 {
		touchHoldMs = 500 // Default touch-and-hold duration
	}
	
	return utils.NewExpression().
		Statement("evt.preventDefault()"). // Prevent text selection on long press
		Statement("clearTimeout($" + h.tooltipID + ".touchTimer)").
		Statement("$" + h.tooltipID + ".touchTimer = setTimeout(() => { " +
			"$" + h.tooltipID + ".touchHeld = true; " +
			"document.getElementById('" + h.tooltipID + "').showPopover(); " +
		"}, " + fmt.Sprintf("%d", touchHoldMs) + ")").
		Build()
}

// BuildTouchEndHandler creates handler for touch end (cancel if released early)
func (h *TooltipHandler) BuildTouchEndHandler() string {
	signals := utils.Signals(h.tooltipID, TooltipSignals{})
	
	return utils.NewExpression().
		Statement("clearTimeout($" + h.tooltipID + ".touchTimer)").
		Statement("$" + h.tooltipID + ".touchTimer = null").
		// Only hide if touch wasn't held long enough
		Conditional(
			"!$" + h.tooltipID + ".touchHeld",
			"document.getElementById('" + h.tooltipID + "').hidePopover()",
			"null",
		).
		Build()
}

// BuildClickOutsideHandler creates handler to dismiss tooltip when clicking outside (like popover)
func (h *TooltipHandler) BuildClickOutsideHandler(triggerSelector string) string {
	signals := utils.Signals(h.tooltipID, TooltipSignals{})
	
	// Close tooltip if clicking outside and it was opened via touch
	condition := "$" + h.tooltipID + ".touchHeld && !evt.target.closest('" + triggerSelector + "')"
	
	return utils.NewConditional(
		condition,
		utils.NewExpression().
			Statement("$" + h.tooltipID + ".touchHeld = false").
			Statement("document.getElementById('" + h.tooltipID + "').hidePopover()").
			Build(),
		"null",
	).Build()
}
