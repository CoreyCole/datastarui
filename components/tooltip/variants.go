package tooltip

import (
	"github.com/coreycole/datastarui/utils"
)

// TooltipContentVariants generates CSS classes for tooltip content (matching PopoverContentVariants pattern)
func TooltipContentVariants(args TooltipContentArgs) string {
	// Base classes matching popover exactly for identical positioning behavior
	base := "rounded-md border bg-popover text-popover-foreground shadow-md outline-none px-3 py-1.5 text-sm pointer-events-none"
	
	// Animation classes will be handled via data-class for reactive state changes
	// The base classes handle the static styling
	// pointer-events-none prevents tooltip from interfering with trigger hover
	
	return utils.TwMerge(base, args.Class)
}