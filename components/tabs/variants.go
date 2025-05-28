package tabs

import (
	"github.com/coreycole/datastarui/utils"
)

// tabsVariants generates the appropriate CSS classes for the Tabs root component
func tabsVariants(className string) string {
	// Base classes from New York v4 - exact copy
	baseClasses := "flex flex-col gap-2"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// tabsListVariants generates the appropriate CSS classes for the TabsList component
func tabsListVariants(className string) string {
	// Base classes from New York v4 - exact copy from shadcn/ui source
	baseClasses := "bg-muted text-muted-foreground inline-flex h-9 w-fit items-center justify-center rounded-lg p-[3px]"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// tabsTriggerVariantsBase generates CSS classes for TabsTrigger without data-[state=active] classes
func tabsTriggerVariantsBase(className string) string {
	// Base classes from New York v4 - exact copy from shadcn/ui source
	baseClasses := "text-foreground dark:text-muted-foreground inline-flex h-[calc(100%-1px)] flex-1 items-center justify-center gap-1.5 rounded-md border border-transparent px-2 py-1 text-sm font-medium whitespace-nowrap transition-[color,box-shadow] focus-visible:ring-[3px] focus-visible:outline-1 focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:outline-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// tabsContentVariants generates the appropriate CSS classes for the TabsContent component
func tabsContentVariants(className string) string {
	// Base classes from New York v4 - exact copy
	baseClasses := "flex-1 outline-none"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}
