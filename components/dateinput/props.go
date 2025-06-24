package dateinput

import "github.com/a-h/templ"

type DateInputProps struct {
	ID          string           // Required for signal namespacing
	Name        string           // For form submission
	Class       string           // Additional CSS classes
	Placeholder string           // Override default "MM/DD/YYYY"
	Value       string           // Initial value (MM/DD/YYYY format)
	Disabled    bool             // Disabled state
	Required    bool             // Required validation
	MinDate     string           // Minimum date (YYYY-MM-DD format)
	MaxDate     string           // Maximum date (YYYY-MM-DD format)
	Postfix     templ.Component  // Optional postfix component (e.g., calendar icon)
	Attributes  templ.Attributes // Additional HTML attributes
}

type DateInputSignals struct {
	InputValue   string `json:"inputValue"`   // "12/25/2023" - display format
	DateValue    string `json:"dateValue"`    // "2023-12-25" - ISO format for forms
	IsValid      bool   `json:"isValid"`      // validation state
	ErrorMessage string `json:"errorMessage"` // validation error message
}
