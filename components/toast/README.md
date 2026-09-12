# Toast Component

A shadcn/ui-parity Toast (Sonner-style) component for DatastarUI providing ephemeral UI feedback.

## Features

- **Multiple Variants**: `default`, `success`, `destructive`, `info`
- **Flexible Positioning**: 6 positions (top/bottom × left/center/right)
- **Auto-dismiss**: Configurable duration with `setTimeout`
- **Manual Dismiss**: Close button appears on hover
- **Smooth Animations**: Fade-in/out and slide animations using `data-state`
- **Accessible**: Proper ARIA roles and live regions
- **Multiple Toasts**: Stack multiple toasts in a container

## Usage

### Basic Example

```go
// 1. Add a ToastContainer with ToastItems
@toast.ToastContainer(toast.ToastContainerArgs{
    ID:       "notifications",
    Position: toast.ToastPositionTopRight,
}) {
    @toast.ToastItem(toast.ToastItemArgs{
        ID:          "success_message",
        Title:       "Copied!",
        Description: "Link copied to clipboard.",
        Variant:     toast.ToastVariantSuccess,
    })
}

// 2. Trigger the toast
@toast.ToastTrigger("success_message", 2000) {
    @button.Button(button.ButtonArgs{Variant: "default"}) {
        Copy Link
    }
}
```

### Clipboard Example (Primary Use Case)

```go
// In your root layout (once)
@toast.ToastContainer(toast.ToastContainerArgs{
    ID:       "app_toasts",
    Position: toast.ToastPositionTopRight,
}) {
    @toast.ToastItem(toast.ToastItemArgs{
        ID:      "clipboard_success",
        Title:   "Copied to clipboard",
        Variant: toast.ToastVariantSuccess,
    })
}

// In your share menu or copy button (anywhere)
@toast.ToastTrigger("clipboard_success", 2000) {
    @button.Button(button.ButtonArgs{
        Variant: "ghost",
        Size:    "sm",
    }) {
        Copy Link
    }
}
```

### Variants

```go
// Success (green)
Variant: toast.ToastVariantSuccess

// Error/Destructive (red)
Variant: toast.ToastVariantDestructive

// Info (blue)
Variant: toast.ToastVariantInfo

// Default (neutral)
Variant: toast.ToastVariantDefault
```

### Positions

```go
Position: toast.ToastPositionTopLeft      // Top-left corner
Position: toast.ToastPositionTopCenter    // Top center
Position: toast.ToastPositionTopRight     // Top-right corner (default)
Position: toast.ToastPositionBottomLeft   // Bottom-left corner
Position: toast.ToastPositionBottomCenter // Bottom center
Position: toast.ToastPositionBottomRight  // Bottom-right corner
```

### No Auto-Dismiss

Set duration to `0` for persistent toasts:

```go
@toast.ToastTrigger("persistent_toast", 0) {
    // Must be manually closed
}
```

## Architecture

### Component Pattern

Follows standard DatastarUI structure:
- `args.go` - Component arguments and types
- `variants.go` - CSS class generation
- `expressions.go` - Datastar expression builders
- `toast.templ` - templ templates with signal struct

### Datastar Integration

- Each toast has its own `open` boolean signal
- ToastTrigger sets `$toastID.open = true`
- Auto-dismiss uses `setTimeout(() => { $toastID.open = false }, duration)`
- Animations driven by `data-state` attribute

### Backend/SSR Compatibility

Toasts can be pre-rendered on the server and shown client-side, or triggered via SSE:

```go
// Backend handler (example)
func HandleCopyLink(c echo.Context) error {
    w := c.Response().Writer
    r := c.Request()
    sse := datastar.NewSSE(w, r)
    
    // Show success toast via SSE
    return sse.ExecuteScript("$copy_success.open = true; setTimeout(() => { $copy_success.open = false }, 2000)")
}
```

## Demo

Visit `/components/toast` to see:
- All variants in action
- Multiple positioning options
- Auto-dismiss behavior
- Manual close functionality
- Multiple toasts stacking

## Testing

E2E tests available in:
- `toast_e2e_helpers_test.go`
- `toast_component_e2e_test.go`

Run with: `just e2e --story toast`
