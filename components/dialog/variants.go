package dialog

import (
	"github.com/coreycole/datastarui/utils"
)

// DialogContentVariants returns the CSS classes for the DialogContent component
func DialogContentVariants(props DialogContentProps) string {
	// Base classes for the dialog content - simplified for better centering
	// The parent dialog element handles positioning, this focuses on content styling
	baseClasses := "grid gap-4 border bg-background p-6 shadow-lg duration-200 sm:rounded-lg relative max-h-[80vh] overflow-y-auto"

	// Size variants
	sizeVariant := ""
	switch props.Size {
	case "sm":
		sizeVariant = "w-full max-w-sm"
	case "md":
		sizeVariant = "w-full max-w-md"
	case "lg":
		sizeVariant = "w-full max-w-lg"
	case "xl":
		sizeVariant = "w-full max-w-xl"
	case "2xl":
		sizeVariant = "w-full max-w-2xl"
	case "full":
		sizeVariant = "w-full max-w-[calc(100%-2rem)]"
	default:
		sizeVariant = "w-full max-w-lg" // Default size
	}

	// Combine base classes with size variant
	classes := baseClasses + " " + sizeVariant

	return utils.TwMerge(classes, props.Class)
}

// DialogOverlayVariants returns the CSS classes for the DialogOverlay component
func DialogOverlayVariants(props DialogOverlayProps) string {
	baseClasses := "fixed inset-0 z-50 bg-black/50 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0"

	return utils.TwMerge(baseClasses, props.Class)
}

// DialogHeaderVariants returns the CSS classes for the DialogHeader component
func DialogHeaderVariants(props DialogHeaderProps) string {
	baseClasses := "flex flex-col gap-2 text-center sm:text-left"

	return utils.TwMerge(baseClasses, props.Class)
}

// DialogFooterVariants returns the CSS classes for the DialogFooter component
func DialogFooterVariants(props DialogFooterProps) string {
	baseClasses := "flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"

	return utils.TwMerge(baseClasses, props.Class)
}

// DialogTitleVariants returns the CSS classes for the DialogTitle component
func DialogTitleVariants(props DialogTitleProps) string {
	baseClasses := "text-lg leading-none font-semibold"

	return utils.TwMerge(baseClasses, props.Class)
}

// DialogDescriptionVariants returns the CSS classes for the DialogDescription component
func DialogDescriptionVariants(props DialogDescriptionProps) string {
	baseClasses := "text-sm text-muted-foreground"

	return utils.TwMerge(baseClasses, props.Class)
}

// DialogCloseVariants returns the CSS classes for the DialogClose component
func DialogCloseVariants(props DialogCloseProps) string {
	baseClasses := "absolute right-4 top-4 rounded-sm transition-all hover:bg-accent hover:text-accent-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none z-10 p-2 bg-white/90 border border-gray-200 shadow-sm hover:shadow-md text-gray-900 hover:text-gray-900"

	return utils.TwMerge(baseClasses, props.Class)
}
