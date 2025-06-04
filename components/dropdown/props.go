package dropdown

import "github.com/a-h/templ"

// DropdownMenuProps defines the properties for the DropdownMenu root component
type DropdownMenuProps struct {
	// ID provides a unique identifier for the dropdown (auto-generated if not provided)
	ID string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes

	// Open controls the open state (for controlled usage)
	Open bool

	// DefaultOpen sets the initial open state
	DefaultOpen bool
}

// DropdownMenuTriggerProps defines the properties for the DropdownMenuTrigger component
type DropdownMenuTriggerProps struct {
	// ID is the dropdown identifier (required to link to the correct dropdown)
	ID string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes

	// AsChild renders the trigger as a child element (for composition)
	AsChild bool

	// Disabled makes the trigger non-interactive
	Disabled bool
}

// DropdownMenuContentProps defines the properties for the DropdownMenuContent component
type DropdownMenuContentProps struct {
	// ID is the dropdown identifier (required to link to the correct dropdown)
	ID string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes

	// Align specifies the alignment relative to the trigger
	// Options: "start", "center", "end"
	Align string

	// Side specifies which side of the trigger to align to
	// Options: "top", "right", "bottom", "left"
	Side string

	// SideOffset specifies the distance from the trigger
	SideOffset int
}

// DropdownMenuItemProps defines the properties for the DropdownMenuItem component
type DropdownMenuItemProps struct {
	// ID is the dropdown identifier (required to link to the correct dropdown)
	ID string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes

	// Inset adds left padding for nested items
	Inset bool

	// Variant defines the visual style of the item
	// Options: "default", "destructive"
	Variant string

	// Disabled makes the item non-interactive
	Disabled bool

	// AsChild renders the item as a child element (for composition)
	AsChild bool
}

// DropdownMenuLabelProps defines the properties for the DropdownMenuLabel component
type DropdownMenuLabelProps struct {
	// ID is the dropdown identifier (optional for labels)
	ID string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes

	// Inset adds left padding for nested labels
	Inset bool
}

// DropdownMenuSeparatorProps defines the properties for the DropdownMenuSeparator component
type DropdownMenuSeparatorProps struct {
	// ID is the dropdown identifier (optional for separators)
	ID string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// DropdownMenuShortcutProps defines the properties for the DropdownMenuShortcut component
type DropdownMenuShortcutProps struct {
	// ID is the dropdown identifier (optional for shortcuts)
	ID string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// DropdownMenuGroupProps defines the properties for the DropdownMenuGroup component
type DropdownMenuGroupProps struct {
	// ID is the dropdown identifier (optional for groups)
	ID string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}
