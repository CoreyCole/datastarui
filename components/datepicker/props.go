package datepicker

import (
	"github.com/a-h/templ"
	"github.com/coreycole/datastarui/components/calendar"
)

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

	// NumberOfMonths defines how many months to display in the calendar
	// Default: 1, for range pickers often use 2
	NumberOfMonths int

	// HideOutsideDays determines whether to hide dates from adjacent months in calendar
	// Default: false (outside days are shown)
	HideOutsideDays bool

	// Disabled dates (comma-separated YYYY-MM-DD format)
	DisabledDates string

	// MinDate sets the minimum selectable date (YYYY-MM-DD format)
	MinDate string

	// MaxDate sets the maximum selectable date (YYYY-MM-DD format)
	MaxDate string

	// Required field for form validation
	Required bool

	// Disabled makes the datepicker non-interactive
	Disabled bool

	// Class allows additional CSS classes to be added to the container
	Class string

	// InputClass allows additional CSS classes to be added to the input field
	InputClass string

	// CalendarClass allows additional CSS classes to be added to the calendar
	CalendarClass string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// DatePickerSignals defines the reactive state for the DatePicker component
type DatePickerSignals struct {
	// Open state of the popover
	Open bool `json:"open"`

	// Input value (what user sees in the input field - YYYY/MM/DD format)
	InputValue string `json:"inputValue"`

	// Selected date in internal format (YYYY-MM-DD)
	SelectedDate string `json:"selectedDate"`

	// Range start date (YYYY-MM-DD format)
	RangeStart string `json:"rangeStart"`

	// Range end date (YYYY-MM-DD format)
	RangeEnd string `json:"rangeEnd"`

	// Current display month (YYYY-MM-DD format, first day of month)
	DisplayMonth string `json:"displayMonth"`

	// Validation state
	IsValid bool `json:"isValid"`

	// Error message for invalid dates
	ErrorMessage string `json:"errorMessage"`

	// Whether user is currently typing (affects calendar sync behavior)
	IsTyping bool `json:"isTyping"`

	// Focused date in calendar (for keyboard navigation)
	FocusedDate string `json:"focusedDate"`

	// Highlighted date (for keyboard navigation within calendar)
	HighlightedDate string `json:"highlightedDate"`
}

// DatePickerInputProps defines properties for the input part of the date picker
type DatePickerInputProps struct {
	// ID for the input field
	ID string

	// Name for form submission
	Name string

	// Placeholder text
	Placeholder string

	// Value for the input (YYYY/MM/DD display format)
	Value string

	// Disabled state
	Disabled bool

	// Required field
	Required bool

	// Class for styling
	Class string

	// Attributes for additional HTML attributes
	Attributes templ.Attributes
}

// DatePickerCalendarProps defines properties for the calendar part of the date picker
type DatePickerCalendarProps struct {
	// Embed calendar props for full calendar functionality
	calendar.CalendarProps
}
