package toast

import (
	"github.com/coreycole/datastarui/utils"
)

// ToastContainerVariants generates CSS classes for the toast container
func ToastContainerVariants(args ToastContainerArgs) string {
	base := "fixed z-[100] flex flex-col gap-2 p-4 pointer-events-none"

	// Position classes
	var position string
	switch args.Position {
	case ToastPositionTopLeft:
		position = "top-0 left-0"
	case ToastPositionTopCenter:
		position = "top-0 left-1/2 -translate-x-1/2"
	case ToastPositionTopRight:
		position = "top-0 right-0"
	case ToastPositionBottomLeft:
		position = "bottom-0 left-0"
	case ToastPositionBottomCenter:
		position = "bottom-0 left-1/2 -translate-x-1/2"
	case ToastPositionBottomRight:
		position = "bottom-0 right-0"
	default:
		position = "top-0 right-0" // Default to top-right
	}

	classes := utils.TwMerge(base, position, args.Class)
	return classes
}

// ToastItemVariants generates CSS classes for a toast item
func ToastItemVariants(args ToastItemArgs) string {
	// Base classes with animation support
	base := "pointer-events-auto group relative flex w-full items-center justify-between space-x-4 overflow-hidden rounded-md border p-6 pr-8 shadow-lg transition-all duration-300 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:slide-out-to-right-full data-[state=open]:slide-in-from-right-full"

	// Variant classes
	var variant string
	switch args.Variant {
	case ToastVariantSuccess:
		variant = "border-green-500/50 bg-green-50 text-green-900 dark:border-green-500/50 dark:bg-green-950 dark:text-green-50"
	case ToastVariantDestructive:
		variant = "border-destructive/50 bg-destructive text-destructive-foreground"
	case ToastVariantInfo:
		variant = "border-blue-500/50 bg-blue-50 text-blue-900 dark:border-blue-500/50 dark:bg-blue-950 dark:text-blue-50"
	default:
		variant = "border bg-background text-foreground"
	}

	classes := utils.TwMerge(base, variant, args.Class)
	return classes
}
