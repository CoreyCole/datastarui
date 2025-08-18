package utils

import "strconv"

// AnchorSide defines the positioning side for anchor-positioned elements
type AnchorSide string

const (
	AnchorSideTop    AnchorSide = "top"
	AnchorSideBottom AnchorSide = "bottom"
	AnchorSideLeft   AnchorSide = "left"
	AnchorSideRight  AnchorSide = "right"
)

// AnchorAlign defines the alignment for anchor-positioned elements
type AnchorAlign string

const (
	AnchorAlignStart  AnchorAlign = "start"
	AnchorAlignCenter AnchorAlign = "center"
	AnchorAlignEnd    AnchorAlign = "end"
)

// GetAnchorPosition generates CSS anchor positioning based on side and align
// This function contains the proven working positioning logic from popover
func GetAnchorPosition(side AnchorSide, align AnchorAlign, sideOffset int) string {
	offset := strconv.Itoa(sideOffset) + "px"

	switch side {
	case AnchorSideTop:
		// Element appears above the anchor
		switch align {
		case AnchorAlignStart:
			return "top: anchor(top); left: anchor(left); translate: 0 calc(-100% - " + offset + ")"
		case AnchorAlignEnd:
			return "top: anchor(top); left: anchor(right); translate: -100% calc(-100% - " + offset + ")"
		default: // center
			return "top: anchor(top); left: anchor(center); translate: -50% calc(-100% - " + offset + ")"
		}
	case AnchorSideBottom:
		// Element appears below the anchor
		switch align {
		case AnchorAlignStart:
			return "top: anchor(bottom); left: anchor(left); translate: 0 " + offset
		case AnchorAlignEnd:
			return "top: anchor(bottom); left: anchor(right); translate: -100% " + offset
		default: // center
			return "top: anchor(bottom); left: anchor(center); translate: -50% " + offset
		}
	case AnchorSideRight:
		// Element appears to the right of the anchor
		switch align {
		case AnchorAlignStart:
			return "top: anchor(top); left: anchor(right); translate: " + offset + " 0"
		case AnchorAlignEnd:
			return "top: anchor(bottom); left: anchor(right); translate: " + offset + " -100%"
		default: // center
			return "top: anchor(center); left: anchor(right); translate: " + offset + " -50%"
		}
	case AnchorSideLeft:
		// Element appears to the left of the anchor
		switch align {
		case AnchorAlignStart:
			return "top: anchor(top); left: anchor(left); translate: calc(-100% - " + offset + ") 0"
		case AnchorAlignEnd:
			return "top: anchor(bottom); left: anchor(left); translate: calc(-100% - " + offset + ") -100%"
		default: // center
			return "top: anchor(center); left: anchor(left); translate: calc(-100% - " + offset + ") -50%"
		}
	default:
		// Default to bottom center
		return "top: anchor(bottom); left: anchor(center); translate: -50% " + offset
	}
}