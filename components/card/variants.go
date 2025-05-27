package card

import (
	"github.com/coreycole/datastarui/templui/utils"
)

// cardVariants generates the appropriate CSS classes for the main Card component
// Based on the New York v4 style from shadcn/ui
func cardVariants(className string) string {
	// Base classes from New York v4 - exact copy
	baseClasses := "bg-card text-card-foreground flex flex-col gap-6 rounded-xl border py-6 shadow-sm"

	// Combine classes
	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// cardHeaderVariants generates the appropriate CSS classes for the CardHeader component
func cardHeaderVariants(className string) string {
	// Base classes from New York v4 - exact copy
	baseClasses := "@container/card-header grid auto-rows-min grid-rows-[auto_auto] items-start gap-1.5 px-6 has-data-[slot=card-action]:grid-cols-[1fr_auto] [.border-b]:pb-6"

	// Combine classes
	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// cardTitleVariants generates the appropriate CSS classes for the CardTitle component
func cardTitleVariants(className string) string {
	// Base classes from New York v4 - exact copy
	baseClasses := "leading-none font-semibold"

	// Combine classes
	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// cardDescriptionVariants generates the appropriate CSS classes for the CardDescription component
func cardDescriptionVariants(className string) string {
	// Base classes from New York v4 - exact copy
	baseClasses := "text-muted-foreground text-sm"

	// Combine classes
	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// cardActionVariants generates the appropriate CSS classes for the CardAction component
func cardActionVariants(className string) string {
	// Base classes from New York v4 - exact copy
	baseClasses := "col-start-2 row-span-2 row-start-1 self-start justify-self-end"

	// Combine classes
	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// cardContentVariants generates the appropriate CSS classes for the CardContent component
func cardContentVariants(className string) string {
	// Base classes from New York v4 - exact copy
	baseClasses := "px-6"

	// Combine classes
	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// cardFooterVariants generates the appropriate CSS classes for the CardFooter component
func cardFooterVariants(className string) string {
	// Base classes from New York v4 - exact copy
	baseClasses := "flex items-center px-6 [.border-t]:pt-6"

	// Combine classes
	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}
