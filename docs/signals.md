# Datastar Signals Management

This document explains how to use the structured signals management system in DatastarUI components.

**IMPORTANT: This document describes the legacy pattern. For new components, use the utility libraries described in [guide.md](./guide.md).**

## Overview

The signals management system provides a clean, type-safe way to create and manage Datastar signals using Go structs with JSON tags. This eliminates the need to manually construct JSON strings and provides better maintainability.

**Note:** Signal IDs are automatically sanitized to replace hyphens with underscores for JavaScript compatibility. For example, `"theme-select"` becomes `"theme_select"` in the generated signals.

## Modern Approach (Recommended)

Use the utility libraries for cleaner, more maintainable code:

```go
// 1. Define signal structure
type MyComponentSignals struct {
    Open     bool   `json:"open"`
    Value    string `json:"value"`
    Count    int    `json:"count"`
}

// 2. Create signals with utilities
signals := utils.Signals(props.ID, MyComponentSignals{
    Open:  false,
    Value: "",
    Count: 0,
})

// 3. Use in templates with expression builders
clickHandler := utils.NewExpression().
    Statement("evt.preventDefault()").
    SetSignal("modal.open", "false").
    Build()

dataClass := utils.HighlightedItem("select.highlighted", index)
```

## Legacy Patterns (For Reference)

### Basic Signal Manager Usage

```go
// Create signals manager
signals := utils.Signals(props.ID, MyComponentSignals{
    Open:    false,
    Value:   "",
    Count:   0,
    Loading: false,
})

// Use in templates
<div data-signals={ signals.DataSignals }>
  <!-- Component content -->
</div>
```

### Signal References

The SignalManager provides helper methods for creating signal references:

```go
signals.Signal("open")      // Returns: "$[props.ID].open"
signals.Toggle("open")      // Returns: "$[props.ID].open = !$[props.ID].open"
signals.Set("value", "'hello'") // Returns: "$[props.ID].value = 'hello'"
signals.Conditional("loading", "'Saving...'", "'Save'")
// Returns: "$[props.ID].loading ? 'Saving...' : 'Save'"
```

## Migration Guide

### From Manual JSON (Avoid):
```go
// Bad - error-prone string concatenation
signalsJSON := "{\"" + props.ID + "\": {\"open\": false, \"count\": 0}}"
toggleExpr := "$" + props.ID + ".open = !$" + props.ID + ".open"
```

### To Modern Utilities (Recommended):
```go
// Good - structured and maintainable
signals := utils.Signals(props.ID, ComponentSignals{
    Open:  false,
    Count: 0,
})

toggleExpr := utils.NewExpression().
    Statement(signals.Toggle("open")).
    Build()
```

## Best Practices

1. **Use utility libraries** - Avoid string concatenation for expressions
2. **Type-safe signals** - Always use JSON tags on struct fields
3. **Descriptive naming** - Use clear, descriptive signal names
4. **Proper initialization** - Set sensible defaults in struct literals
5. **Namespace by ID** - Use component ID to avoid signal conflicts
6. **Expression builders** - Use `utils.NewExpression()` for complex logic

## See Also

- [Development Guide](./guide.md) - Comprehensive patterns using utilities
- [Debugging Guide](./debugging.md) - Testing and troubleshooting
- [Playwright Testing](./playwright.md) - Browser automation testing