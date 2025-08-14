package input

import "github.com/a-h/templ"

type InputArgs struct {
	Type        string           // Input type (text, password, email, etc.)
	Class       string           // Additional CSS classes
	Placeholder string           // Placeholder text
	Value       string           // Input value
	Name        string           // Input name attribute
	ID          string           // Input ID attribute
	FormID      string           // Form ID for automatic data-bind (optional)
	Disabled    bool             // Whether input is disabled
	Required    bool             // Whether input is required
	Attributes  templ.Attributes // Additional HTML attributes
}
