package textarea

import "github.com/coreycole/datastarui/utils"

func textareaVariants(className string, borderless bool) string {
	// Extract EXACT base classes from shadcn/ui New York v4 textarea
	baseClasses := "placeholder:text-muted-foreground flex field-sizing-content min-h-16 w-full rounded-md bg-transparent px-3 py-2 text-base text-foreground transition-[color,box-shadow] outline-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm"

	// Border and focus styles (only when not borderless)
	borderClasses := "border-input border shadow-xs dark:bg-input/30 focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive"

	classes := []string{baseClasses}
	if !borderless {
		classes = append(classes, borderClasses)
	}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}
