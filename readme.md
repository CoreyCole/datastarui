# DatastarUI

A Go/templ port of [shadcn/ui](https://ui.shadcn.com/) components that maintains pixel-perfect visual and behavioral parity with minimal JavaScript (lightweight 15KB [Datastar](https://data-star.dev/) library for reactivity).

See [datastar-ui.com](https://datastar-ui.com) for component demos.

## ✨ Features

- 🚀 **Server-side rendered** components with Go/templ
- ⚡ **Reactive UI** powered by Datastar signals
- 🎨 **Identical styling** to shadcn/ui using Tailwind CSS
- 📦 **Lightweight** - only 15KB Datastar runtime
- 🔧 **Type-safe** component props with Go structs
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
│   │   ├── types.go     # Props and types
│   │   └── variants.go  # CSS class variants
│   └── dropdown/        # Dropdown component
├── pages/               # Page templates
│   └── components/      # Component demo pages
├── layouts/             # Layout templates
├── static/              # Static assets
│   └── css/             # Tailwind CSS files
└── main.go              # Server entry point
```

## 🧩 Component Architecture

Each component follows a consistent pattern:

### 1. Template File (`component.templ`)

```go
package button

templ Button(props ButtonProps) {
    <button
        type={ props.Type }
        class={ buttonVariants(props.Variant, props.Size, props.Class) }
        disabled?={ props.Disabled }
        { props.Attributes... }
    >
        { children... }
    </button>
}
```

### 2. Props Definition (`props.go`)

```go
package button

type ButtonProps struct {
    Variant    string           // "default", "destructive", "outline"
    Size       string           // "default", "sm", "lg", "icon"
    Class      string           // Additional CSS classes
    Attributes templ.Attributes // HTML attributes
    Disabled   bool             // Interactive state
    Type       string           // Button type
}
```

### 3. CSS Variants (`variants.go`)

```go
package button

func buttonVariants(variant, size, className string) string {
    // Exact CSS classes from shadcn/ui
    baseClasses := "inline-flex items-center justify-center..."

    variantClasses := map[string]string{
        "default": "bg-primary text-primary-foreground...",
        "outline": "border border-input bg-background...",
    }

    return utils.TwMerge(baseClasses, variantClasses[variant], className)
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
1. **Follow the architecture** outlined in the development guide
1. **Maintain pixel-perfect accuracy** to the original
1. **Add Datastar reactivity** where needed
1. **Create comprehensive demos** showing all variants

See the [Development Guide](./.cursor/rules/datastar.mdc) for detailed implementation instructions.

## 📖 Documentation

- [Component Development Guide](./.cursor/rules/datastar.mdc) - Detailed implementation patterns
- [Datastar Documentation](https://data-star.dev/) - Reactivity framework
- [templ Documentation](https://templ.guide/) - Go templating
- [shadcn/ui](https://ui.shadcn.com/) - Original component library

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
