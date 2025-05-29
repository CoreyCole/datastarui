package popover

import "github.com/a-h/templ"

// PopoverProps defines the properties for the root popover container
type PopoverProps struct {
	ID    string
	Class string
}

// PopoverTriggerProps defines the properties for the popover trigger button
type PopoverTriggerProps struct {
	ID         string // Optional: ID for the trigger element itself
	Class      string
	PopoverID  string // Required: ID of the popover content to control (for popovertarget)
	AnchorName string // Optional: CSS anchor name for positioning (when using anchor positioning)
	AsChild    bool   // Whether to render as child element instead of button
	Attributes templ.Attributes
}

// PopoverContentProps defines the properties for the popover content
type PopoverContentProps struct {
	ID         string // Required: Must match PopoverTriggerProps.PopoverID
	Class      string
	UseAnchor  bool   // Whether to use CSS anchor positioning
	AnchorName string // CSS anchor name to position relative to (should match trigger's AnchorName)
	Side       string // Positioning side: "top", "right", "bottom", "left"
	Align      string // Alignment: "start", "center", "end"
	SideOffset int    // Offset in pixels from the anchor (default: 8)
	Attributes templ.Attributes
}
