package button

import (
	"github.com/coreycole/datastarui/utils"
)

// buttonVariants generates the appropriate CSS classes for the button based on variant and size
// This mimics the class-variance-authority (cva) functionality from the original shadcn/ui
func buttonVariants(variant, size, className string) string {
	// Base classes that are always applied - exact copy from New York v4
	baseClasses := "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive"

	// Variant classes - exact copy from New York v4
	variantClasses := map[string]string{
		"default":     "bg-primary text-primary-foreground shadow-xs hover:bg-primary/90",
		"destructive": "bg-destructive text-white shadow-xs hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40",
		"outline":     "border bg-background shadow-xs hover:bg-accent hover:text-accent-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50",
		"secondary":   "bg-secondary text-secondary-foreground shadow-xs hover:bg-secondary/80",
		"ghost":       "hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50",
		"link":        "text-primary underline-offset-4 hover:underline",
	}

	// Size classes - exact copy from New York v4
	sizeClasses := map[string]string{
		"default": "h-9 px-4 py-2 has-[>svg]:px-3",
		"sm":      "h-8 rounded-md gap-1.5 px-3 has-[>svg]:px-2.5",
		"lg":      "h-10 rounded-md px-6 has-[>svg]:px-4",
		"icon":    "size-9",
	}

	// Set defaults if not provided
	if variant == "" {
		variant = "default"
	}
	if size == "" {
		size = "default"
	}

	// Get the appropriate classes
	variantClass := variantClasses[variant]
	sizeClass := sizeClasses[size]

	// Combine all classes
	classes := []string{baseClasses}
	if variantClass != "" {
		classes = append(classes, variantClass)
	}
	if sizeClass != "" {
		classes = append(classes, sizeClass)
	}
	if className != "" {
		classes = append(classes, className)
	}

	// Use the utility function to merge classes (similar to cn() in shadcn/ui)
	return utils.TwMerge(classes...)
}
