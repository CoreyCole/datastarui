package dropdown

import "github.com/coreycole/datastarui/utils"

// dropdownMenuVariants generates CSS classes for the main DropdownMenu root component
func dropdownMenuVariants(className string) string {
	// Add relative positioning for proper dropdown positioning
	baseClasses := "relative"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// dropdownMenuTriggerVariants generates CSS classes for the DropdownMenuTrigger component
func dropdownMenuTriggerVariants(className string) string {
	// No specific base classes for the trigger in shadcn/ui dropdown
	classes := []string{}

	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// dropdownMenuContentVariants generates CSS classes for the DropdownMenuContent component
func dropdownMenuContentVariants(className string) string {
	// Extract EXACT base classes from shadcn/ui dropdown - v4 registry
	baseClasses := "z-50 min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// dropdownMenuItemVariants generates CSS classes for the DropdownMenuItem component
func dropdownMenuItemVariants(variant, className string, inset, disabled bool) string {
	// Extract EXACT base classes from shadcn/ui dropdown
	baseClasses := "focus:bg-accent focus:text-accent-foreground relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"

	// Map EXACT variant classes from shadcn/ui
	variantClasses := map[string]string{
		"default":     "",
		"destructive": "data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 dark:data-[variant=destructive]:focus:bg-destructive/20 data-[variant=destructive]:focus:text-destructive data-[variant=destructive]:*:[svg]:!text-destructive",
	}

	// Apply defaults
	if variant == "" {
		variant = "default"
	}

	classes := []string{baseClasses}
	if variantClasses[variant] != "" {
		classes = append(classes, variantClasses[variant])
	}
	if inset {
		classes = append(classes, "data-[inset]:pl-8")
	}
	if !disabled {
		classes = append(classes, "[&_svg:not([class*='text-'])]:text-muted-foreground")
	}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// dropdownMenuLabelVariants generates CSS classes for the DropdownMenuLabel component
func dropdownMenuLabelVariants(className string, inset bool) string {
	// Extract EXACT base classes from shadcn/ui dropdown
	baseClasses := "px-2 py-1.5 text-sm font-medium"

	classes := []string{baseClasses}
	if inset {
		classes = append(classes, "data-[inset]:pl-8")
	}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// dropdownMenuSeparatorVariants generates CSS classes for the DropdownMenuSeparator component
func dropdownMenuSeparatorVariants(className string) string {
	// Extract EXACT base classes from shadcn/ui dropdown
	baseClasses := "bg-border -mx-1 my-1 h-px"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// dropdownMenuShortcutVariants generates CSS classes for the DropdownMenuShortcut component
func dropdownMenuShortcutVariants(className string) string {
	// Extract EXACT base classes from shadcn/ui dropdown
	baseClasses := "text-muted-foreground ml-auto text-xs tracking-widest"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// dropdownMenuGroupVariants generates CSS classes for the DropdownMenuGroup component
func dropdownMenuGroupVariants(className string) string {
	// No specific base classes for the group element in shadcn/ui dropdown
	classes := []string{}

	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}
