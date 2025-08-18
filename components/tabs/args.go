package tabs

import "github.com/a-h/templ"

// TabsArgs defines the properties for the Tabs root component
type TabsArgs struct {
	// ID is a unique identifier for this tabs instance (required for signal management)
	ID string

	// DefaultValue sets the initial active tab
	DefaultValue string

	// Value controls the active tab (for controlled usage)
	Value string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// TabsListArgs defines the properties for the TabsList component
type TabsListArgs struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// TabsTriggerArgs defines the properties for the TabsTrigger component
type TabsTriggerArgs struct {
	// ID is the tabs instance ID (must match the parent Tabs ID)
	ID string

	// Value identifies this trigger and associates it with content
	Value string

	// Disabled makes the trigger non-interactive
	Disabled bool

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// TabsContentArgs defines the properties for the TabsContent component
type TabsContentArgs struct {
	// ID is the tabs instance ID (must match the parent Tabs ID)
	ID string

	// Value identifies this content and associates it with a trigger
	Value string

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}
