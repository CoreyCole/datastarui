package label

import "github.com/coreycole/datastarui/utils"

func labelVariants(className string) string {
	// Extract EXACT base classes from shadcn/ui New York v4 label
	baseClasses := "flex items-center gap-2 text-sm leading-none font-medium select-none group-data-[disabled=true]:pointer-events-none group-data-[disabled=true]:opacity-50 peer-disabled:cursor-not-allowed peer-disabled:opacity-50"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}
