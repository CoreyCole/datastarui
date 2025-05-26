# DatastarUI Component Development Guide

## Project Overview

DatastarUI is a Go/templ port of shadcn/ui components that maintains pixel-perfect visual and behavioral parity while eliminating JavaScript dependencies (except for the 15KB Datastar library for reactivity).

## Component Architecture

### File Structure

Each component follows this structure:

```
components/[component-name]/
├── [component-name].templ     # Main component template
├── types.go                   # Component props and types
└── variants.go                # CSS class generation (CVA-style)
```

### Reference Implementation: Button Component

Use the button component as the gold standard for all implementations:

#### 1. Source Material (`ui/apps/v4/registry/new-york-v4/ui/button.tsx`)

- Always use the **New York style** components from the registry
- Extract the exact `buttonVariants` cva configuration
- Preserve all CSS classes, variants, and default values
- Maintain the same prop interface and behavior

#### 2. Component Template (`components/button/button.templ`)

```go
package button

templ Button(props ButtonProps) {
    {{
        // Set default values
        buttonType := props.Type
        if buttonType == "" {
            buttonType = "button"
        }

        // Generate CSS classes using variant system
        classes := buttonVariants(props.Variant, props.Size, props.Class)
    }}
    if props.AsChild {
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

#### 3. Types Definition (`components/button/types.go`)

```go
package button

import "github.com/a-h/templ"

type ButtonProps struct {
    // Mirror all React component props
    Variant    string           // "default", "destructive", "outline", "secondary", "ghost", "link"
    Size       string           // "default", "sm", "lg", "icon"
    AsChild    bool             // For composition patterns
    Class      string           // Additional CSS classes
    Attributes templ.Attributes // HTML attributes
    Disabled   bool             // Interactive state
    Type       string           // Button type attribute
}
```

#### 4. Variants System (`components/button/variants.go`)

```go
package button

import "github.com/coreycole/datastarui/templui/utils"

