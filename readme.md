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
│   │   ├── args.go      # Component arguments
│   │   └── variants.go  # CSS class variants
│   └── select/          # Select component (fully refactored)
├── utils/               # Utility libraries
│   ├── signals.go       # Signal management with namespacing
│   ├── expressions.go   # Datastar expression builders
│   └── data_class.go    # Conditional CSS class helpers
├── pages/               # Page templates
│   └── components/      # Component demo pages
└── main.go              # Server entry point
```

## 🧩 Component Architecture

Each component follows a pattern with utility-driven Datastar integration:

- [Template File (`dialog.templ`)](./components/dialog/dialog.templ)
- [Expressions (`expressions.go`)](./components/dialog/expressions.go)
- [Args Definition (`args.go`)](./components/dialog/args.go)
- [CSS Variants (`variants.go`)](./components/dialog/variants.go)

### Template File (`dialog.templ`)

```go
package dialog

import "github.com/coreycole/datastarui/utils"

// DialogSignals defines the signal structure for dialog components
type DialogSignals struct {
	Open bool `json:"open"`
}

// Dialog container - pure Datastar signal-based modal using data-show
templ Dialog(args DialogArgs) {
	{{
		// Create signals using the new structured system with proper initial state
		signals := utils.Signals(args.ID, DialogSignals{
			Open: args.DefaultOpen,
		})

		// Dialog backdrop overlay
		backdropClasses := "fixed inset-0 z-50 bg-black/50 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0"

		// Dialog positioning classes - centered on screen
		dialogClasses := "fixed left-[50%] top-[50%] z-50 w-full max-w-lg translate-x-[-50%] translate-y-[-50%] data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95"

		// Generate the CSS classes for the inner dialog content container using our variant system
		containerClasses := DialogVariants(args)

		// Create dialog handler for clean expressions
		dialogHandler := NewDialogHandler(signals)
		backdropClickHandler := dialogHandler.BuildBackdropClickHandler()
		escapeHandler := dialogHandler.BuildEscapeHandler()
	}}
	<div data-signals={ signals.DataSignals }>
		<!-- Dialog backdrop overlay -->
		<div
			data-show={ signals.Signal("open") }
			class={ backdropClasses }
			data-on-click={ backdropClickHandler }
			data-on-keydown__window={ escapeHandler }
			if !args.DefaultOpen {
				style="display: none;"
			}
		>
			<!-- Dialog content container -->
			<div
				id={ args.ID }
				class={ dialogClasses }
				role="dialog"
				aria-modal="true"
				tabindex="-1"
				data-on-click="evt.stopPropagation()"
				data-on-mount="evt.target.focus()"
			>
				<div class={ containerClasses }>
					{ children... }
				</div>
			</div>
		</div>
	</div>
}
```

### Expressions (`expressions.go`)

```go
package dialog

import (
	"fmt"
	"github.com/coreycole/datastarui/utils"
)

// DialogHandler creates handlers for Dialog component functionality
type DialogHandler struct {
	signals *utils.SignalManager
}

// NewDialogHandler creates a dialog handler
func NewDialogHandler(signals *utils.SignalManager) *DialogHandler {
	return &DialogHandler{
		signals: signals,
	}
}

// BuildBackdropClickHandler creates the backdrop click handler for closing dialog
func (d *DialogHandler) BuildBackdropClickHandler() string {
	return d.signals.ConditionalAction("evt.target === evt.currentTarget", "open", "false")
}

// BuildEscapeHandler creates an escape key handler for closing dialog
func (d *DialogHandler) BuildEscapeHandler() string {
	condition := fmt.Sprintf("evt.key === 'Escape' && %s", d.signals.Signal("open"))
	return d.signals.ConditionalAction(condition, "open", "false")
}

// BuildCloseHandler creates a close handler with optional return value
func (d *DialogHandler) BuildCloseHandler(returnValue string) string {
	expr := utils.NewExpression().Statement(d.signals.Set("open", "false"))
	
	if returnValue != "" {
		expr.Statement(d.signals.SetString("returnValue", returnValue))
	}
	
	return expr.Build()
}
```

### Args Definition (`args.go`)

```go
package dialog

import "github.com/a-h/templ"

// DialogArgs defines the args for the Dialog container (using Datastar signals)
type DialogArgs struct {
	ID          string
	DefaultOpen bool // Whether the dialog should be open by default
	Class       string
	Attributes  templ.Attributes
}

// DialogTriggerArgs defines the args for the DialogTrigger component
type DialogTriggerArgs struct {
	DialogID   string
	AsChild    bool
	Class      string
	Attributes templ.Attributes
}
```

### CSS Variants (`variants.go`)

```go
package dialog

import (
	"github.com/coreycole/datastarui/utils"
)

// DialogVariants returns the CSS classes for the main Dialog container component
func DialogVariants(args DialogArgs) string {
	// Dialog-specific styling - optimized for modal dialogs with consistent padding
	baseClasses := "max-w-lg w-full max-h-[90vh] overflow-auto bg-background border shadow-lg rounded-lg p-6"

	return utils.TwMerge(baseClasses, args.Class)
}
```

## 🎨 Design System

The project uses tailwind classes from [shadcn/ui](https://ui.shadcn.com/) components new york v4.

## 🤝 Contributing

1. **Pick a component** from the [shadcn/ui registry](https://ui.shadcn.com/docs/components)
1. **Follow the utility-driven architecture** using expression builders
1. **Create comprehensive demos** showing all variants

## 📖 Documentation

- [Playwright Testing](./docs/playwright.md) - Browser automation testing
- [Datastar Documentation](https://data-star.dev/) - Reactivity framework
- [templ Documentation](https://templ.guide/) - Go templating engine
- [shadcn/ui](https://ui.shadcn.com/) - Original component library

## 📄 License

MIT License - see [LICENSE](./LICENSE) for details.
