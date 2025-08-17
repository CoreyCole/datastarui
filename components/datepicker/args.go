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

	// OpenOnFocus automatically opens the calendar when input is focused
	OpenOnFocus bool

	// DisablePopoverOpenOnFocus prevents the popover from opening on input focus
	// When true, popover only opens when calendar icon is clicked
	// When false (default), popover opens on input focus
	DisablePopoverOpenOnFocus bool

	// PopoverPosition controls where the popover appears relative to the input
	// Options: "bottom" (default), "top"
	PopoverPosition string

	// Class allows additional CSS classes to be added to the container
	Class string

	// InputClass allows additional CSS classes to be added to the input field
	InputClass string

	// CalendarClass allows additional CSS classes to be added to the calendar
	CalendarClass string

	// Note: CalendarID and DatePickerInputsID are now auto-generated internally
	// Signal coordination is handled automatically using the main ID

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// DatePickerCalendarProps defines properties for the calendar part of the date picker
type DatePickerCalendarProps struct {
	// Embed calendar args for full calendar functionality
	calendar.CalendarProps
}

// DatePickerPopoverProps defines properties for the popover calendar
type DatePickerPopoverProps struct {
	ID                  string
	DateInputID         string
	Mode                string
	NumberOfMonths      int
	MinDate             string
	MaxDate             string
	DisabledDates       string
	HideOutsideDays     bool
	CloseOnSelect       bool
	Class               string
	Position            string
	InitialSelectedDate string
	InitialDisplayMonth string
	InitialRangeStart   string
	InitialRangeEnd     string
	DatePickerInputsID  string // ID of date picker inputs to send signals to for synchronization
}
