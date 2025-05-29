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
		switch align {
		case "start":
			return "bottom: anchor(top); left: anchor(left); translate: 0 -" + offset
		case "end":
			return "bottom: anchor(top); right: anchor(right); translate: 0 -" + offset
		default: // center
			return "bottom: anchor(top); left: anchor(center); translate: -50% -" + offset
		}
	case "right":
		switch align {
		case "start":
			return "left: anchor(right); top: anchor(top); translate: " + offset + " 0"
		case "end":
			return "left: anchor(right); bottom: anchor(bottom); translate: " + offset + " 0"
		default: // center
			return "left: anchor(right); top: anchor(center); translate: " + offset + " -50%"
		}
	case "left":
		switch align {
		case "start":
			return "right: anchor(left); top: anchor(top); translate: -" + offset + " 0"
		case "end":
			return "right: anchor(left); bottom: anchor(bottom); translate: -" + offset + " 0"
		default: // center
			return "right: anchor(left); top: anchor(center); translate: -" + offset + " -50%"
		}
	default: // bottom
		switch align {
		case "start":
			return "top: anchor(bottom); left: anchor(left); translate: 0 " + offset
		case "end":
			return "top: anchor(bottom); right: anchor(right); translate: 0 " + offset
		default: // center
			return "top: anchor(bottom); left: anchor(center); translate: -50% " + offset
		}
	}
}
