# Datastar Signals Management

This document explains how to use the structured signals management system in DatastarUI components.

## Overview

The signals management system provides a clean, type-safe way to create and manage Datastar signals using Go structs with JSON tags. This eliminates the need to manually construct JSON strings and provides better maintainability.

**Note:** Signal IDs are automatically sanitized to replace hyphens with underscores for JavaScript compatibility. For example, `"theme-select"` becomes `"theme_select"` in the generated signals.

## Basic Usage

### 1. Define Signal Structure

First, define a struct with JSON tags for your signal properties:

```go
type MyComponentSignals struct {
    Open     bool   `json:"open"`
    Value    string `json:"value"`
    Count    int    `json:"count"`
    Loading  bool   `json:"loading"`
}
```

### 2. Create Signals Manager

Use the `utils.Signals()` function to create a signals manager:

```go
signals := utils.Signals(props.ID, MyComponentSignals{
    Open:    false,
    Value:   "",
    Count:   0,
    Loading: false,
})
```

### 3. Use in Templates

Use the `DataSignals` property in your template:

```go
templ () {
    <div data-signals={ signals.DataSignals }>
      <!-- Your component content -->
    </div>
}
```

## Signal References

The SignalManager provides several helper methods for creating signal references:

### Signal

Get a reference to a signal property:

```go
signals.Signal("open") // Returns: "$[props.ID].open"
```

### Toggle

Create a toggle expression for boolean properties:

```go
signals.Toggle("open") // Returns: "$[props.ID].open = !$[props.ID].open"
```

### Set

Create a set expression for any property:

```go
signals.Set("value", "'hello'") // Returns: "$[props.ID].value = 'hello'"
signals.Set("count", "0")       // Returns: "$[props.ID].count = 0"
```

### Conditional

Create conditional expressions:

```go
signals.Conditional("loading", "'Saving...'", "'Save'")
// Returns: "$[props.ID].loading ? 'Saving...' : 'Save'"
```

## Complete Example

Here's a complete example of a modal component using the new signals system:

```go
package modal

import "github.com/coreycole/datastarui/utils"

// ModalSignals defines the signal structure
type ModalSignals struct {
    Open bool `json:"open"`
}

templ Modal(props ModalProps) {
    {{
        signals := utils.Signals(props.ID, ModalSignals{
            Open: props.DefaultOpen,
        })
    }}
    <div data-signals={ signals.DataSignals }>
        <!-- Backdrop -->
        <div
            data-show={ signals.Signal("open") }
            class="fixed inset-0 bg-black/50"
            data-on-click={ signals.Set("open", "false") }
        >
            <!-- Modal content -->
            <div class="modal-content">
                { children... }
            </div>
        </div>
    </div>
}

templ ModalTrigger(props ModalTriggerProps) {
    {{
        signals := utils.Signals(props.TargetID, ModalSignals{})
    }}
    <button
        data-on-click={ signals.Set("open", "true") }
        { props.Attributes... }
    >
        { children... }
    </button>
}
```

## Multi-Namespace Signals

For complex components that need multiple signal namespaces, use `MultiSignals`:

```go
signals := utils.MultiSignals(map[string]interface{}{
    "form": FormSignals{
        Name:  "",
        Email: "",
    },
    "ui": UISignals{
        ShowHelp: false,
        Theme:    "light",
    },
})
```

Use `MultiSignalRef` for references:

```go
utils.MultiSignalRef("form", "name") // Returns: "$form.name"
```

## Benefits

1. **Type Safety**: Signal structures are defined with Go structs
2. **Maintainability**: No more manual JSON string construction
3. **Autocomplete**: IDE support for signal properties
4. **Consistency**: Standardized approach across all components
5. **Namespace Support**: Automatic ID-based namespacing prevents conflicts
6. **Helper Methods**: Built-in methods for common signal operations
7. **Datastar Alignment**: Method names match Datastar attribute conventions

## Migration from Manual JSON

### Before (Manual JSON):

```go
signalsJSON := "{\"" + props.ID + "\": {\"open\": false, \"count\": 0}}"
toggleExpr := "$" + props.ID + ".open = !$" + props.ID + ".open"
```

### After (Structured Signals):

```go
type ComponentSignals struct {
    Open  bool `json:"open"`
    Count int  `json:"count"`
}

signals := utils.Signals(props.ID, ComponentSignals{
    Open:  false,
    Count: 0,
})
toggleExpr := signals.Toggle("open")
```

## Best Practices

1. **Always use JSON tags** on your signal struct fields
2. **Use descriptive signal names** that match your component's purpose
3. **Initialize with sensible defaults** in your struct literals
4. **Namespace by component ID** to avoid signal conflicts
5. **Use helper methods** instead of constructing expressions manually
6. **Keep signal structures focused** - one struct per component type

## JavaScript Expression Best Practices

When creating event handlers with Datastar, use the `ConditionalAction` helper for safe conditional logic:

### ✅ Recommended:

```go
// Use ConditionalAction for safe conditional expressions
handler := signals.ConditionalAction("evt.target === evt.currentTarget", "open", "false")
// Results in: evt.target === evt.currentTarget ? ($component.open = false) : void 0
```

### Multiple Actions:

```go
// Use semicolons to separate multiple statements
handler := signals.Set("loading", "true") + "; " + signals.Set("message", "'Processing...'")

// Or use ConditionalMultiAction for conditional multiple actions
handler := signals.ConditionalMultiAction("condition",
    signals.Set("action1", "value1"),
    signals.Set("action2", "value2"))
```
