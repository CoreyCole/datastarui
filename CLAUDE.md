# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DatastarUI is a Go/templ port of shadcn/ui components that maintains pixel-perfect visual and behavioral parity while using server-side rendering with minimal JavaScript (only the 15KB Datastar library for reactivity).

## Common Commands

### Development
- `just tailwind` - Start Tailwind CSS watcher (monitors .templ files for class changes)
- `just watch` - Start Go server with live reload using Air
- `just build` - Generate templ files and build binary
- `just run` - Build and run the binary
- `just install` - Install all dependencies (air, templ, go modules)

### Important Notes
- **DO NOT** run `go run main.go` or `go build` directly - use `just` commands
- **DO NOT** try to run the compiled binary - the developer has live reload running
- Templ files are auto-generated, but you can run `templ generate` to check for errors
- Tailwind CSS runs in watch mode automatically during development

## Architecture

### Tech Stack
- **Go 1.24+** with Echo framework for HTTP server
- **templ** - Go templating engine for type-safe HTML templates
- **Datastar** - 15KB JavaScript library for reactivity via data-* attributes
- **Tailwind CSS** - Utility-first CSS framework, exact classes from shadcn/ui
- **tailwind-merge-go** - For merging Tailwind classes safely

### Project Structure
```
datastarui/
├── components/                    # Reusable UI components
│   ├── button/
│   │   ├── button.templ          # Component template
│   │   ├── types.go              # Props and types
│   │   └── variants.go           # CSS variants
├── pages/
│   ├── components/               # Component demo pages
│   └── home_page.templ           # Home page
├── layouts/                      # Page layouts and navigation
├── static/css/                   # Tailwind CSS files
└── main.go                       # Server and routing
```

### Component Pattern
Each component follows a consistent 3-file pattern:
1. **Template** (`component.templ`) - templ markup with Datastar attributes
2. **Types** (`types.go`) - Go structs defining component props
3. **Variants** (`variants.go`) - CSS class generation matching shadcn/ui exactly

## Datastar Development Guidelines

### Signal Naming Convention
- Use lowercase with underscores in `props.ID` (e.g., `user_profile`, `item_list`)
- **Never** use uppercase, dashes, or periods in signal names
- Components should validate `props.ID` format and throw errors for invalid names

### Signal Architecture
- Signals are **globally scoped** on the page
- Use `props.ID` as namespace: `$user_profile.name`, `$item_list.selected`
- Only leaf nodes are valid signals (not intermediate namespaces)
- Initialize signals in component template: `data-signals='{"user_profile": {"name": "", "active": false}}'`

### Key Patterns
- **Props down, events up** - encapsulate state, communicate via defined interfaces
- **Server-driven state** - backend is single source of truth
- **Hypermedia approach** - let server determine available actions
- **Minimal JavaScript** - use data-* attributes for reactivity when possible

### Common Datastar Attributes
- `data-signals` - Initialize component state
- `data-on-click` - Handle click events
- `data-text` - Display signal values
- `data-show`/`data-hide` - Conditional visibility
- `data-attr-*` - Dynamic HTML attributes
- `data-bind-*` - Two-way data binding
- `data-indicator-*` - Loading states for fetch requests

## templ Specific Guidelines

### Required Patterns
- All page content must be inside `@l.Root()` layout wrapper
- Use `templ.SafeURL()` for href attributes to prevent XSS
- Use `templ.Attributes` for flexible HTML attribute passing
- Implement conditional rendering with `switch` statements

### Component Composition
- Nest components properly for multi-part components (e.g., breadcrumbs)
- Follow existing import patterns in main.go for routing
- Create comprehensive demo pages showing all component variants

## Development Workflow

### Adding New Components
1. **Analyze** source shadcn/ui component for structure and variants
2. **Create** component directory with 3-file pattern
3. **Implement** exact CSS classes using tailwind-merge-go
4. **Add** Datastar reactivity with proper signal patterns
5. **Create** demo page in `pages/components/[component]/`
6. **Update** routing in main.go and sidebar.go
7. **Test** all variants and interactive behaviors

### CSS and Styling
- Use exact CSS classes from shadcn/ui for pixel-perfect parity
- Leverage Tailwind CSS design tokens (colors, spacing, typography)
- Support dark mode with CSS custom properties
- Maintain accessibility with proper ARIA attributes

## TODO: Testing Strategy
Need to implement automated testing for Datastar component demo pages to verify:
- Component rendering across all variants
- Interactive behavior with Datastar signals  
- Visual regression testing
- Accessibility compliance