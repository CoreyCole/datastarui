package tooltip

import (
	"github.com/a-h/templ"
	"github.com/coreycole/datastarui/utils"
)

// TooltipTriggerArgs defines the properties for the tooltip trigger (matching PopoverTriggerArgs pattern)
type TooltipTriggerArgs struct {
	ID            string // Optional: ID for the trigger element itself
	Class         string
	TooltipID     string // Required: ID of the tooltip content to control
	DelayDuration int    // Delay in milliseconds before showing tooltip (default 700)
	Attributes    templ.Attributes
}

// TooltipContentArgs defines the properties for the tooltip content (matching PopoverContentArgs pattern)
type TooltipContentArgs struct {
	ID         string            // Required: Must match TooltipTriggerArgs.TooltipID
	Class      string
	UseAnchor  bool              // Whether to use CSS anchor positioning
	Side       utils.AnchorSide  // Positioning side
	Align      utils.AnchorAlign // Alignment
	SideOffset int               // Offset in pixels from the anchor (default: 4)
	Attributes templ.Attributes
}


