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

	// AsChild renders the button as a child element (for composition)
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
