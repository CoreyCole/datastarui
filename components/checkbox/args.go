package checkbox

import (
	"github.com/a-h/templ"
)

// CheckboxProps defines the properties for the Checkbox component
type CheckboxProps struct {
	// Core HTML attributes
	ID       string
	Name     string
	Value    string
	Checked  bool
	Disabled bool
	Required bool

	// Styling
	ClassName string

	// Accessibility
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	AriaInvalid     bool

	// Additional HTML attributes
	Attributes templ.Attributes
}

// CheckboxIndicatorProps defines the properties for the CheckboxIndicator component
type CheckboxIndicatorProps struct {
	// Styling
	ClassName string

	// Additional HTML attributes
	Attributes templ.Attributes
}
