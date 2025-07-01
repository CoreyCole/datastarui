# Templ Syntax Guide

A comprehensive guide to templ syntax with real-world examples from the DatastarUI component library.

## Table of Contents

1. [Package Structure & Imports](#package-structure--imports)
2. [Basic Component Syntax](#basic-component-syntax)
3. [Go Code Blocks](#go-code-blocks)
4. [Props and Type Safety](#props-and-type-safety)
5. [Template Expressions](#template-expressions)
6. [Conditional Rendering](#conditional-rendering)
7. [For Loops](#for-loops)
8. [Template Composition](#template-composition)
9. [Attributes and Spreading](#attributes-and-spreading)
10. [Comments](#comments)
11. [State Management Integration](#state-management-integration)
12. [Advanced Patterns](#advanced-patterns)

---

## Package Structure & Imports

Templ files start with a package declaration and imports, just like regular Go files:

```go
package checkbox

import "github.com/coreycole/datastarui/utils"

// You can define types and structs alongside your templates
type CheckboxSignals struct {
	Checked bool `json:"checked"`
}
```

**Key Points:**

- Package name must match the directory name
- Import statements work exactly like Go
- You can mix Go code and templ components in the same file

---

## Basic Component Syntax

Components are defined with the `templ` keyword followed by a function signature:

```go
// Simple component with props
templ Button(props ButtonProps) {
	{{
		// Go code block for logic
		buttonType := props.Type
		if buttonType == "" {
			buttonType = "button"
		}
		classes := buttonVariants(props.Variant, props.Size, props.Class)
	}}
	<button
		type={ buttonType }
		class={ classes }
		disabled?={ props.Disabled }
		{ props.Attributes... }
	>
		{ children... }
	</button>
}
```

**Key Points:**

- Components compile to functions returning `templ.Component`
- Props are passed as function parameters
- `{ children... }` renders child content
- `{{ }}` blocks contain Go code
- `{ }` expressions output dynamic content

---

## Go Code Blocks

Go code blocks (`{{ }}`) allow complex logic before rendering:

```go
templ Tabs(props TabsProps) {
	{{
		// Generate unique ID if not provided
		tabsID := props.ID
		if tabsID == "" {
			b := make([]byte, 4)
			rand.Read(b)
			tabsID = fmt.Sprintf("tabs_%x", b)
		}

		// Set up default values
		defaultValue := props.DefaultValue
		if defaultValue == "" {
			defaultValue = "tab1"
		}

		// Create state management
		signals := utils.Signals(tabsID, TabsSignals{
			Active: defaultValue,
		})
	}}
	<div
		data-signals={ signals.DataSignals }
		class={ tabsVariants(props.Class) }
	>
		{ children... }
	</div>
}
```

**Best Practices:**

- Use Go blocks for complex logic, calculations, and setup
- Keep blocks focused and readable
- Variables declared in blocks are available in the template
- Use proper Go error handling when needed

---

## Props and Type Safety

Define strongly-typed props for your components:

```go
// ButtonProps defines the properties for the Button component
type ButtonProps struct {
	// Variant defines the visual style of the button
	// Options: "default", "destructive", "outline", "secondary", "ghost", "link"
	Variant string

	// Size defines the size of the button
	// Options: "default", "sm", "lg", "icon"
	Size string

	// AsChild renders the button as a child element (for composition)
	AsChild bool

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes

	// Disabled makes the button non-interactive
	Disabled bool

	// Type specifies the button type (button, submit, reset)
	Type string
}

// Usage in component
templ Button(props ButtonProps) {
	{{
		classes := buttonVariants(props.Variant, props.Size, props.Class)
	}}
	<button
		type={ props.Type }
		class={ classes }
		disabled?={ props.Disabled }
		{ props.Attributes... }
	>
		{ children... }
	</button>
}
```

**Key Features:**

- `templ.Attributes` for spreading HTML attributes
- Document expected values in comments
- Use Go's type system for validation
- Optional fields work naturally with Go's zero values

---

## Template Expressions

Dynamic content and attributes use curly braces:

```go
templ Checkbox(props CheckboxProps) {
	{{
		signalRef := signals.Signal("checked")
		toggleExpr := signals.Toggle("checked")
		allClasses := checkboxVariants(props.ClassName)
	}}
	<button
		id={ props.ID }
		class={ allClasses }
		data-on-click={ toggleExpr }
		data-attr-aria-checked={ signalRef + " ? 'true' : 'false'" }
		data-attr-data-state={ signalRef + " ? 'checked' : 'unchecked'" }
	>
		<svg
			class="h-4 w-4"
			data-attr-style={ signalRef + " ? 'opacity: 1' : 'opacity: 0'" }
		>
			<path d="M20 6 9 17l-5-5"></path>
		</svg>
	</button>
}
```

**Expression Types:**

- `{ variable }` - Simple variable output
- `{ expression }` - Go expressions (concatenation, function calls)
- `{ condition ? "true" : "false" }` - Ternary-like expressions
- `attribute={ value }` - Dynamic attributes
- `attribute?={ bool }` - Conditional attributes (present/absent)

---

## Conditional Rendering

Use Go's `if/else` statements for conditional rendering:

```go
templ Button(props ButtonProps) {
	if props.AsChild {
		<!-- When AsChild is true, render as span -->
		<span
			class={ classes }
			if props.Disabled {
				aria-disabled="true"
			}
			{ props.Attributes... }
		>
			{ children... }
		</span>
	} else {
		<!-- Default button element -->
		<button
			type={ buttonType }
			class={ classes }
			disabled?={ props.Disabled }
			{ props.Attributes... }
		>
			{ children... }
		</button>
	}
}
```

**Conditional Patterns:**

- Standard `if/else` blocks
- Nested conditionals
- Inline conditionals in attributes
- Guard clauses for early returns

---

## For Loops

Iterate over slices and maps with standard Go syntax:

```go
templ renderSelectOptions(selectID string, options []SelectOption) {
	{{
		// Group options by their Group field
		groupedOptions := make(map[string][]SelectOption)
		ungroupedOptions := []SelectOption{}

		for _, option := range options {
			if option.Group != "" {
				groupedOptions[option.Group] = append(groupedOptions[option.Group], option)
			} else {
				ungroupedOptions = append(ungroupedOptions, option)
			}
		}

		globalIndex := 0
	}}

	<!-- Render ungrouped options first -->
	if len(ungroupedOptions) > 0 {
		for _, option := range ungroupedOptions {
			@SelectItem(SelectItemProps{
				ID:       selectID,
				Value:    option.Value,
				Index:    globalIndex,
				Disabled: option.Disabled,
			}) {
				{ option.Label }
			}
			{{ globalIndex++ }}
		}
	}

	<!-- Render grouped options -->
	for groupName, groupOptions := range groupedOptions {
		@SelectGroup(SelectGroupProps{}) {
			if groupName != "" {
				@SelectLabel(SelectLabelProps{}) {
					{ groupName }
				}
			}
			for _, option := range groupOptions {
				@SelectItem(SelectItemProps{
					ID:       selectID,
					Value:    option.Value,
					Index:    globalIndex,
					Disabled: option.Disabled,
				}) {
					{ option.Label }
				}
				{{ globalIndex++ }}
			}
		}
	}
}
```

**Loop Features:**

- Range over slices: `for _, item := range items`
- Range over maps: `for key, value := range map`
- Access to index and value
- Nested loops supported
- Complex data processing in Go blocks

---

## Template Composition

Build complex components from simpler ones:

```go
// Card system with multiple composable parts
templ Card(props CardProps) {
	{{
		classes := cardVariants(props.Class)
	}}
	<div
		data-slot="card"
		class={ classes }
		{ props.Attributes... }
	>
		{ children... }
	</div>
}

templ CardHeader(props CardHeaderProps) {
	{{
		classes := cardHeaderVariants(props.Class)
	}}
	<div
		data-slot="card-header"
		class={ classes }
		{ props.Attributes... }
	>
		{ children... }
	</div>
}

templ CardContent(props CardContentProps) {
	{{
		classes := cardContentVariants(props.Class)
	}}
	<div
		data-slot="card-content"
		class={ classes }
		{ props.Attributes... }
	>
		{ children... }
	</div>
}

// Usage:
// @Card(CardProps{}) {
//   @CardHeader(CardHeaderProps{}) {
//     @CardTitle(CardTitleProps{}) { "Title" }
//   }
//   @CardContent(CardContentProps{}) {
//     <p>Content goes here</p>
//   }
// }
```

**Composition Patterns:**

- `{ children... }` for slot-based composition
- `@ComponentName()` for component calls
- Nested component structures
- Props-based customization
- `data-slot` attributes for CSS targeting

---

## Attributes and Spreading

Handle dynamic attributes elegantly:

```go
templ Button(props ButtonProps) {
	<button
		type={ props.Type }
		class={ classes }
		disabled?={ props.Disabled }
		{ props.Attributes... }
	>
		{ children... }
	</button>
}

// Advanced attribute handling
templ Checkbox(props CheckboxProps) {
	{{
		signalRef := signals.Signal("checked")
		// Datastar object syntax for conditional classes
		activeClassesObj := "{'bg-primary': " + signalRef + ", 'text-white': " + signalRef + "}"
	}}
	<button
		class={ baseClasses }
		data-class={ activeClassesObj }
		data-attr-aria-checked={ signalRef + " ? 'true' : 'false'" }
		{ props.Attributes... }
	>
		{ children... }
	</button>
}
```

**Attribute Features:**

- `{ props.Attributes... }` spreads all attributes
- `attribute?={ bool }` for conditional presence
- `data-attr-*` for dynamic attribute names
- `data-class` for conditional CSS classes
- Merge attributes with `utils.MergeAttributes()`

---

## Comments

Use HTML comments for documentation and Go comments for code:

```go
templ Checkbox(props CheckboxProps) {
	{{
		// Create signals using the new structured system
		signals := utils.Signals(props.ID, CheckboxSignals{
			Checked: props.Checked,
		})

		// Generate the CSS classes using the variant system
		classes := checkboxVariants(props.ClassName)
	}}
	<div data-signals={ signals.DataSignals }>
		<!-- Visible checkbox button -->
		<button
			type="button"
			role="checkbox"
			{ props.Attributes... }
		>
			<!-- Checkmark icon -->
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
				<path d="M20 6 9 17l-5-5"></path>
			</svg>
		</button>
		<!-- Hidden input for form submission -->
		<input
			type="checkbox"
			name={ props.Name }
			class="sr-only"
			tabindex="-1"
		/>
	</div>
}
```

---

## State Management Integration

Integrate with Datastar for reactive state:

```go
templ Tabs(props TabsProps) {
	{{
		// Create signals using the structured system
		signals := utils.Signals(tabsID, TabsSignals{
			Active: activeTab,
		})
	}}
	<div
		data-signals={ signals.DataSignals }
		class={ classes }
	>
		{ children... }
	</div>
}

templ TabsTrigger(props TabsTriggerProps) {
	{{
		signals := utils.Signals(props.ID, TabsSignals{})
		signalRef := signals.Signal("active")
		clickExpr := signals.Set("active", "'"+props.Value+"'")

		// Conditional classes using Datastar object syntax
		activeClassesObj := "{'bg-background': " + signalRef + " === '" + props.Value + "'}"
	}}
	<button
		class={ baseClasses }
		data-class={ activeClassesObj }
		data-on-click={ clickExpr }
		data-attr-data-state={ signalRef + " === '" + props.Value + "' ? 'active' : 'inactive'" }
	>
		{ children... }
	</button>
}
```

**State Management Features:**

- `data-signals` for state initialization
- `data-on-*` for event handlers
- `data-attr-*` for reactive attributes
- `data-class` for conditional styling
- Signal expressions for reactivity

---

## Advanced Patterns

Complex components with multiple concerns:

```go
templ Select(props SelectProps) {
	{{
		// Generate unique ID if not provided
		selectID := props.ID
		if selectID == "" {
			b := make([]byte, 4)
			rand.Read(b)
			selectID = fmt.Sprintf("select_%x", b)
		}

		// Determine initial state
		initialValue := props.Value
		if initialValue == "" {
			initialValue = props.DefaultValue
		}

		// Find initial label from options
		initialLabel := ""
		if initialValue != "" && len(props.Options) > 0 {
			for _, option := range props.Options {
				if option.Value == initialValue {
					initialLabel = option.Label
					break
				}
			}
		}

		// Create comprehensive state
		signals := utils.Signals(selectID, SelectSignals{
			Open:        props.DefaultOpen,
			Value:       initialValue,
			Label:       initialLabel,
			Highlighted: -1,
		})

		// Event handlers
		clickOutsideHandler := signals.Signal("open") + " ? " + signals.Set("open", "false") + " : null"
	}}
	<div
		data-select-id={ selectID }
		data-signals={ signals.DataSignals }
		data-on-click__outside={ clickOutsideHandler }
		class={ selectVariants(props.Class) }
	>
		if props.Name != "" {
			<!-- Hidden form input -->
			<input
				type="hidden"
				name={ props.Name }
				data-bind={ signals.Signal("value") }
				if props.Required {
					required
				}
			/>
		}

		if len(props.Options) > 0 {
			<!-- Auto-render when options provided -->
			@SelectTrigger(SelectTriggerProps{ID: selectID}) {
				@SelectValue(SelectValueProps{
					ID: selectID,
					Placeholder: props.Placeholder,
				})
			}
			@SelectContent(SelectContentProps{ID: selectID}) {
				@renderSelectOptions(selectID, props.Options)
			}
		} else {
			<!-- Manual composition -->
			{ children... }
		}
	</div>
}
```

**Advanced Features:**

- Automatic vs manual composition
- Complex state management
- Event handling patterns
- Form integration
- Accessibility attributes
- Performance optimization

---

## Best Practices

1. **Type Safety**: Always define prop types with documentation
2. **Composition**: Prefer small, composable components
3. **State Management**: Use structured signals for complex state
4. **Accessibility**: Include ARIA attributes and semantic HTML
5. **Performance**: Minimize Go code in render loops
6. **Documentation**: Comment complex logic and component APIs
7. **Error Handling**: Handle edge cases in Go blocks
8. **Naming**: Use consistent naming conventions for components and props

---

## Running Templ

Generate Go code from your templ files:

```bash
# Install templ
go install github.com/a-h/templ/cmd/templ@latest

# Generate code
templ generate

# Watch for changes
templ generate --watch
```

**Integration:**

- Generated code is type-safe Go
- Components implement `templ.Component` interface
- Use with any Go web framework
- Perfect for server-side rendering
- Integrates with build tools and CI/CD
