package selectcomponent

import "github.com/a-h/templ"

// SelectOption represents a single option in the select dropdown
type SelectOption struct {
	// Value is the value of this option
	Value string `json:"value"`

	// Label is the display text for this option
	Label string `json:"label"`

	// Disabled makes this option non-selectable
	Disabled bool `json:"disabled,omitempty"`

	// Group is the optional group name for this option
	Group string `json:"group,omitempty"`
}

// SelectProps defines the properties for the main Select component
type SelectProps struct {
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
	Options []SelectOption

	// Name for form submission
	Name string

	// Disabled makes the select non-interactive
	Disabled bool

	// Required for form validation
	Required bool

	// Placeholder text to show when no value is selected
	Placeholder string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// SelectTriggerProps defines the properties for the SelectTrigger component
type SelectTriggerProps struct {
	// ID must match the parent Select ID for signal scoping
	ID string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes

	// Disabled makes the trigger non-interactive
	Disabled bool
}

// SelectValueProps defines the properties for the SelectValue component
type SelectValueProps struct {
	// ID must match the parent Select ID for signal scoping
	ID string

	// Placeholder text to show when no value is selected
	Placeholder string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// SelectContentProps defines the properties for the SelectContent component
type SelectContentProps struct {
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

// SelectItemProps defines the properties for the SelectItem component
type SelectItemProps struct {
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

// SelectLabelProps defines the properties for the SelectLabel component
type SelectLabelProps struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// SelectSeparatorProps defines the properties for the SelectSeparator component
type SelectSeparatorProps struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// SelectGroupProps defines the properties for the SelectGroup component
type SelectGroupProps struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}
