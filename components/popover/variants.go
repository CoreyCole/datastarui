package popover

import (
	"strconv"

	"github.com/coreycole/datastarui/utils"
)

// PopoverContentVariants generates CSS classes for popover content
func PopoverContentVariants(props PopoverContentProps) string {
	// Base classes matching shadcn/ui popover with native popover API
	base := "w-72 rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-none"

	return utils.TwMerge(base, props.Class)
}

// GetAnchorPosition generates CSS anchor positioning based on side and align
func GetAnchorPosition(side, align string, sideOffset int) string {
	offset := strconv.Itoa(sideOffset) + "px"

	switch side {
	case "top":
		// Popover appears above the anchor
		switch align {
		case "start":
			return "top: anchor(top); left: anchor(left); translate: 0 calc(-100% - " + offset + ")"
		case "end":
			return "top: anchor(top); left: anchor(right); translate: -100% calc(-100% - " + offset + ")"
		default: // center
			return "top: anchor(top); left: anchor(center); translate: -50% calc(-100% - " + offset + ")"
		}
	case "bottom":
		// Popover appears below the anchor
		switch align {
		case "start":
			return "top: anchor(bottom); left: anchor(left); translate: 0 " + offset
		case "end":
			return "top: anchor(bottom); left: anchor(right); translate: -100% " + offset
		default: // center
			return "top: anchor(bottom); left: anchor(center); translate: -50% " + offset
		}
	case "right":
		// Popover appears to the right of the anchor
		switch align {
		case "start":
			return "top: anchor(top); left: anchor(right); translate: " + offset + " 0"
		case "end":
			return "top: anchor(bottom); left: anchor(right); translate: " + offset + " -100%"
		default: // center
			return "top: anchor(center); left: anchor(right); translate: " + offset + " -50%"
		}
	case "left":
		// Popover appears to the left of the anchor
		switch align {
		case "start":
			return "top: anchor(top); left: anchor(left); translate: calc(-100% - " + offset + ") 0"
		case "end":
			return "top: anchor(bottom); left: anchor(left); translate: calc(-100% - " + offset + ") -100%"
		default: // center
			return "top: anchor(center); left: anchor(left); translate: calc(-100% - " + offset + ") -50%"
		}
	default:
		// Default to bottom
		return "top: anchor(bottom); left: anchor(center); translate: -50% " + offset
	}
}
