# Playwright MCP Quick Reference

Quick reference for common Playwright MCP commands used during component testing and refactoring.

## 🚀 Setup Commands

```bash
# Start development environment
just up
just docker-tail app  # Verify server is running

# Check compilation status
just docker-tail app  # Look for "Complete [updates=XX]"
```

## 🌐 Navigation Commands

```bash
# Navigate to component page
mcp__playwright__playwright_navigate http://localhost:4242/components/button
mcp__playwright__playwright_navigate http://localhost:4242/components/select
mcp__playwright__playwright_navigate http://localhost:4242/components/calendar
```

## 🖱️ Interaction Commands

### Clicking Elements
```bash
# Click with specific selectors
mcp__playwright__playwright_click "button[type='button']"
mcp__playwright__playwright_click "[data-select-id='test'] [data-slot='select-trigger']"
mcp__playwright__playwright_click "[data-datepicker-id='single_date'] button[data-on-click*='open']"

# Click calendar day buttons (specific to avoid overlaps)
mcp__playwright__playwright_click "[data-datepicker-id='range_date'] [data-slot='datepicker-popover'] button:has-text('15')"
```

### Keyboard Interaction
```bash
# Press individual keys
mcp__playwright__playwright_press_key "ArrowDown"
mcp__playwright__playwright_press_key "ArrowUp"
mcp__playwright__playwright_press_key "Enter"
mcp__playwright__playwright_press_key "Escape"
mcp__playwright__playwright_press_key "Tab"
mcp__playwright__playwright_press_key " "  # Space

# Key combinations (if supported)
mcp__playwright__playwright_press_key "Shift+Tab"
```

### Form Input
```bash
# Fill input fields
mcp__playwright__playwright_fill "input[type='text']" "test value"
mcp__playwright__playwright_fill "[data-datepicker-id='single_date'] input" "12/25/2024"

# Clear input
mcp__playwright__playwright_fill "input[type='text']" ""
```

### Hover Actions
```bash
# Hover over elements
mcp__playwright__playwright_hover "button[variant='outline']"
mcp__playwright__playwright_hover ".hover-target"
```

## 📸 Visual Testing Commands

### Screenshots
```bash
# Full page screenshot
mcp__playwright__playwright_screenshot name="component-baseline"
mcp__playwright__playwright_screenshot name="calendar-month-view"
mcp__playwright__playwright_screenshot name="select-dropdown-open"

# Element-specific screenshot
mcp__playwright__playwright_screenshot name="button-variants" selector="[data-component='button']"

# Full page screenshot
mcp__playwright__playwright_screenshot name="responsive-mobile" fullPage=true
```

### Screen Resizing
```bash
# Change viewport size for responsive testing
mcp__playwright__playwright_evaluate "window.resizeTo(768, 1024)"  # Tablet
mcp__playwright__playwright_evaluate "window.resizeTo(375, 667)"   # Mobile
mcp__playwright__playwright_evaluate "window.resizeTo(1280, 720)"  # Desktop
```

## 🐛 Debugging Commands

### Console Monitoring
```bash
# Check all console messages
mcp__playwright__playwright_console_logs type="all" clear=true

# Check only errors
mcp__playwright__playwright_console_logs type="error"

# Check specific log types
mcp__playwright__playwright_console_logs type="warn"
mcp__playwright__playwright_console_logs type="log"
mcp__playwright__playwright_console_logs type="debug"

# Clear previous logs and get fresh output
mcp__playwright__playwright_console_logs type="all" clear=true
```

### DOM Inspection
```bash
# Get visible HTML structure
mcp__playwright__playwright_get_visible_html

# Get specific element HTML
mcp__playwright__playwright_get_visible_html selector="[data-slot='select-content']"

# Get HTML with size limit
mcp__playwright__playwright_get_visible_html maxLength=5000

# Get clean HTML (remove scripts, comments)
mcp__playwright__playwright_get_visible_html cleanHtml=true removeScripts=true
```

