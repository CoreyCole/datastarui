## Debugging and HTML Comparison

When implementing components, it's crucial to compare the actual HTML output between DatastarUI and shadcn/ui to ensure pixel-perfect accuracy.

### HTML Comparison Workflow with shad-diffs/

All HTML comparison files should be stored in `/d/cdev/shad-diffs` to keep the root directory clean.

**Step 1: Capture HTML Files**

```bash
# Navigate to the project root
cd /mnt/d/cdev

# Capture DatastarUI HTML (assuming dev server on port 4242)
curl -s http://localhost:4242/components/tabs > shad-diffs/datastarui-tabs.html

# Capture shadcn/ui HTML (assuming docs server on port 3333)
curl -s http://localhost:3333/docs/components/tabs > shad-diffs/shadcn-tabs.html

# Or capture from live sites if local servers aren't running
curl -s https://ui.shadcn.com/docs/components/tabs > shad-diffs/shadcn-tabs-live.html
```

**Step 2: Extract and Compare Specific Patterns**

```bash
# Compare grid classes between implementations
echo "=== DatastarUI Grid Classes ==="
cat shad-diffs/datastarui-tabs.html | rg 'grid-cols-' -o
echo "=== shadcn/ui Grid Classes ==="
cat shad-diffs/shadcn-tabs.html | rg 'grid-cols-' -o

# Compare all class attributes
echo "=== DatastarUI Classes ==="
cat shad-diffs/datastarui-tabs.html | rg 'class="[^"]*"' -o
echo "=== shadcn/ui Classes ==="
cat shad-diffs/shadcn-tabs.html | rg 'class="[^"]*"' -o
```

**Step 3: Extract Specific Patterns with ripgrep**

```bash
# Find data-slot attributes with context
cat shad-diffs/datastarui-tabs.html | rg 'data-slot=' -A 1 -B 1

# Extract specific CSS classes (e.g., button variants)
cat shad-diffs/datastarui-button.html | rg 'class="[^"]*button[^"]*"' -o

# Find Tailwind spacing classes
cat shad-diffs/datastarui-tabs.html | rg 'p-[0-9]|m-[0-9]|px-[0-9]|py-[0-9]' -o

# Find ARIA attributes
cat shad-diffs/datastarui-tabs.html | rg 'aria-[a-z-]+=' -o

# Extract data attributes (useful for Datastar bindings)
cat shad-diffs/datastarui-tabs.html | rg 'data-[a-z-]+=' -o
```

**Step 4: Compare Component Structures**

```bash
# Find component wrapper elements in DatastarUI
cat shad-diffs/datastarui-tabs.html | rg '<div[^>]*data-' -A 2 -B 1

# Find component wrapper elements in shadcn/ui
cat shad-diffs/shadcn-tabs.html | rg '<div[^>]*class="[^"]*tabs[^"]*"' -A 2 -B 1

# Extract button elements with context
cat shad-diffs/datastarui-button.html | rg '<button' -A 3 -B 1
```

**Step 5: Advanced Pattern Analysis**

```bash
# Find all Tailwind color classes
cat shad-diffs/datastarui-button.html | rg 'bg-[a-z-]+|text-[a-z-]+|border-[a-z-]+' -o

# Extract hover and focus states
cat shad-diffs/datastarui-button.html | rg 'hover:|focus:' -o

# Find responsive classes
cat shad-diffs/datastarui-tabs.html | rg 'sm:|md:|lg:|xl:' -o

# Compare layout classes (flex, grid, etc.)
echo "=== DatastarUI Layout Classes ==="
cat shad-diffs/datastarui-tabs.html | rg 'flex|grid|inline' -o
echo "=== shadcn/ui Layout Classes ==="
cat shad-diffs/shadcn-tabs.html | rg 'flex|grid|inline' -o
```

**Step 6: Side-by-Side Component Analysis**

```bash
# Create a comparison script for easy analysis
cat > shad-diffs/compare-tabs.sh << 'EOF'
#!/bin/bash
echo "=== TABS LAYOUT COMPARISON ==="
echo "DatastarUI tabs container:"
cat shad-diffs/datastarui-tabs.html | rg 'data-slot="tabs"' -A 5 -B 1
echo -e "\nshadcn/ui tabs container:"
cat shad-diffs/shadcn-tabs.html | rg 'class="[^"]*tabs[^"]*"' -A 5 -B 1

echo -e "\n=== TABS LIST COMPARISON ==="
echo "DatastarUI tabs list:"
cat shad-diffs/datastarui-tabs.html | rg 'role="tablist"' -A 3 -B 1
echo -e "\nshadcn/ui tabs list:"
cat shad-diffs/shadcn-tabs.html | rg 'role="tablist"' -A 3 -B 1
EOF

chmod +x shad-diffs/compare-tabs.sh
./shad-diffs/compare-tabs.sh
```

### Debugging Checklist

1. **HTML Structure**: Compare the DOM structure between implementations
2. **CSS Classes**: Ensure exact class names match shadcn/ui
3. **Data Attributes**: Verify Datastar bindings are working correctly. **IMPORTANT**: shadcn does not use datastar. we will need to reason about how to use datastar to simplify and achieve the same functionallity as the react component.
4. **Interactive States**: Test all interactive states (onclick, hover, active, disabled)
5. **Accessibility**: Check ARIA attributes and roles
6. **Responsive Behavior**: Test on different screen sizes

# IMPORTANT

The developer is running a live reload server wathcing for file changes. Do not try to run the compiled binary as it is already running. Templ files will be automatically generated, but feel free to run `templ generate` to check for errors.
