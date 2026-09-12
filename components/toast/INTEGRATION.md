# Toast Integration Guide for Vamos

## SSR-Safe Root Placement

**IMPORTANT**: Mount `ToastContainer` once in your root layout, not per-page. This ensures:
- No chrome remounting when workbench content changes
- Toasts can be triggered from any page/component
- SSR-safe with consistent DOM structure

### Recommended Pattern

```go
// In layouts/root.templ or your app's root layout
templ Root(args RootArgs) {
	<!DOCTYPE html>
	<html>
		<head>...</head>
		<body>
			<!-- Toast container at root level -->
			@toast.ToastContainer(toast.ToastContainerArgs{
				ID:       "app_toasts",
				Position: toast.ToastPositionTopRight,
			}) {
				<!-- Pre-define toast items -->
				@toast.ToastItem(toast.ToastItemArgs{
					ID:          "clipboard_success",
					Title:       "Copied to clipboard",
					Description: "",
					Variant:     toast.ToastVariantSuccess,
				})
				@toast.ToastItem(toast.ToastItemArgs{
					ID:          "error_message",
					Title:       "Error",
					Description: "",
					Variant:     toast.ToastVariantDestructive,
				})
			}
			
			<!-- Your app content -->
			<main>
				{ children... }
			</main>
		</body>
	</html>
}
```

## Public API

### 1. ToastContainer (Mount Once)

```go
@toast.ToastContainer(toast.ToastContainerArgs{
	ID:       string        // Required: Unique container ID
	Position: ToastPosition // Optional: Defaults to TopRight
	Class:    string        // Optional: Additional CSS classes
})
```

**Positions:**
- `toast.ToastPositionTopRight` (default)
- `toast.ToastPositionTopCenter`
- `toast.ToastPositionTopLeft`
- `toast.ToastPositionBottomRight`
- `toast.ToastPositionBottomCenter`
- `toast.ToastPositionBottomLeft`

### 2. ToastItem (Pre-define in Container)

```go
@toast.ToastItem(toast.ToastItemArgs{
	ID:          string       // Required: Unique toast ID
	Title:       string       // Required: Toast title
	Description string       // Optional: Additional text
	Variant:     ToastVariant // Optional: Defaults to Default
	Class:       string       // Optional: Additional CSS
	Attributes:  templ.Attributes // Optional
})
```

**Variants:**
- `toast.ToastVariantDefault`
- `toast.ToastVariantSuccess`
- `toast.ToastVariantDestructive`
- `toast.ToastVariantInfo`

### 3. ToastTrigger (Use Anywhere)

```go
@toast.ToastTrigger(
	toastID  string, // Required: ID of ToastItem to show
	duration int,    // Required: Auto-dismiss ms (0 = no auto-dismiss)
)
```

### 4. ShowToastExpr (Direct Expression - for ClientActions)

```go
toast.ShowToastExpr(toastID string, durationMs int) string
```

Returns a Datastar expression string for direct use in `data-on:click` or ClientActions. Use this when you need to compose toast display with other operations (e.g., clipboard write).

**Example - Overflow menu with clipboard:**
```go
<button
	data-on:click={ 
		"navigator.clipboard.writeText('https://example.com').then(() => {" +
		toast.ShowToastExpr("clipboard_success", 2000) +
		"})" 
	}
>
	Copy Link
</button>
```

## Usage Example: Clipboard

```go
// 1. In root layout (once)
@toast.ToastContainer(toast.ToastContainerArgs{
	ID: "app_toasts",
}) {
	@toast.ToastItem(toast.ToastItemArgs{
		ID:      "clipboard_success",
		Title:   "Copied to clipboard",
		Variant: toast.ToastVariantSuccess,
	})
}

// 2. In any component/page (e.g., workbench share menu)
@toast.ToastTrigger("clipboard_success", 2000) {
	@button.Button(button.ButtonArgs{
		Variant: "ghost",
		Size:    "sm",
	}) {
		<svg>...</svg> Copy Link
	}
}
```

## Advanced Usage

### Compose with Clipboard Operations

Use `ShowToastExpr()` to show toast after successful clipboard write:

```go
import "github.com/coreycole/datastarui/components/toast"

// In overflow menu or share button
<button
	data-on:click={ 
		"navigator.clipboard.writeText('" + shareURL + "')" +
		".then(() => {" + toast.ShowToastExpr("clipboard_success", 2000) + "})" +
		".catch(() => {" + toast.ShowToastExpr("clipboard_error", 3000) + "})"
	}
>
	Copy Link
</button>
```

### Backend SSE Trigger

```go
// Backend SSE handler
sse.ExecuteScript("$clipboard_success.open = true; setTimeout(() => { $clipboard_success.open = false }, 2000)")
```

## No Local Toast API Required

Vamos will CLI-copy this component into `pkg/datastarui` after merge. No vamos-specific toast API needed — use the public templ API directly.

## Ephemeral UI Signals

Each `ToastItem` creates one signal: `$toastID.open` (boolean). No other state is managed — keep it simple.
