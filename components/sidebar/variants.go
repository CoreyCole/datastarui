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
	baseClasses := "flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-9 w-9 md:hidden cursor-pointer"
	
	return utils.TwMerge(baseClasses, args.Class)
}

// SidebarCloseButtonVariants returns the CSS classes for the SidebarCloseButton component
func SidebarCloseButtonVariants(args SidebarCloseButtonArgs) string {
	// Close button - similar to trigger but no responsive hiding
	baseClasses := "inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-9 w-9 cursor-pointer"
	
	return utils.TwMerge(baseClasses, args.Class)
}

// SidebarContentVariants returns the CSS classes for the SidebarContent component
func SidebarContentVariants(args SidebarContentArgs) string {
	// Shared content structure with solid background (matches main content)
	baseClasses := "flex h-full w-full flex-col bg-background"
	
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

// SidebarNavLinkMobileVariants returns the CSS classes for mobile SidebarNavLink
func SidebarNavLinkMobileVariants(args SidebarNavLinkArgs) string {
	isActive := args.Item.Href == args.CurrentPath
	
	// Mobile: larger text for touch interfaces
	baseClasses := "group relative flex w-full items-center rounded-md p-2 text-2xl font-medium transition-colors cursor-pointer"
	
	if isActive {
		// Active state: use accent background and foreground (not sidebar-accent)
		baseClasses = utils.TwMerge(baseClasses, "bg-accent text-accent-foreground font-medium border border-accent")
	} else {
		// Inactive state: default text color with hover effect
		baseClasses = utils.TwMerge(baseClasses, "text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:bg-sidebar-accent focus-visible:text-sidebar-accent-foreground")
	}
	
	return baseClasses
}

// SidebarNavLinkDesktopVariants returns the CSS classes for desktop SidebarNavLink
func SidebarNavLinkDesktopVariants(args SidebarNavLinkArgs) string {
	isActive := args.Item.Href == args.CurrentPath
	
	// Desktop: smaller, more compact text with shadcn/ui styling
	baseClasses := "group relative flex w-full items-center rounded-md p-2 text-[0.8rem] font-medium outline-hidden transition-[width,height,padding] focus-visible:ring-2 focus-visible:ring-sidebar-ring cursor-pointer"
	
	if isActive {
		// Active state: use accent background and foreground (not sidebar-accent)
		baseClasses = utils.TwMerge(baseClasses, "bg-accent text-accent-foreground font-medium border border-accent")
	} else {
		// Inactive state: default text color with hover effect
		baseClasses = utils.TwMerge(baseClasses, "text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:bg-sidebar-accent focus-visible:text-sidebar-accent-foreground")
	}
	
	return baseClasses
}
