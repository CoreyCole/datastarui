package form

import "github.com/a-h/templ"

type FormArgs struct {
	ID         string           // Form ID
	Action     string           // Form action URL
	Class      string           // Additional CSS classes
	Attributes templ.Attributes // Additional HTML attributes
}

type FormItemArgs struct {
	Class      string           // Additional CSS classes
	Attributes templ.Attributes // Additional HTML attributes
}

type FormLabelArgs struct {
	For        string           // HTML for attribute (links to input ID)
	HasError   bool             // Whether the field has an error
	Class      string           // Additional CSS classes
	Attributes templ.Attributes // Additional HTML attributes
}

type FormControlArgs struct {
	ID              string           // Control ID
	AriaDescribedBy string           // ARIA described by attribute
	AriaInvalid     bool             // Whether the control is invalid
	Attributes      templ.Attributes // Additional HTML attributes
}

type FormDescriptionArgs struct {
	ID         string           // Description ID
	Class      string           // Additional CSS classes
	Attributes templ.Attributes // Additional HTML attributes
}

type FormMessageArgs struct {
	ID         string           // Message ID
	Message    string           // Error message text
	Class      string           // Additional CSS classes
	Attributes templ.Attributes // Additional HTML attributes
}
