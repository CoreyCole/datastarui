# DatastarUI

A Go/templ port of [shadcn/ui](https://ui.shadcn.com/) components that maintains pixel-perfect visual and behavioral parity with minimal JavaScript (lightweight 15KB [Datastar](https://data-star.dev/) library for reactivity).

See [datastar-ui.com](https://datastar-ui.com) for component demos.

## ✨ Features

- 🚀 **Server-side rendered** components with Go/templ
- ⚡ **Reactive UI** powered by Datastar signals
- 🎨 **Identical styling** to shadcn/ui using Tailwind CSS
- 📦 **Lightweight** - only 15KB Datastar runtime
- 🔧 **Type-safe** component args with Go structs
- 🌙 **Dark mode** support built-in
- ♿ **Accessible** with proper ARIA attributes

## 🚀 Quick Start

### Prerequisites

- [Go 1.24+](https://golang.org/dl/)
- [Just](https://github.com/casey/just) - command runner
- [Air](https://github.com/cosmtrek/air) - live reload for Go
- [templ](https://templ.guide/) - Go templating engine
- [tailwindcss standalone CLI](https://tailwindcss.com/blog/standalone-cli)

### Development Setup

```bash
# start the Tailwind CSS watcher:
just tailwind

# start the Go server with live reload:
just watch
```

see demo site at [http://localhost:4242](http://localhost:4242)

The development environment will automatically:

- ✅ Rebuild Go templates when `.templ` files change
- ✅ Recompile CSS when Tailwind classes are added/removed
- ✅ Restart the server when Go code changes

## 🏗️ Project Structure

```
datastarui/
├── components/          # Reusable UI components
│   ├── button/          # Button component
│   │   ├── button.templ # Template file
│   │   ├── types.go     # Args and types
│   │   └── variants.go  # CSS class variants
│   └── select/          # Select component (fully refactored)
├── utils/               # Utility libraries
│   ├── signals.go       # Signal management with namespacing
│   ├── expressions.go   # Datastar expression builders
│   └── data_class.go    # Conditional CSS class helpers
├── docs/                # Documentation
│   ├── guide.md         # Development patterns guide
│   ├── debugging.md     # Testing and troubleshooting
│   ├── playwright.md    # Browser automation testing
│   ├── signals.md       # Signal management reference
│   └── templ.md         # templ syntax guide
├── pages/               # Page templates
│   └── components/      # Component demo pages
├── layouts/             # Layout templates
├── static/              # Static assets
│   └── css/             # Tailwind CSS files
└── main.go              # Server entry point
```

## 🧩 Component Architecture

Each component follows a consistent 3-file pattern with utility-driven Datastar integration:

### 1. Template File (`component.templ`)

```go
package button

import "github.com/coreycole/datastarui/utils"

type ButtonSignals struct {
    Clicked bool `json:"clicked"`
    Loading bool `json:"loading"`
}

templ Button(args ButtonProps) {
    {{
        // Use utilities for signal management
        signals := utils.Signals(args.ID, ButtonSignals{
            Clicked: false,
            Loading: args.Loading,
        })
        
        // Use expression builders for click handlers
        clickHandler := utils.NewExpression().
            Statement("evt.preventDefault()").
            SetSignal("clicked", "true").
            Build()
    }}
    <button
        type={ args.Type }
        class={ buttonVariants(args.Variant, args.Size, args.Class) }
        disabled?={ args.Disabled }
        data-signals={ signals.DataSignals }
        data-on-click={ clickHandler }
        { args.Attributes... }
    >
        { children... }
    </button>
}
```

### 2. Args Definition (`types.go`)

```go
package button

import "github.com/a-h/templ"

type ButtonProps struct {
    ID         string           // Component identifier for signals
    Variant    string           // "default", "destructive", "outline"
    Size       string           // "default", "sm", "lg", "icon"
    Class      string           // Additional CSS classes
    Attributes templ.Attributes // HTML attributes
    Disabled   bool             // Interactive state
    Loading    bool             // Loading state
    Type       string           // Button type
}
```

### 3. CSS Variants (`variants.go`)

```go
package button

import "github.com/coreycole/datastarui/utils"

func buttonVariants(variant, size, className string) string {
    // Exact CSS classes from shadcn/ui
    baseClasses := "inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"

    variantClasses := map[string]string{
        "default":     "bg-primary text-primary-foreground hover:bg-primary/90",
        "destructive": "bg-destructive text-destructive-foreground hover:bg-destructive/90",
        "outline":     "border border-input bg-background hover:bg-accent hover:text-accent-foreground",
        "secondary":   "bg-secondary text-secondary-foreground hover:bg-secondary/80",
        "ghost":       "hover:bg-accent hover:text-accent-foreground",
        "link":        "text-primary underline-offset-4 hover:underline",
    }

    sizeClasses := map[string]string{
        "default": "h-10 px-4 py-2",
        "sm":      "h-9 rounded-md px-3",
        "lg":      "h-11 rounded-md px-8",
        "icon":    "h-10 w-10",
    }

    return utils.TwMerge(baseClasses, variantClasses[variant], sizeClasses[size], className)
}
```

## 🎨 Design System

The project uses the same design tokens as shadcn/ui:

- **Colors**: CSS custom properties for theming
- **Typography**: Tailwind's font system
- **Spacing**: Consistent spacing scale
- **Shadows**: Subtle elevation system
- **Borders**: Rounded corners and borders

## 🤝 Contributing

1. **Pick a component** from the [shadcn/ui registry](https://ui.shadcn.com/docs/components)
2. **Follow the utility-driven architecture** using expression builders
3. **Maintain pixel-perfect accuracy** to the original
4. **Use structured Datastar integration** with utilities
5. **Create comprehensive demos** showing all variants

See the [Development Guide](./docs/guide.md) for detailed implementation patterns.

## 📖 Documentation

- [Development Guide](./docs/guide.md) - Comprehensive utility patterns and best practices
- [Signals Management](./docs/signals.md) - Signal management patterns (legacy reference)
- [Debugging Guide](./docs/debugging.md) - Testing and troubleshooting components
- [templ Syntax Guide](./docs/templ.md) - Complete templ templating reference
- [Playwright Testing](./docs/playwright.md) - Browser automation testing
- [Datastar Documentation](https://data-star.dev/) - Reactivity framework
- [templ Documentation](https://templ.guide/) - Go templating engine
- [shadcn/ui](https://ui.shadcn.com/) - Original component library

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
