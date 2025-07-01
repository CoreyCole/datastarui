# Debugging and Component Analysis

This guide explains how to debug DatastarUI components and ensure pixel-perfect accuracy with shadcn/ui.

## HTML Comparison Workflow

When implementing components, compare the actual HTML output between DatastarUI and shadcn/ui to ensure pixel-perfect accuracy.

### Setup Comparison Environment

All HTML comparison files should be stored in `/d/cdev/shad-diffs` to keep the project root clean:

```bash
# Create comparison directory
mkdir -p /d/cdev/shad-diffs
cd /mnt/d/cdev
```

### Step 1: Capture HTML Files

```bash
# Capture DatastarUI HTML (dev server on port 4242)
curl -s http://localhost:4242/components/tabs > shad-diffs/datastarui-tabs.html

# Capture shadcn/ui HTML from live site
curl -s https://ui.shadcn.com/docs/components/tabs > shad-diffs/shadcn-tabs-live.html
```

### Step 2: Extract and Compare Patterns

```bash
# Compare CSS classes between implementations
echo "=== DatastarUI Classes ==="
cat shad-diffs/datastarui-tabs.html | rg 'class="[^"]*"' -o | head -20

echo "=== shadcn/ui Classes ==="
cat shad-diffs/shadcn-tabs.html | rg 'class="[^"]*"' -o | head -20

# Compare grid layout classes
echo "=== Grid Classes Comparison ==="
echo "DatastarUI:"
cat shad-diffs/datastarui-tabs.html | rg 'grid-cols-|grid-rows-' -o
echo "shadcn/ui:"
cat shad-diffs/shadcn-tabs.html | rg 'grid-cols-|grid-rows-' -o
```

### Step 3: Analyze Component Structure

```bash
# Find component wrapper elements
echo "=== DatastarUI Structure ==="
cat shad-diffs/datastarui-tabs.html | rg '<div[^>]*data-slot=' -A 2 -B 1

echo "=== shadcn/ui Structure ==="
cat shad-diffs/shadcn-tabs.html | rg '<div[^>]*class="[^"]*tabs[^"]*"' -A 2 -B 1

# Extract button elements
cat shad-diffs/datastarui-button.html | rg '<button' -A 3 -B 1
```

### Step 4: Advanced Pattern Analysis

```bash
# Find Tailwind utility classes
echo "=== Color Classes ==="
cat shad-diffs/datastarui-button.html | rg 'bg-[a-z-]+|text-[a-z-]+|border-[a-z-]+' -o

echo "=== Interactive States ==="
cat shad-diffs/datastarui-button.html | rg 'hover:|focus:|active:|disabled:' -o

echo "=== Layout Classes ==="
cat shad-diffs/datastarui-tabs.html | rg 'flex|grid|inline|block' -o

echo "=== Spacing Classes ==="
cat shad-diffs/datastarui-tabs.html | rg 'p-[0-9]|m-[0-9]|px-[0-9]|py-[0-9]|gap-[0-9]' -o
```

### Step 5: Component-Specific Analysis

Create comparison scripts for detailed analysis:

```bash
cat > shad-diffs/compare-component.sh << 'EOF'
#!/bin/bash
COMPONENT=$1

echo "=== $COMPONENT COMPARISON ==="
echo "DatastarUI structure:"
cat shad-diffs/datastarui-$COMPONENT.html | rg "data-slot=\"$COMPONENT\"" -A 5 -B 1

echo -e "\nshadcn/ui structure:"
cat shad-diffs/shadcn-$COMPONENT.html | rg "class=\"[^\"]*$COMPONENT[^\"]*\"" -A 5 -B 1

echo -e "\n=== ARIA ATTRIBUTES ==="
echo "DatastarUI ARIA:"
cat shad-diffs/datastarui-$COMPONENT.html | rg 'aria-[a-z-]+=' -o | sort | uniq

echo "shadcn/ui ARIA:"
cat shad-diffs/shadcn-$COMPONENT.html | rg 'aria-[a-z-]+=' -o | sort | uniq
EOF

chmod +x shad-diffs/compare-component.sh
```

## Debugging Datastar Expressions

### Console Error Detection

Use the Playwright MCP tools for runtime error detection:

```bash
# Navigate to component page
mcp__playwright__playwright_navigate http://localhost:4242/components/select

# Check for Datastar expression errors
mcp__playwright__playwright_console_logs type="all" clear=true

# Interact with component and check for errors
mcp__playwright__playwright_click "[data-select-id='test'] button"
mcp__playwright__playwright_console_logs type="error"
```

