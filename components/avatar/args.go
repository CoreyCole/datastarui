package avatar

import "github.com/a-h/templ"

// AvatarArgs defines the properties for the Avatar root component
type AvatarArgs struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes

	// BackgroundColor sets the background color (e.g., "#3b82f6" or "hsl(210, 100%, 50%)")
	BackgroundColor string

	// TextColor sets the text color (defaults to white if not specified)
	TextColor string
}
