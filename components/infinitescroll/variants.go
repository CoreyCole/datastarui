package infinitescroll

import (
	"github.com/coreycole/datastarui/utils"
)

// InfiniteScrollVariants generates CSS classes for the InfiniteScroll component
func InfiniteScrollVariants(args InfiniteScrollArgs) string {
	baseClasses := "relative"
	if args.Class != "" {
		return utils.TwMerge(baseClasses, args.Class)
	}
	return baseClasses
}

// HostVariants generates CSS classes for the Host container
func HostVariants(args HostArgs) string {
	baseClasses := "relative overflow-auto"
	if args.Class != "" {
		return utils.TwMerge(baseClasses, args.Class)
	}
	return baseClasses
}

// ItemsVariants generates CSS classes for the Items container
func ItemsVariants(args ItemsArgs) string {
	baseClasses := "space-y-0"
	if args.Class != "" {
		return utils.TwMerge(baseClasses, args.Class)
	}
	return baseClasses
}

// SentinelVariants generates CSS classes for sentinel elements
func SentinelVariants(args SentinelArgs) string {
	baseClasses := "h-px w-full"
	if args.Class != "" {
		return utils.TwMerge(baseClasses, args.Class)
	}
	return baseClasses
}

// LoadingVariants generates CSS classes for loading indicators
func LoadingVariants(args LoadingArgs) string {
	baseClasses := "flex items-center justify-center p-4"
	if args.Class != "" {
		return utils.TwMerge(baseClasses, args.Class)
	}
	return baseClasses
}
