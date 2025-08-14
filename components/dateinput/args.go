package dateinput

import "github.com/a-h/templ"

type DateInputArgs struct {
	// Basic args
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

	// Signal coordination
	CalendarID string // ID of calendar to send signals to for synchronization

	// Single mode args
	Value string // Initial value (MM/DD/YYYY format)

	// Range mode args
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
