package sidebar

import "github.com/a-h/templ"

// SidebarMobileArgs defines the args for the SidebarMobile component
type SidebarMobileArgs struct {
	ID          string            // Shared signal ID
	DefaultOpen bool              // Whether the mobile sidebar should be open by default
	Class       string            // Additional CSS classes
	Attributes  templ.Attributes  // Additional HTML attributes
}

// SidebarDesktopArgs defines the args for the SidebarDesktop component
type SidebarDesktopArgs struct {
	ID         string            // Same shared signal ID (for consistency)
	Class      string            // Additional CSS classes
	Attributes templ.Attributes  // Additional HTML attributes
}

// SidebarTriggerArgs defines the args for the SidebarTrigger component
type SidebarTriggerArgs struct {
	ID         string            // CRITICAL: Same shared signal ID
	Class      string            // Additional CSS classes
	Attributes templ.Attributes  // Additional HTML attributes
}

// SidebarCloseButtonArgs defines the args for the SidebarCloseButton component
type SidebarCloseButtonArgs struct {
	ID         string            // CRITICAL: Same shared signal ID
	Class      string            // Additional CSS classes
	Attributes templ.Attributes  // Additional HTML attributes
}

// SidebarContentArgs defines the args for the SidebarContent component
type SidebarContentArgs struct {
	ID         string            // Same shared signal ID
	Class      string            // Additional CSS classes
	Attributes templ.Attributes  // Additional HTML attributes
}

// SidebarHeaderArgs defines the args for the SidebarHeader component
type SidebarHeaderArgs struct {
	Class      string            // Additional CSS classes
	Attributes templ.Attributes  // Additional HTML attributes
}

// SidebarFooterArgs defines the args for the SidebarFooter component
type SidebarFooterArgs struct {
	Class      string            // Additional CSS classes
	Attributes templ.Attributes  // Additional HTML attributes
}

// SidebarNavArgs defines the args for the SidebarNav component
type SidebarNavArgs struct {
	SidebarID   string          // Reference to the shared sidebar signal ID
	Sections    []SidebarSection // Navigation sections
	CurrentPath string          // Current page path for active state
}

// SidebarNavLinkArgs defines the args for the SidebarNavLink component
type SidebarNavLinkArgs struct {
	SidebarID   string      // Reference to the shared sidebar signal ID
	Item        SidebarItem // Navigation item
	CurrentPath string      // Current page path for active state
}

// SidebarSection represents a section in the sidebar navigation
type SidebarSection struct {
	Title string        // Section title
	Label string        // Optional label (e.g., "New", "Beta")
	Items []SidebarItem // Navigation items in this section
}

// SidebarItem represents a navigation item
type SidebarItem struct {
	Title string // Display title
	Href  string // Navigation URL
	Label string // Optional label (e.g., "New", "Beta")
}