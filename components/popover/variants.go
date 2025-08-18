package popover

import (
	"github.com/coreycole/datastarui/utils"
)

// PopoverContentVariants generates CSS classes for popover content
func PopoverContentVariants(args PopoverContentArgs) string {
	// Base classes matching shadcn/ui popover with native popover API
	base := "w-72 rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-none"

	return utils.TwMerge(base, args.Class)
}

