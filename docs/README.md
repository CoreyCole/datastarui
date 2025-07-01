# DatastarUI Documentation

Welcome to the DatastarUI documentation! This directory contains comprehensive guides for building Go/templ components with Datastar integration.

## 📚 Documentation Overview

### Core Guides

- **[Development Guide](./guide.md)** - **START HERE** - Comprehensive patterns for building components with utility libraries
- **[Debugging Guide](./debugging.md)** - Testing, troubleshooting, and ensuring pixel-perfect accuracy
- **[templ Syntax Guide](./templ.md)** - Complete reference for templ templating engine

### Testing & Automation

- **[Playwright Testing](./playwright.md)** - Browser automation testing with Playwright MCP
- **[Playwright Commands](./playwright-commands.md)** - Quick reference for testing commands

### Planning & Process

- **[Refactoring Plan](./refactoring-plan.md)** - Systematic approach to modernizing all components
- **[Component Checklist](./component-checklist.md)** - Template for tracking refactoring progress

### Reference

- **[Signals Management](./signals.md)** - Legacy signal patterns (use guide.md for modern approach)

## 🚀 Quick Start

1. **New to DatastarUI?** Start with the [Development Guide](./guide.md)
2. **Building components?** Follow the utility patterns in the guide
3. **Debugging issues?** Check the [Debugging Guide](./debugging.md)
4. **Writing tests?** Use [Playwright Testing](./playwright.md)

## 🛠️ Modern Development Approach

DatastarUI uses utility libraries for clean, maintainable Datastar integration:

### Core Utilities

```go
// Signal management
signals := utils.Signals("component_id", ComponentSignals{...})

// Expression building  
expr := utils.NewExpression().
    Statement("evt.preventDefault()").
    SetSignal("modal.open", "false").
    Build()

// Conditional classes
dataClass := utils.HighlightedItem("select.highlighted", index)
```

### Key Benefits

- ✅ **No string concatenation** - Type-safe expression building
- ✅ **Reusable patterns** - Standardized handlers for common UI patterns
- ✅ **Error prevention** - Proper escaping and syntax handling
- ✅ **Maintainability** - Clean, readable code

## 📖 Documentation Principles

Our documentation follows these principles:

1. **Example-driven** - Show real code from working components
2. **Progressive complexity** - Start simple, build up to advanced patterns
3. **Best practices** - Highlight the recommended approach
4. **Common pitfalls** - Explain what to avoid and why

## 🔄 Migration from Legacy Patterns

If you're working with older DatastarUI code that uses string concatenation:

1. Read the [Development Guide](./guide.md) for modern patterns
2. Use [Signals Management](./signals.md) as a reference for legacy code
3. Follow the migration examples in the guide
4. Test thoroughly using [Debugging Guide](./debugging.md)

## 🤝 Contributing to Documentation

When adding or updating documentation:

1. **Follow existing patterns** - Use the same structure and formatting
2. **Include examples** - Show working code from actual components
3. **Test examples** - Ensure all code examples actually work
4. **Cross-reference** - Link to related documentation
5. **Update this index** - Add new docs to the overview above

## 📋 Documentation Checklist

For each new component or pattern:

- [ ] Implementation follows [Development Guide](./guide.md) patterns
- [ ] Uses utility libraries instead of string concatenation  
- [ ] Has proper signal management with `utils.Signals()`
- [ ] Uses expression builders for complex Datastar expressions
- [ ] Includes tests following [Playwright Testing](./playwright.md)
- [ ] Documentation updated in relevant guides

## 🎯 Next Steps

Ready to build components? Head to the [Development Guide](./guide.md) and start with the utility patterns!