### Common Expression Issues

1. **Syntax Errors**: Missing quotes, parentheses, or semicolons
2. **Signal References**: Incorrect signal paths or undefined signals  
3. **Type Mismatches**: Comparing numbers with strings
4. **Conditional Logic**: Improper ternary operator usage

### Testing Expression Generation

Verify that utility builders generate correct expressions:

```go
// Test expression building
expr := utils.NewExpression().
    Statement("console.log('Testing')").
    SetSignal("test.value", "'hello'").
    Build()
// Should output: "console.log('Testing'); $test.value = 'hello'"

// Test conditional classes
dataClass := utils.HighlightedItem("select.highlighted", 0)
// Should output: {'bg-accent': $select.highlighted === 0, 'text-accent-foreground': $select.highlighted === 0}
```

## Development Server Debugging

### Check Compilation Status

Always verify that changes compile successfully:

```bash
# Check development server logs
just docker-tail app

# Look for:
# - "templ has changed"
# - "building..."
# - "(✓) Complete [updates=XX duration=XXXms]"
```

### Common Compilation Issues

1. **templ Syntax Errors**: Invalid templ syntax in .templ files
2. **Go Compilation Errors**: Type errors or missing imports
3. **CSS Generation**: Tailwind class detection issues

## Visual Regression Testing

### Screenshot Comparison

Use Playwright for visual testing:

```bash
# Take reference screenshots
mcp__playwright__playwright_navigate http://localhost:4242/components/button
mcp__playwright__playwright_screenshot name="button-variants"

# Compare different states
mcp__playwright__playwright_hover "button[variant='outline']"
mcp__playwright__playwright_screenshot name="button-outline-hover"
```

### Responsive Testing

```bash
# Test different viewport sizes
mcp__playwright__playwright_evaluate "window.resizeTo(768, 1024)"
mcp__playwright__playwright_screenshot name="button-tablet"

mcp__playwright__playwright_evaluate "window.resizeTo(375, 667)"
mcp__playwright__playwright_screenshot name="button-mobile"
```

## Debugging Checklist

### Component Implementation
- [ ] HTML structure matches shadcn/ui
- [ ] CSS classes are identical
- [ ] Interactive states work correctly
- [ ] ARIA attributes are present
- [ ] Responsive behavior is consistent

### Datastar Integration
- [ ] Signals are properly initialized
- [ ] Expressions compile without errors
- [ ] Event handlers work as expected
- [ ] State updates trigger re-renders
- [ ] No JavaScript console errors

### Development Environment
- [ ] Server is running and accessible
- [ ] templ files compile successfully
- [ ] Tailwind classes are generated
- [ ] Live reload is working
- [ ] Browser devtools show no errors

### Accessibility
- [ ] Keyboard navigation works
- [ ] Screen reader attributes
- [ ] Focus management
- [ ] Color contrast compliance
- [ ] Proper semantic markup

## Performance Debugging

### Datastar Performance

Monitor signal updates and expression evaluation:

```bash
# Enable Datastar debugging
mcp__playwright__playwright_evaluate "window.ds.debug = true"

# Check for excessive re-renders
mcp__playwright__playwright_console_logs type="log" search="datastar"
```

### Bundle Size Analysis

Check the impact of utilities on bundle size:

```bash
# Analyze generated templ files
find components -name "*_templ.go" -exec wc -l {} \; | sort -n

# Check expression complexity
rg "BuildClickHandler|BuildKeyboardHandler" components/ -A 5 -B 5
```

## Troubleshooting Guide

### Common Issues and Solutions

| Issue | Cause | Solution |
|-------|-------|----------|
| Console errors | Invalid Datastar expressions | Use expression builders |
| Classes not applying | Tailwind not detecting classes | Check CSS generation |
| No reactivity | Missing data-signals | Initialize signals properly |
| Wrong styling | Incorrect CSS classes | Compare with shadcn/ui |
| Keyboard navigation broken | Missing keyboard handlers | Use keyboard expression builders |

### Getting Help

1. **Check console logs** for JavaScript errors
2. **Verify HTML output** against shadcn/ui
3. **Test expressions** using browser devtools
4. **Compare CSS classes** using ripgrep
5. **Use Playwright MCP** for automated testing