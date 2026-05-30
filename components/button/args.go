package button

import "github.com/a-h/templ"

// ButtonArgs defines the properties for the Button component
type ButtonArgs struct {
	// Variant defines the visual style of the button
	// Options: "default", "destructive", "outline", "secondary", "ghost", "link"
	Variant string

	// Size defines the size of the button
	// Options: "default", "sm", "lg", "icon"
	Size string

	// AsChild renders wrapper-styled children for compatibility only.
	// It does not transfer attributes into the child element.
	AsChild bool

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes

	// Disabled makes the button non-interactive
	Disabled bool

	// Type specifies the button type (button, submit, reset)
	Type string
}

// LinkButtonArgs defines the properties for an anchor rendered with button styling.
type LinkButtonArgs struct {
	// Href is the destination URL for the link.
	Href string

	// Variant defines the visual style of the button link.
	// Options: "default", "destructive", "outline", "secondary", "ghost", "link"
	Variant string

	// Size defines the size of the button link.
	// Options: "default", "sm", "lg", "icon"
	Size string

	// Target specifies where to open the linked document.
	Target string

	// Rel specifies the relationship of the target object to the link object.
	Rel string

	// Class allows additional CSS classes to be added.
	Class string

	// Attributes allows additional HTML attributes to be added.
	Attributes templ.Attributes
}
