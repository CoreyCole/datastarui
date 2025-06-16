package calendar

import "github.com/a-h/templ"

// CalendarProps defines the properties for the Calendar component
type CalendarProps struct {
	// ID for the calendar (auto-generated if not provided)
	ID string

	// Mode defines the selection mode
	// Options: "single", "range"
	Mode string

	// DefaultDate sets the initial month to display (YYYY-MM-DD format)
	// Will be normalized to the first day of the month
	DefaultDate string

	// SelectedDate for single mode (YYYY-MM-DD format)
	SelectedDate string

	// RangeStart for range mode (YYYY-MM-DD format)
	RangeStart string

	// RangeEnd for range mode (YYYY-MM-DD format)
	RangeEnd string

	// Disabled dates (comma-separated YYYY-MM-DD format)
	Disabled string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// CalendarGridProps defines the properties for the calendar day grid
type CalendarGridProps struct {
	// ID of the parent calendar
	ID string
}

// CalendarHeaderProps defines the properties for the calendar header
type CalendarHeaderProps struct {
	// ID of the parent calendar
	ID string
}

// CalendarDayProps defines the properties for individual calendar day buttons
type CalendarDayProps struct {
	// ID of the parent calendar
	ID string

	// DayIndex is the index of this day (0-41 for 6 weeks)
	DayIndex int
}