func buttonVariants(variant, size, className string) string {
    // Extract EXACT base classes from shadcn/ui cva
    baseClasses := "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive"

    // Map EXACT variant classes from cva
    variantClasses := map[string]string{
        "default":     "bg-primary text-primary-foreground shadow-xs hover:bg-primary/90",
        "destructive": "bg-destructive text-destructive-foreground shadow-xs hover:bg-destructive/90",
        "outline":     "border border-input bg-background shadow-xs hover:bg-accent hover:text-accent-foreground",
        "secondary":   "bg-secondary text-secondary-foreground shadow-xs hover:bg-secondary/80",
        "ghost":       "hover:bg-accent hover:text-accent-foreground",
        "link":        "text-primary underline-offset-4 hover:underline",
    }

    // Map EXACT size classes from cva
    sizeClasses := map[string]string{
        "default": "h-9 px-4 py-2 has-[>svg]:px-3",
        "sm":      "h-8 rounded-md gap-1.5 px-3 has-[>svg]:px-2.5",
        "lg":      "h-10 rounded-md px-6 has-[>svg]:px-4",
        "icon":    "size-9",
    }

    // Apply defaults and merge classes
    if variant == "" { variant = "default" }
    if size == "" { size = "default" }

    classes := []string{baseClasses}
    if variantClasses[variant] != "" {
        classes = append(classes, variantClasses[variant])
    }
    if sizeClasses[size] != "" {
        classes = append(classes, sizeClasses[size])
    }
    if className != "" {
        classes = append(classes, className)
    }

    return utils.TwMerge(classes...)
}
```

## Implementation Requirements

### 1. Pixel-Perfect Accuracy

- **Copy CSS classes exactly** from the shadcn/ui source
- **Preserve all Tailwind utilities** including complex selectors like `[&_svg]:size-4`
- **Maintain exact spacing, sizing, and visual hierarchy**
- **Support both light and dark themes** using CSS variables

### 2. Behavioral Parity

- **Mirror all React component props** and their default values
- **Implement the same composition patterns** (AsChild, Slot-like behavior)
- **Preserve accessibility attributes** and ARIA states
- **Handle disabled states** and focus management identically

### 3. Datastar Integration

- **Use Datastar signals** for reactive state: `data-signals-count="0"`
- **Handle events** with Datastar: `data-on-click="$count++"`
- **Bind reactive content**: `data-text="$count"`, `data-show="$visible"`

### 4. Go/templ Patterns

- **Use templ.Attributes** for flexible HTML attribute passing
- **Implement conditional rendering** with templ's `switch` statements

```go
templ userTypeDisplay(userType string) {
	switch userType {
		case "test":
			<span>{ "Test user" }</span>
		case "admin":
			<span>{ "Admin user" }</span>
		default:
			<span>{ "Unknown user" }</span>
	}
}
```

- **Handle optional props** with Go's zero-value defaults
- **Use the utils.TwMerge** function for class composition

## Development Environment

### templ templates

- when making changes to templ files, be sure to run `templ generate` to check for any compile errors

### Tailwind CSS Watch Process

**IMPORTANT**: Tailwind CSS is running in watch mode during development. The developer has `tailwindcss --watch` running automatically, which means:

- **DO NOT run Tailwind commands** unless specifically checking for CSS compilation errors
- CSS changes in `static/css/index.css` are automatically compiled to `static/css/build.css`
- The watch process monitors all `.templ` files for class changes
- Only run `tailwind` or similar commands if there are CSS compilation issues that need debugging

## Development Workflow

### Step 1: Analyze Source Component

1. Locate the component in `ui/apps/v4/registry/new-york-v4/ui/[component].tsx`
2. Extract the `cva` configuration (base classes, variants, defaults)
3. Identify all props and their types
4. Note any special behaviors or composition patterns

### Step 2: Create Component Structure

1. Create the component directory: `components/[component-name]/`
2. Define types in `types.go` with all necessary props
3. Implement variants in `variants.go` with exact CSS classes
4. Build the template in `[component-name].templ`

### Step 3: Add Interactivity

1. **Identify reactive needs**: What state changes are required?
2. **Implement Datastar patterns**: Use appropriate signals and bindings
3. **Test interactions**: Ensure smooth, JavaScript-like behavior
4. **Optimize performance**: Minimize unnecessary re-renders

## CSS Variables and Theming

The project uses shadcn/ui's CSS variable system for theming. Ensure all components reference these variables:

```css
:root {
  --background: 0 0% 100%;
  --foreground: 240 10% 3.9%;
  --primary: 240 5.9% 10%;
  --primary-foreground: 0 0% 98%;
  /* ... all shadcn/ui variables ... */
}

.dark {
  --background: 240 10% 3.9%;
  --foreground: 0 0% 98%;
  /* ... dark mode overrides ... */
}
```

## Quality Standards

### Visual Fidelity

- ✅ Identical appearance to shadcn/ui in light mode
- ✅ Identical appearance to shadcn/ui in dark mode
- ✅ Proper hover, focus, and active states
- ✅ Correct spacing and typography

### Functional Parity

- ✅ All props work as expected
- ✅ Composition patterns (AsChild) function correctly
- ✅ Accessibility attributes are preserved
- ✅ Keyboard navigation works identically

### Code Quality

- ✅ Clean, readable Go code
- ✅ Proper error handling and defaults
- ✅ Efficient class generation
- ✅ Comprehensive type definitions

## Example Usage

```go
// Basic usage
@button.Button(button.ButtonProps{})
    Click me
@button.Button

// With variants
@button.Button(button.ButtonProps{
    Variant: "destructive",
    Size: "lg",
})
    Delete
@button.Button

// With Datastar reactivity
@button.Button(button.ButtonProps{
    Attributes: templ.Attributes{
        "data-on-click": "$count++",
    },
})
    <span data-text="$count"></span>
@button.Button
```

This approach ensures every component is a faithful, pixel-perfect recreation of shadcn/ui while leveraging Go's type safety and Datastar's reactivity.
