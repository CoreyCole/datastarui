# DatastarUI Development Guide

This guide explains how to build clean, maintainable Datastar components using templ and our utility libraries.

## Table of Contents
- [Core Utilities](#core-utilities)
- [Signal Management](#signal-management)
- [Expression Builders](#expression-builders)
- [Conditional Classes](#conditional-classes)
- [Component Patterns](#component-patterns)
- [Best Practices](#best-practices)

## Core Utilities

### Signal Management (`utils/signals.go`)

The `SignalManager` provides a structured way to handle Datastar signals with proper namespacing:

```go
// Define your component's signal structure
type SelectSignals struct {
    Open        bool   `json:"open"`
    Value       string `json:"value"`
    Label       string `json:"label"`
    Highlighted int    `json:"highlighted"`
}

// Create signals in your component
signals := utils.Signals("my_select", SelectSignals{
    Open:        false,
    Value:       "default",
    Label:       "Select an option",
    Highlighted: -1,
})
```

**Signal Manager Methods:**
- `signals.DataSignals` - Returns the JSON for `data-signals` attribute
- `signals.Signal("prop")` - Returns signal reference: `$my_select.prop`
- `signals.Toggle("open")` - Toggle boolean: `$my_select.open = !$my_select.open`
- `signals.Set("value", "'new'")` - Set value: `$my_select.value = 'new'`
- `signals.Conditional("open", "true", "false")` - Ternary expression

### Expression Builders (`utils/expressions.go`)

Build complex Datastar expressions without string concatenation:

```go
// Simple expression
expr := utils.NewExpression().
    Statement("evt.preventDefault()").
    SetSignal("select.open", "true").
    Build()
// Output: "evt.preventDefault(); $select.open = true"

// Conditional expression
expr := utils.NewExpression().
    Conditional(
        "$select.open",
        "$select.open = false",
        "null"
    ).
    Build()
// Output: "$select.open ? $select.open = false : null"
```

### Conditional Classes (`utils/data_class.go`)

Create dynamic CSS classes using Datastar's `data-class` attribute:

```go
// Highlight active item
dataClass := utils.HighlightedItem("select.highlighted", index)
// Output: {'bg-accent': $select.highlighted === 0, 'text-accent-foreground': $select.highlighted === 0}

// Custom conditional classes
dataClass := utils.NewDataClass().
    Add("border-primary", "$tabs.active === 'tab1'").
    Add("font-semibold", "$tabs.active === 'tab1'").
    Build()
```

## Component Patterns

### Basic Component Structure

```go
// types.go - Define props and signals
type SelectProps struct {
    ID          string
    Options     []SelectOption
    Placeholder string
    Class       string
}

type SelectSignals struct {
    Open  bool   `json:"open"`
    Value string `json:"value"`
}
```

```go
// select.templ - Component template
templ Select(props SelectProps) {
    {{
        // Initialize signals
        signals := utils.Signals(props.ID, SelectSignals{
            Open:  false,
            Value: props.DefaultValue,
        })
        
        // Build expressions
        clickHandler := signals.Toggle("open")
    }}
    <div 
        data-signals={ signals.DataSignals }
        data-on-click={ clickHandler }
    >
        { children... }
    </div>
}
```

### Keyboard Navigation Pattern

Use expression builders for complex keyboard handlers:

```go
// Build keyboard navigation for dropdown
handler := utils.NewSelectContentHandler(props.ID, signals)
keyboardExpr := handler.BuildKeyboardHandler()

// In template
<div data-on-keydown__window={ keyboardExpr }>
```

### Dynamic Classes Pattern

```go
// Conditional highlighting
templ SelectItem(props SelectItemProps) {
    {{
        highlightClass := utils.HighlightedItem(
            fmt.Sprintf("%s.highlighted", props.ID), 
            props.Index
        )
    }}
    <div 
        class="base-classes"
        data-class={ highlightClass }
    >
        { children... }
    </div>
}
```

## Expression Builder Patterns

### Select Component Handlers

```go
// Trigger button handler
triggerHandler := utils.NewSelectTriggerHandler(selectID, signals)
clickExpr := triggerHandler.BuildClickHandler()
keyExpr := triggerHandler.BuildKeyboardHandler()

// Content navigation handler  
contentHandler := utils.NewSelectContentHandler(selectID, signals)
navExpr := contentHandler.BuildKeyboardHandler()

// Item selection handler
itemHandler := utils.NewSelectItemHandler(selectID, itemValue)
selectExpr := itemHandler.BuildClickHandler()
```

### Custom Keyboard Handlers

```go
// Build custom keyboard handler
keyHandler := utils.NewKeyboardHandler("Enter", "Space").
    OnKeys(func(expr *DatastarExpression) *DatastarExpression {
        return expr.
            PreventDefault().
            SetSignal("modal.open", "false").
            Statement("console.log('Modal closed')")
    }).
    Build()
```

## Best Practices

### 1. Signal Naming Conventions
- Use lowercase with underscores for component IDs: `user_profile`, `item_list`
- Never use uppercase, dashes, or periods in signal names
- Validate IDs in your components

### 2. Expression Building
- Use expression builders instead of string concatenation
- Use comma operator for multiple statements in conditionals: `(stmt1, stmt2, stmt3)`
- Always escape string values in expressions

### 3. Component Composition
```go
// Parent provides signals context
templ SelectWrapper(props WrapperProps) {
    @Select(SelectProps{
        ID: "my_select",
        Options: props.Options,
    }) {
        @SelectTrigger(SelectTriggerProps{ID: "my_select"})
        @SelectContent(SelectContentProps{ID: "my_select"})
    }
}
```

### 4. Avoiding Common Pitfalls

**DON'T: String concatenation**
```go
// Bad
clickHandler := "$" + id + ".open = !" + "$" + id + ".open"
```

**DO: Use utilities**
```go
// Good
clickHandler := signals.Toggle("open")
```

**DON'T: Mix quote styles**
```go
// Bad - will cause syntax errors
expr := `$select.value = "` + value + `"`
```

**DO: Consistent escaping**
```go
// Good
expr := signals.Set("value", fmt.Sprintf("'%s'", value))
```

## Advanced Patterns

### Multi-Signal Coordination
```go
// Update multiple related signals
expr := utils.NewExpression().
    SetSignal("datepicker.open", "false").
    SetSignal("datepicker.selectedDate", "evt.target.dataset.date").
    SetSignal("datepicker.view", "'day'").
    Build()
```

### Complex Conditionals
```go
// Nested conditions with proper parentheses
expr := utils.NewExpression().
    Conditional(
        "evt.key === 'Enter' && $form.valid",
        "($form.submitting = true, @post('/api/submit'))",
        "evt.key === 'Escape' ? $form.reset() : null"
    ).
    Build()
```

### Event Modifiers
```go
// Click outside handler
<div data-on-click__outside={ signals.Set("open", "false") }>

// Window-level keyboard handler
<div data-on-keydown__window={ keyboardHandler }>

// Debounced input
<input data-on-input__debounce_500ms={ signals.Set("search", "evt.target.value") }>
```

## Testing Expressions

Always test your expressions in the browser console:
1. Open DevTools
2. Check for Datastar errors in console
3. Use `mcp__playwright__playwright_console_logs` to capture runtime errors
4. Verify signal updates with Datastar DevTools

## Migration Guide

Converting string concatenation to utilities:

**Before:**
```go
clickHandler := fmt.Sprintf("(evt) => { %s = !%s; %s = %s ? 0 : -1 }", 
    openSignal, openSignal, highlightSignal, openSignal)
```

**After:**
```go
clickHandler := utils.NewExpression().
    Statement(signals.Toggle("open")).
    SetSignal("highlighted", signals.Conditional("open", "0", "-1")).
    Build()
```

This approach provides:
- Type safety where possible
- Readable, maintainable code
- Proper escaping and syntax
- Reusable patterns