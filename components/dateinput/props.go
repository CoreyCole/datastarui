package dateinput

import "github.com/a-h/templ"

type DateInputProps struct {
	// Basic props
	ID          string           // Required for signal namespacing
	Name        string           // For form submission (single mode) or base name (range mode)
	Class       string           // Additional CSS classes
	Placeholder string           // Override default placeholder
	Disabled    bool             // Disabled state
	Required    bool             // Required validation
	MinDate     string           // Minimum date (YYYY-MM-DD format)
	MaxDate     string           // Maximum date (YYYY-MM-DD format)
	Postfix     templ.Component  // Optional postfix component (e.g., calendar icon)
	Attributes  templ.Attributes // Additional HTML attributes

	// Single mode props
	Value string // Initial value (MM/DD/YYYY format)

	// Range mode props
	Mode             string // "single" (default) or "range"
	StartValue       string // Initial start date value (MM/DD/YYYY format)
	EndValue         string // Initial end date value (MM/DD/YYYY format)
	StartName        string // Form field name for start date (overrides Name_start)
	EndName          string // Form field name for end date (overrides Name_end)
	StartPlaceholder string // Placeholder for start date input
	EndPlaceholder   string // Placeholder for end date input
	EndDateOptional  bool   // If true, shows "Set an end date" checkbox
	Separator        string // Text between inputs (default: "to")
	Orientation      string // "horizontal" (default) or "vertical"
}

type DateInputSignals struct {
	// Single mode signals
	InputValue string `json:"inputValue"` // "12/25/2023" - display format
	DateValue  string `json:"dateValue"`  // "2023-12-25" - ISO format for forms

	// Range mode signals
	StartInputValue string `json:"startInputValue"` // Start date display format
	StartDateValue  string `json:"startDateValue"`  // Start date ISO format
	EndInputValue   string `json:"endInputValue"`   // End date display format
	EndDateValue    string `json:"endDateValue"`    // End date ISO format
	EndDateEnabled  bool   `json:"endDateEnabled"`  // Checkbox state for optional end date

	// Shared validation
	IsValid      bool   `json:"isValid"`      // Overall validation state
	ErrorMessage string `json:"errorMessage"` // Validation error message
}
