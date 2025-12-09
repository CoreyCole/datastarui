package textarea

import "github.com/a-h/templ"

type TextareaArgs struct {
	Class       string           // Additional CSS classes
	Placeholder string           // Placeholder text
	Value       string           // Textarea value
	Name        string           // Textarea name attribute
	ID          string           // Textarea ID attribute
	FormID      string           // Form ID for automatic data-bind (optional)
	Rows        int              // Number of rows (default: 4)
	Disabled    bool             // Whether textarea is disabled
	Required    bool             // Whether textarea is required
	Borderless  bool             // Whether to remove border and focus styles
	Attributes  templ.Attributes // Additional HTML attributes
}
