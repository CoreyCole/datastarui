package popover

import (
	"github.com/a-h/templ"
	"github.com/coreycole/datastarui/utils"
)

// PopoverTriggerArgs defines the properties for the popover trigger
type PopoverTriggerArgs struct {
	ID         string // Optional: ID for the trigger element itself
	Class      string
	PopoverID  string // Required: ID of the popover content to control (for popovertarget)
	AnchorName string // Optional: CSS anchor name for positioning (when using anchor positioning)
	Attributes templ.Attributes
}

// PopoverContentArgs defines the properties for the popover content
type PopoverContentArgs struct {
	ID         string           // Required: Must match PopoverTriggerArgs.PopoverID
	Class      string
	UseAnchor  bool             // Whether to use CSS anchor positioning
	AnchorName string           // CSS anchor name to position relative to (should match trigger's AnchorName)
	Side       utils.AnchorSide // Positioning side
	Align      utils.AnchorAlign // Alignment
	SideOffset int              // Offset in pixels from the anchor (default: 8)
	Attributes templ.Attributes
}
