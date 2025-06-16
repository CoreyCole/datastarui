package datepicker

import "github.com/coreycole/datastarui/utils"

// datePickerVariants generates CSS classes for the datepicker container
func datePickerVariants(className string) string {
	base := "relative w-full max-w-sm"

	return utils.TwMerge(base, className)
}
