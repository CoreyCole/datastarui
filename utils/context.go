package utils

import (
	"context"
)

// contextKey is a type for context keys
type contextKey string

const (
	// IsMobileKey is the context key for mobile detection
	IsMobileKey contextKey = "isMobile"
)

// WithMobile adds the mobile detection result to the context
func WithMobile(ctx context.Context, isMobile bool) context.Context {
	return context.WithValue(ctx, IsMobileKey, isMobile)
}

// GetMobile retrieves the mobile detection result from the context
func GetMobile(ctx context.Context) bool {
	if isMobile, ok := ctx.Value(IsMobileKey).(bool); ok {
		return isMobile
	}
	return false
}