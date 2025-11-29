package selectcomponent

import "github.com/a-h/templ"

// SelectOptionArgs represents a single option in the select dropdown
type SelectOptionArgs struct {
	// Value is the value of this option
	Value string `json:"value"`

	// Label is the display text for this option
	Label string `json:"label"`

	// Disabled makes this option non-selectable
	Disabled bool `json:"disabled,omitempty"`

	// Group is the optional group name for this option
	Group string `json:"group,omitempty"`
}

// SelectArgs defines the properties for the main Select component
type SelectArgs struct {
	// ID is used for scoping datastar signals
	ID string

	// Open controls whether the select is open (controlled)
	Open bool

	// DefaultOpen sets the initial open state (uncontrolled)
	DefaultOpen bool

	// Value is the current selected value
	Value string

	// DefaultValue is the initial value when uncontrolled
	DefaultValue string

	// Options is a slice of options to render automatically
	Options []SelectOptionArgs

	// Name for form submission
	Name string

	// Disabled makes the select non-interactive
	Disabled bool

	// Required for form validation
	Required bool

	// Placeholder text to show when no value is selected
	Placeholder string

	// OnChange is an optional Datastar expression to execute when value changes
	// Example: "@post('/api/update')"
	OnChange string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// SelectTriggerArgs defines the properties for the SelectTrigger component
type SelectTriggerArgs struct {
	// ID must match the parent Select ID for signal scoping
	ID string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes

	// Disabled makes the trigger non-interactive
	Disabled bool
}

// SelectValueArgs defines the properties for the SelectValue component
type SelectValueArgs struct {
	// ID must match the parent Select ID for signal scoping
	ID string

	// Placeholder text to show when no value is selected
	Placeholder string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// SelectContentArgs defines the properties for the SelectContent component
type SelectContentArgs struct {
	// ID must match the parent Select ID for signal scoping
	ID string

	// Position determines how the content is positioned relative to the trigger
	// Options: "item-aligned", "popper"
	Position string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// SelectItemArgs defines the properties for the SelectItem component
type SelectItemArgs struct {
	// ID must match the parent Select ID for signal scoping
	ID string

	// Value is the value of this item
	Value string

	// Index is the position of this item in the list (for keyboard navigation)
	Index int

	// Disabled makes the item non-selectable
	Disabled bool

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// SelectLabelArgs defines the properties for the SelectLabel component
type SelectLabelArgs struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// SelectSeparatorArgs defines the properties for the SelectSeparator component
type SelectSeparatorArgs struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// SelectGroupArgs defines the properties for the SelectGroup component
type SelectGroupArgs struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}
