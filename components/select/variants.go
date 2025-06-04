package selectcomponent

import "github.com/coreycole/datastarui/utils"

// selectVariants returns the CSS classes for the main select container
func selectVariants(class string) string {
	baseClasses := "relative"

	classes := []string{baseClasses}
	if class != "" {
		classes = append(classes, class)
	}

	return utils.TwMerge(classes...)
}

// selectTriggerVariants returns the CSS classes for the select trigger
func selectTriggerVariants(class string) string {
	baseClasses := "flex h-9 w-full items-center justify-between whitespace-nowrap rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50 [&>span]:line-clamp-1"

	classes := []string{baseClasses}
	if class != "" {
		classes = append(classes, class)
	}

	return utils.TwMerge(classes...)
}

// selectContentVariants returns the CSS classes for the select content
func selectContentVariants(position string, class string) string {
	baseClasses := "z-50 max-h-96 min-w-[8rem] overflow-hidden rounded-md border bg-popover text-popover-foreground shadow-md data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2"

	classes := []string{baseClasses}
	if position == "popper" {
		classes = append(classes, "data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1")
	}
	if class != "" {
		classes = append(classes, class)
	}

	return utils.TwMerge(classes...)
}

// selectItemVariants returns the CSS classes for the select item
func selectItemVariants(class string) string {
	baseClasses := "relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-2 pr-8 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"

	classes := []string{baseClasses}
	if class != "" {
		classes = append(classes, class)
	}

	return utils.TwMerge(classes...)
}

// selectLabelVariants returns the CSS classes for the select label
func selectLabelVariants(class string) string {
	baseClasses := "px-2 py-1.5 text-sm font-semibold"

	classes := []string{baseClasses}
	if class != "" {
		classes = append(classes, class)
	}

	return utils.TwMerge(classes...)
}

// selectSeparatorVariants returns the CSS classes for the select separator
func selectSeparatorVariants(class string) string {
	baseClasses := "-mx-1 my-1 h-px bg-muted"

	classes := []string{baseClasses}
	if class != "" {
		classes = append(classes, class)
	}

	return utils.TwMerge(classes...)
}

// selectViewportVariants returns the CSS classes for the select viewport
func selectViewportVariants(position string, class string) string {
	baseClasses := "p-1"

	classes := []string{baseClasses}
	if position == "popper" {
		classes = append(classes, "h-[var(--radix-select-trigger-height)] w-full min-w-[var(--radix-select-trigger-width)]")
	}
	if class != "" {
		classes = append(classes, class)
	}

	return utils.TwMerge(classes...)
}
