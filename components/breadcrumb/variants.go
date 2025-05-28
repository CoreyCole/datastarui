package breadcrumb

import "github.com/coreycole/datastarui/utils"

// breadcrumbVariants generates CSS classes for the main Breadcrumb nav component
func breadcrumbVariants(className string) string {
	// No specific base classes for the nav element in shadcn/ui breadcrumb
	classes := []string{}

	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// breadcrumbListVariants generates CSS classes for the BreadcrumbList ol component
func breadcrumbListVariants(className string) string {
	// Extract EXACT base classes from shadcn/ui breadcrumb
	baseClasses := "text-muted-foreground flex flex-wrap items-center gap-1.5 text-sm break-words sm:gap-2.5"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// breadcrumbItemVariants generates CSS classes for the BreadcrumbItem li component
func breadcrumbItemVariants(className string) string {
	// Extract EXACT base classes from shadcn/ui breadcrumb
	baseClasses := "inline-flex items-center gap-1.5"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// breadcrumbLinkVariants generates CSS classes for the BreadcrumbLink a component
func breadcrumbLinkVariants(className string) string {
	// Extract EXACT base classes from shadcn/ui breadcrumb
	baseClasses := "hover:text-foreground transition-colors"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// breadcrumbPageVariants generates CSS classes for the BreadcrumbPage span component
func breadcrumbPageVariants(className string) string {
	// Extract EXACT base classes from shadcn/ui breadcrumb
	baseClasses := "text-foreground font-normal"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// breadcrumbSeparatorVariants generates CSS classes for the BreadcrumbSeparator li component
func breadcrumbSeparatorVariants(className string) string {
	// Extract EXACT base classes from shadcn/ui breadcrumb
	baseClasses := "[&>svg]:size-3.5"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}

// breadcrumbEllipsisVariants generates CSS classes for the BreadcrumbEllipsis span component
func breadcrumbEllipsisVariants(className string) string {
	// Extract EXACT base classes from shadcn/ui breadcrumb
	baseClasses := "flex size-9 items-center justify-center"

	classes := []string{baseClasses}
	if className != "" {
		classes = append(classes, className)
	}

	return utils.TwMerge(classes...)
}
