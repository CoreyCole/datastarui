package datepicker

import "github.com/a-h/templ"

// DatePickerProps defines the properties for the DatePicker component
type DatePickerProps struct {
	// ID for the datepicker (auto-generated if not provided)
	ID string

	// Name for form submission
	Name string

	// Mode defines the selection mode
	// Options: "single", "range"
	Mode string

	// Placeholder text for input field
	Placeholder string

	// DefaultDate sets the initial month to display (YYYY-MM-DD format)
	// Will be normalized to the first day of the month
	DefaultDate string

	// SelectedDate for single mode (YYYY-MM-DD format)
	SelectedDate string

	// RangeStart for range mode (YYYY-MM-DD format)
	RangeStart string

	// RangeEnd for range mode (YYYY-MM-DD format)
	RangeEnd string

	// Required field for form validation
	Required bool

	// Disabled makes the datepicker non-interactive
	Disabled bool

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}
