package toast

import "github.com/a-h/templ"

// ToastVariant represents the visual style of the toast
type ToastVariant string

const (
	ToastVariantDefault     ToastVariant = "default"
	ToastVariantSuccess     ToastVariant = "success"
	ToastVariantDestructive ToastVariant = "destructive"
	ToastVariantInfo        ToastVariant = "info"
)

// ToastPosition represents where toasts appear on screen
type ToastPosition string

const (
	ToastPositionTopLeft     ToastPosition = "top-left"
	ToastPositionTopCenter   ToastPosition = "top-center"
	ToastPositionTopRight    ToastPosition = "top-right"
	ToastPositionBottomLeft  ToastPosition = "bottom-left"
	ToastPositionBottomCenter ToastPosition = "bottom-center"
	ToastPositionBottomRight ToastPosition = "bottom-right"
)

// ToastContainerArgs defines the properties for the toast container
type ToastContainerArgs struct {
	ID       string        // Required: Unique ID for the toast container
	Position ToastPosition // Position on screen (default: top-right)
	Class    string
}

// ToastItemArgs defines the properties for a single toast notification
type ToastItemArgs struct {
	ID          string       // Required: Unique ID for this toast
	Title       string       // Toast title
	Description string       // Toast description (optional)
	Variant     ToastVariant // Toast variant (default, success, destructive, info)
	Duration    int          // Auto-dismiss duration in milliseconds (0 = no auto-dismiss)
	Class       string
	Attributes  templ.Attributes
}
