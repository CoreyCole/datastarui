package dialog

import (
	"github.com/coreycole/datastarui/utils"
)

// DialogVariants returns the CSS classes for the main Dialog container component
func DialogVariants(props DialogProps) string {
	// Dialog-specific styling - optimized for modal dialogs with consistent padding
	baseClasses := "max-w-lg w-full max-h-[90vh] overflow-auto bg-background border shadow-lg rounded-lg p-6"

	return utils.TwMerge(baseClasses, props.Class)
}

// DialogContentVariants returns the CSS classes for the DialogContent component
func DialogContentVariants(props DialogContentProps) string {
	// Base classes for dialog content - no horizontal padding since container handles it
	baseClasses := "py-4"

	return utils.TwMerge(baseClasses, props.Class)
}

// DialogOverlayVariants returns the CSS classes for the DialogOverlay component
func DialogOverlayVariants(props DialogOverlayProps) string {
	baseClasses := "fixed inset-0 z-50 bg-black/50 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0"

	return utils.TwMerge(baseClasses, props.Class)
}

// DialogHeaderVariants returns the CSS classes for the DialogHeader component
func DialogHeaderVariants(props DialogHeaderProps) string {
	// Dialog header styling - no padding since container handles it, just spacing between elements
	baseClasses := "flex flex-col gap-2 text-left"

	return utils.TwMerge(baseClasses, props.Class)
}

// DialogFooterVariants returns the CSS classes for the DialogFooter component
func DialogFooterVariants(props DialogFooterProps) string {
	// Footer with top margin to separate from content, no horizontal padding
	baseClasses := "flex flex-row gap-3 justify-end pt-4"

	return utils.TwMerge(baseClasses, props.Class)
}

// DialogTitleVariants returns the CSS classes for the DialogTitle component
func DialogTitleVariants(props DialogTitleProps) string {
	// Dialog title styling - focused on readability and hierarchy
	baseClasses := "text-lg leading-none font-semibold"

	return utils.TwMerge(baseClasses, props.Class)
}

// DialogDescriptionVariants returns the CSS classes for the DialogDescription component
func DialogDescriptionVariants(props DialogDescriptionProps) string {
	// Dialog description styling - similar to card but optimized for dialog context
	baseClasses := "text-sm text-muted-foreground"

	return utils.TwMerge(baseClasses, props.Class)
}

// DialogCloseVariants returns the CSS classes for the DialogClose component
func DialogCloseVariants(props DialogCloseProps) string {
	baseClasses := "absolute right-4 top-4 rounded-sm transition-all hover:bg-accent hover:text-accent-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none z-10 p-2 bg-white/90 border border-gray-200 shadow-sm hover:shadow-md text-gray-900 hover:text-gray-900"

	return utils.TwMerge(baseClasses, props.Class)
}
