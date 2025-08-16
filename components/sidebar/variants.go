package sidebar

import (
	"github.com/coreycole/datastarui/utils"
)

// SidebarDesktopVariants returns the CSS classes for the SidebarDesktop component
func SidebarDesktopVariants(args SidebarDesktopArgs) string {
	// Desktop sidebar - normal flow positioning, hidden on mobile
	baseClasses := "w-64 shrink-0 bg-sidebar text-sidebar-foreground hidden md:block"
	
	return utils.TwMerge(baseClasses, args.Class)
}

// SidebarTriggerVariants returns the CSS classes for the SidebarTrigger component
func SidebarTriggerVariants(args SidebarTriggerArgs) string {
	// Mobile hamburger menu button
	baseClasses := "inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-9 w-9 md:hidden"
	
	return utils.TwMerge(baseClasses, args.Class)
}

// SidebarCloseButtonVariants returns the CSS classes for the SidebarCloseButton component
func SidebarCloseButtonVariants(args SidebarCloseButtonArgs) string {
	// Close button - similar to trigger but no responsive hiding
	baseClasses := "inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-9 w-9"
	
	return utils.TwMerge(baseClasses, args.Class)
}

// SidebarContentVariants returns the CSS classes for the SidebarContent component
func SidebarContentVariants(args SidebarContentArgs) string {
	// Shared content structure
	baseClasses := "flex h-full w-full flex-col bg-sidebar"
	
	return utils.TwMerge(baseClasses, args.Class)
}

// SidebarHeaderVariants returns the CSS classes for the SidebarHeader component
func SidebarHeaderVariants(args SidebarHeaderArgs) string {
	// Header section styling
	baseClasses := "flex h-14 items-center px-4"
	
	return utils.TwMerge(baseClasses, args.Class)
}

// SidebarFooterVariants returns the CSS classes for the SidebarFooter component
func SidebarFooterVariants(args SidebarFooterArgs) string {
	// Footer section styling - auto margin top to push to bottom
	baseClasses := "mt-auto flex items-center px-4"
	
	return utils.TwMerge(baseClasses, args.Class)
}

// SidebarNavLinkVariants returns the CSS classes for the SidebarNavLink component
func SidebarNavLinkVariants(args SidebarNavLinkArgs) string {
	isActive := args.Item.Href == args.CurrentPath
	
	baseClasses := "group relative flex h-8 w-full items-center rounded-lg px-2 text-sm transition-colors hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
	
	if isActive {
		baseClasses = utils.TwMerge(baseClasses, "bg-sidebar-accent text-sidebar-accent-foreground font-medium")
	} else {
		baseClasses = utils.TwMerge(baseClasses, "text-sidebar-foreground")
	}
	
	return baseClasses
}
