package form

import "github.com/a-h/templ"

type FormProps struct {
	ID         string           // Form ID
	Action     string           // Form action URL
	Method     string           // Form method (GET, POST, etc.)
	Class      string           // Additional CSS classes
	Attributes templ.Attributes // Additional HTML attributes
}

type FormItemProps struct {
	Class      string           // Additional CSS classes
	Attributes templ.Attributes // Additional HTML attributes
}

type FormLabelProps struct {
	For        string           // HTML for attribute (links to input ID)
	HasError   bool             // Whether the field has an error
	Class      string           // Additional CSS classes
	Attributes templ.Attributes // Additional HTML attributes
}

type FormControlProps struct {
	ID              string           // Control ID
	AriaDescribedBy string           // ARIA described by attribute
	AriaInvalid     bool             // Whether the control is invalid
	Attributes      templ.Attributes // Additional HTML attributes
}

type FormDescriptionProps struct {
	ID         string           // Description ID
	Class      string           // Additional CSS classes
	Attributes templ.Attributes // Additional HTML attributes
}

type FormMessageProps struct {
	ID         string           // Message ID
	Message    string           // Error message text
	Class      string           // Additional CSS classes
	Attributes templ.Attributes // Additional HTML attributes
}
