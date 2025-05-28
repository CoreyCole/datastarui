package checkbox

import (
	"github.com/coreycole/datastarui/utils"
)

// checkboxVariants generates the appropriate CSS classes for the Checkbox component
func checkboxVariants(className string) string {
	// Base classes from New York v4 - exact copy from shadcn/ui source
	baseClasses := "peer border-input dark:bg-input/30 data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground dark:data-[state=checked]:bg-primary data-[state=checked]:border-primary focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive size-4 shrink-0 rounded-[4px] border shadow-xs transition-shadow outline-none focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// checkboxIndicatorVariants generates the appropriate CSS classes for the CheckboxIndicator component
func checkboxIndicatorVariants(className string) string {
	// Base classes from New York v4 - exact copy from shadcn/ui source
	baseClasses := "flex items-center justify-center text-current transition-none"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}
