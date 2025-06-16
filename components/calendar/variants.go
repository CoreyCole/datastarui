package calendar

import "github.com/coreycole/datastarui/utils"

// calendarVariants generates CSS classes for the calendar container
func calendarVariants(className string) string {
	base := "p-3"

	return utils.TwMerge(base, className)
}
