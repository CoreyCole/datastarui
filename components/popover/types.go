package popover

import "github.com/a-h/templ"

// PopoverProps defines the properties for the root popover container
type PopoverProps struct {
	ID    string
	Class string
}

// PopoverTriggerProps defines the properties for the popover trigger button
type PopoverTriggerProps struct {
	ID         string
	Class      string
	PopoverID  string // ID of the popover to control
	AsChild    bool
	Attributes templ.Attributes
}

// PopoverContentProps defines the properties for the popover content
type PopoverContentProps struct {
	ID         string
	Class      string
	Side       string // top, right, bottom, left
	Align      string // start, center, end
	SideOffset int    // offset in pixels
	PopoverAPI string // "auto", "manual", or "" for auto
	Attributes templ.Attributes
}
