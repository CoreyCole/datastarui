package sheet

import (
	"github.com/coreycole/datastarui/utils"
)

// SheetVariants returns the CSS classes for the main Sheet container component
func SheetVariants(args SheetArgs) string {
	// Sheet-specific styling based on side
	baseClasses := ""

	switch args.Side {
	case "top":
		baseClasses = "inset-x-0 top-0 border-b"
	case "bottom":
		baseClasses = "inset-x-0 bottom-0 border-t"
	case "left":
		baseClasses = "inset-y-0 left-0 h-full w-3/4 border-r sm:max-w-sm"
	default: // "right"
		baseClasses = "inset-y-0 right-0 h-full w-3/4 border-l sm:max-w-sm"
	}

	return utils.TwMerge(baseClasses, args.Class)
}

// SheetContentVariants returns the CSS classes for the SheetContent component
func SheetContentVariants(args SheetContentArgs) string {
	// Base classes for sheet content - relative positioning for close button
	baseClasses := "relative flex flex-col"

	return utils.TwMerge(baseClasses, args.Class)
}

// SheetHeaderVariants returns the CSS classes for the SheetHeader component
func SheetHeaderVariants(args SheetHeaderArgs) string {
	// Sheet header styling
	baseClasses := "flex flex-col gap-1.5 p-4"

	return utils.TwMerge(baseClasses, args.Class)
}

// SheetFooterVariants returns the CSS classes for the SheetFooter component
func SheetFooterVariants(args SheetFooterArgs) string {
	// Footer styling - auto margin top to push to bottom
	baseClasses := "mt-auto flex flex-col gap-2 p-4"

	return utils.TwMerge(baseClasses, args.Class)
}

// SheetTitleVariants returns the CSS classes for the SheetTitle component
func SheetTitleVariants(args SheetTitleArgs) string {
	// Sheet title styling
	baseClasses := "text-foreground font-semibold"

	return utils.TwMerge(baseClasses, args.Class)
}

// SheetDescriptionVariants returns the CSS classes for the SheetDescription component
func SheetDescriptionVariants(args SheetDescriptionArgs) string {
	// Sheet description styling
	baseClasses := "text-muted-foreground text-sm"

	return utils.TwMerge(baseClasses, args.Class)
}