### JavaScript Evaluation
```bash
# Execute JavaScript
mcp__playwright__playwright_evaluate "console.log('Testing expression')"
mcp__playwright__playwright_evaluate "document.querySelectorAll('[data-signals]').length"

# Get element properties
mcp__playwright__playwright_evaluate "document.querySelector('button').disabled"
mcp__playwright__playwright_evaluate "getComputedStyle(document.querySelector('.select-item')).backgroundColor"
```

## 🎯 Component-Specific Selectors

### Select Component
```bash
# Trigger button
"[data-select-id='component_id'] [data-slot='select-trigger']"

# Dropdown content
"[data-select-id='component_id'] [data-slot='select-content']"

# Select items
"[data-select-id='component_id'] [data-select-item]"
"[data-select-id='component_id'] [data-select-item][data-value='option1']"
```

### DatePicker Component
```bash
# Calendar icon button
"[data-datepicker-id='picker_id'] button[data-on-click*='open']"

# Popover container
"[data-datepicker-id='picker_id'] [data-slot='datepicker-popover']"

# Calendar day buttons (when popover is visible)
"[data-datepicker-id='picker_id'] [data-slot='datepicker-popover'] button:has-text('15')"

# Input fields
"[data-datepicker-id='picker_id'] input[type='text']"
```

### Calendar Component
```bash
# Month navigation
"[data-slot='calendar'] button[data-on-click*='currentDate']"

# Day buttons
"[data-slot='calendar'] button[data-calendar-day]"

# Today button
"[data-slot='calendar'] button:has-text('Today')"
```

### Tabs Component
```bash
# Tab triggers
"[data-slot='tabs-list'] button[role='tab']"
"[data-slot='tabs-list'] button[data-value='tab1']"

# Tab content
"[data-slot='tabs-content'][data-value='tab1']"
```

### Dialog/Modal Component
```bash
# Trigger button
"[data-dialog-trigger]"

# Modal backdrop
"[data-slot='dialog-overlay']"

# Close button
"[data-slot='dialog-close']"
```

## 🔄 Common Testing Workflows

### Basic Component Test
```bash
# 1. Navigate
mcp__playwright__playwright_navigate http://localhost:4242/components/button

# 2. Check for errors
mcp__playwright__playwright_console_logs type="all" clear=true

# 3. Take baseline screenshot
mcp__playwright__playwright_screenshot name="button-baseline"

# 4. Test interactions
mcp__playwright__playwright_click "button[variant='default']"
mcp__playwright__playwright_hover "button[variant='outline']"

# 5. Check for errors again
mcp__playwright__playwright_console_logs type="error"

# 6. Take final screenshot
mcp__playwright__playwright_screenshot name="button-tested"
```

### Keyboard Navigation Test
```bash
# 1. Open component
mcp__playwright__playwright_click "[data-select-id='test'] button"

# 2. Test arrow keys
mcp__playwright__playwright_press_key "ArrowDown"
mcp__playwright__playwright_screenshot name="select-arrow-down"

mcp__playwright__playwright_press_key "ArrowUp"
mcp__playwright__playwright_screenshot name="select-arrow-up"

# 3. Test selection
mcp__playwright__playwright_press_key "Enter"
mcp__playwright__playwright_screenshot name="select-selected"

# 4. Check console
mcp__playwright__playwright_console_logs type="error"
```

### Error Detection Workflow
```bash
# 1. Clear console and navigate
mcp__playwright__playwright_console_logs type="all" clear=true
mcp__playwright__playwright_navigate http://localhost:4242/components/calendar

# 2. Interact with component
mcp__playwright__playwright_click "[data-calendar-day]"

# 3. Check for JavaScript errors
mcp__playwright__playwright_console_logs type="error"

# 4. Check for Datastar expression errors
mcp__playwright__playwright_console_logs type="all" search="GenerateExpression"
```

## 🧹 Cleanup

```bash
# Close browser when done
mcp__playwright__playwright_close
```

## 💡 Tips

1. **Use specific selectors** - Avoid generic selectors like `button` that match multiple elements
2. **Check console regularly** - Always check for errors after interactions
3. **Take screenshots liberally** - Visual verification is crucial
4. **Test keyboard navigation** - Ensure accessibility is preserved
5. **Clear console logs** - Use `clear=true` to get fresh error detection
6. **Wait for animations** - Some components may need brief delays after interactions