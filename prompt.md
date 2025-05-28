# DatastarUI Component Development Guide

## Project Overview

DatastarUI is a Go/templ port of shadcn/ui components that maintains pixel-perfect visual and behavioral parity while eliminating JavaScript dependencies (except for the 15KB Datastar library for reactivity).

### Key Features

- **Pixel-perfect shadcn/ui components** ported to Go/templ
- **Zero JavaScript dependencies** (except 15KB Datastar for reactivity)
- **Type-safe form handling** with protobuf integration
- **Folder-based component organization** with dedicated demo pages
- **Datastar reactive patterns** for interactive components
- **Server-side form parsing** into structured protobuf messages

### Project Structure

```
datastarui/
├── components/                    # Reusable UI components
│   ├── button/
│   │   ├── button.templ          # Component template
│   │   ├── types.go              # Props and types
│   │   └── variants.go           # CSS variants
│   ├── form/                     # Form components
│   ├── input/                    # Input component
│   └── ...                       # Other components
├── pages/
│   ├── components/               # Component demo pages
│   │   ├── button/
│   │   │   └── button_page.templ # Button demo page
│   │   ├── tabs/
│   │   │   ├── tabs_page.templ   # Tabs demo page
│   │   │   ├── tabs_page.proto   # Form definitions
│   │   │   └── tabs_page.pb.go   # Generated protobuf code
│   │   └── ...                   # Other demo pages
│   ├── home_page.templ           # Home page
│   └── docs_page.templ           # Documentation page
├── layouts/                      # Page layouts and navigation
├── static/                       # Static assets (CSS, JS)
└── main.go                       # Server and routing
```

## Component Architecture

### File Structure

Each component follows this structure:

```
components/[component-name]/
├── [component-name].templ     # Main component template
├── types.go                   # Component props and types
└── variants.go                # CSS class generation (CVA-style)
```

### Demo Page Structure

Each component demo page follows this structure:

```
pages/components/[component-name]/
├── [component-name]_page.templ    # Demo page template
└── [component-name]_page.proto    # Protobuf definitions for forms
```

### Protobuf Integration

Components that handle forms include protobuf definitions for type-safe form parsing:

```
pages/components/tabs/
├── tabs_page.templ     # Demo page with form examples
└── tabs_page.proto     # Form message definitions
```

Example protobuf definition:

```protobuf
syntax = "proto3";

package pages.components;

option go_package = "github.com/coreycole/datastarui/pages/components";

// Account form data
message AccountForm {
  string name = 1;
  string username = 2;
}

// Password form data
message PasswordForm {
  string current_password = 1;
  string new_password = 2;
}

// Response messages
message FormResponse {
  bool success = 1;
  string message = 2;
  repeated string errors = 3;
}
```

### Form Component Architecture

Form components follow a specific pattern that integrates Datastar reactivity with protobuf type safety:

#### 1. Form Structure

```go
// Form components use the shadcn/ui form pattern
@form.Form(form.FormProps{ID: "account-form", Method: "POST", Action: "/api/account"}) {
    @form.FormItem(form.FormItemProps{}) {
        @form.FormLabel(form.FormLabelProps{For: "username"}) {
            Username
        }
        @form.FormControl(form.FormControlProps{ID: "username"}) {
            @input.Input(input.InputProps{
                ID:          "username",
                Name:        "username",
                Placeholder: "Enter username",
            })
        }
        @form.FormDescription(form.FormDescriptionProps{}) {
            This is your public display name.
        }
        @form.FormMessage(form.FormMessageProps{}) {
            // Error messages appear here
        }
    }
}
```

#### 2. Datastar Integration

```go
// Use Datastar signals for reactive form state
<div data-signals="{formData: {username: ''}, errors: {}, isSubmitting: false}">
    @input.Input(input.InputProps{
        Attributes: templ.Attributes{
            "data-model": "$formData.username",
            "data-class": "$errors.username ? 'border-destructive' : ''",
        },
    })
</div>
```

#### 3. Server-Side Form Parsing

```go
// Example HTTP handler using the automatic form decoder
func handleAccountFormSubmission(c echo.Context) error {
    // Use the automatic decoder with protobuf reflection
    form, response, err := tabs.DecodeAccountForm(c.Request())
    if err != nil {
        return c.JSON(500, map[string]string{"error": err.Error()})
    }

    // Check if validation failed
    if !response.Success {
        return c.JSON(400, response)
    }

    // Process the valid form data
    log.Printf("Received account form: Name=%s, Username=%s", form.Name, form.Username)

    // Save to database, etc.
    // ...

    return c.JSON(200, &tabs.FormResponse{
        Success: true,
        Message: "Account updated successfully",
    })
}
```

#### 4. Complete Form Handler Example

Using the automatic form decoder with protobuf reflection:

```go
// Example HTTP handler using the form decoder
func handleAccountFormSubmission(c echo.Context) error {
    // Use the automatic decoder
    form, response, err := tabs.DecodeAccountForm(c.Request())
    if err != nil {
        return c.JSON(500, map[string]string{"error": err.Error()})
    }

    // Check if validation failed
    if !response.Success {
        return c.JSON(400, response)
    }

    // Process the valid form data
    log.Printf("Received account form: Name=%s, Username=%s", form.Name, form.Username)

    // Save to database, etc.
    // ...

    return c.JSON(200, &tabs.FormResponse{
        Success: true,
        Message: "Account updated successfully",
    })
}

// Example route setup in main.go
func setupFormRoutes(e *echo.Echo) {
    e.POST("/api/account", handleAccountFormSubmission)
    e.POST("/api/password", handlePasswordFormSubmission)
}

// Password form handler
func handlePasswordFormSubmission(c echo.Context) error {
    form, response, err := tabs.DecodePasswordForm(c.Request())
    if err != nil {
        return c.JSON(500, map[string]string{"error": err.Error()})
    }

    if !response.Success {
        return c.JSON(400, response)
    }

    // Validate current password, update to new password
    log.Printf("Password change request for user")

    return c.JSON(200, &tabs.FormResponse{
        Success: true,
        Message: "Password updated successfully",
    })
}
```

#### 5. Form Decoder Features

The automatic form decoder provides:

- **Type-safe parsing**: Uses protobuf reflection to automatically parse form fields
- **Automatic validation**: Built-in validation for common field types (strings, numbers, booleans)
- **Custom validation rules**: Field-specific validation (username format, password length, etc.)
- **Error handling**: Detailed error messages for validation failures
- **Support for all protobuf types**: strings, integers, floats, booleans, enums, bytes

**Supported Field Types:**

- `string` - Text fields with length and format validation
- `int32`, `int64` - Integer fields with range validation
- `uint32`, `uint64` - Unsigned integer fields
- `float`, `double` - Floating point numbers
- `bool` - Boolean fields (accepts "true"/"false", "1"/"0")
- `bytes` - Binary data fields
- `enum` - Enumeration fields (by name or number)

**Built-in Validation Rules:**

- **Name fields**: 2-50 characters
- **Username fields**: 3-20 characters, alphanumeric + underscore/hyphen only
- **Password fields**: 8-128 characters minimum

#### 6. Datastar Form Submission Patterns

DatastarUI uses a consistent pattern for form submissions with Datastar:

**Route Pattern**: `POST /forms/{form-name}`

**Frontend Pattern**:

```go
// Datastar signals for form state
<div data-signals="{
    formData: {field1: '', field2: ''},
    submitting: false,
    response: null
}">
    <form data-on-submit="
        $submitting = true;
        $response = null;
        fetch('/forms/form-name', {
            method: 'POST',
            headers: {'Content-Type': 'application/x-www-form-urlencoded'},
            body: new URLSearchParams({
                field1: $formData.field1,
                field2: $formData.field2
            })
        })
        .then(r => r.json())
        .then(data => {
            $response = data;
            $submitting = false;
            if (data.success) {
                // Clear form on success
                $formData.field1 = '';
                $formData.field2 = '';
            }
        })
        .catch(err => {
            $response = {success: false, message: 'Network error', errors: [err.message]};
            $submitting = false;
        });
        event.preventDefault();
    ">
        <!-- Form fields with data-model bindings -->
        @input.Input(input.InputProps{
            Attributes: templ.Attributes{
                "data-model": "$formData.field1",
            },
        })

        <!-- Error display -->
        <div data-show="$response && !$response.success" class="text-sm text-destructive">
            <span data-text="$response.message"></span>
        </div>

        <!-- Success display -->
        <div data-show="$response && $response.success" class="text-sm text-green-600">
            <span data-text="$response.message"></span>
        </div>

        <!-- Submit button with loading state -->
        @button.Button(button.ButtonProps{
            Type: "submit",
            Attributes: templ.Attributes{
                "data-disabled": "$submitting",
            },
        }) {
            <span data-show="!$submitting">Submit</span>
            <span data-show="$submitting">Submitting...</span>
        }
    </form>
</div>
```

**Backend Pattern**:

```go
// Route handler using protobuf decoder
e.POST("/forms/account", func(c echo.Context) error {
    form, response, err := tabs.DecodeAccountForm(c.Request())
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    if !response.Success {
        return c.JSON(http.StatusBadRequest, response)
    }

    // Process the validated form data
    log.Printf("Form submitted: %+v", form)

    return c.JSON(http.StatusOK, &tabs.FormResponse{
        Success: true,
        Message: "Form submitted successfully!",
        Errors:  []string{},
    })
})
```

**Key Features**:

- **Reactive form state**: Real-time updates with Datastar signals
- **Loading states**: Visual feedback during submission
- **Error handling**: Field-specific and general error display
- **Success feedback**: Clear success messages
- **Form clearing**: Automatic form reset on successful submission
- **Type safety**: Protobuf validation and parsing

### 3. Datastar Integration

- **Use Datastar signals** for reactive state: `data-signals-count="0"`
- **Handle events** with Datastar: `data-on-click="$count++"`
- **Bind reactive content**: `data-text="$count"`, `data-show="$visible"`

#### Datastar Expressions

Datastar expressions are JavaScript-like strings evaluated by Datastar attributes. They provide powerful declarative reactivity with some key differences from standard JavaScript.

**Signal Access**: Signals use the `$` prefix and are automatically converted to `ctx.signals.signal('name').value`:

```go
// Basic signal usage
data-text="$count"                   // Displays signal value
data-show="$isVisible"               // Conditional visibility
data-text="$user.name"               // Nested signal access
data-text="$items.length"            // JavaScript properties work

// Signal assignment
data-on-click="$count++"             // Increment signal
data-on-click="$user.name = 'John'"  // Set signal value
```

**Actions**: Actions use the `@` prefix for

```go
data-on-click="@post('/api/endpoint')"           // POST request
data-on-click="$count++; @post('/api/count')"   // Multiple statements
```

**JavaScript Operators**: Standard JavaScript operators are available:

```go
// Logical operators
data-on-click="$isReady && @post('/launch')"
data-show="$user && $user.isActive"

// Ternary operator - works with most attributes
data-attr-class="$theme === 'dark' ? 'bg-black text-white' : 'bg-white text-black'"
data-text="$count > 0 ? $count + ' items' : 'No items'"
data-attr-disabled="$loading ? 'true' : 'false'"
data-show="$user ? true : false"  // Though just $user is simpler

// Nested ternary (use sparingly)
data-text="$status === 'loading' ? 'Loading...' : $status === 'error' ? 'Error occurred' : 'Success'"

// Ternary with complex expressions
data-attr-href="$user.isAdmin ? '/admin/dashboard' : '/user/profile'"
data-text="$items.length === 0 ? 'No items found' : $items.length === 1 ? '1 item' : $items.length + ' items'"

// Comparison operators
data-show="$status === 'success'"
data-show="$count >= 10"
```

**Important Note**: The ternary operator works with most Datastar attributes (`data-text`, `data-show`, `data-attr-*`, etc.) but **NOT** with `data-class`, which requires object syntax.

**Multiple Statements**: Use semicolons to separate statements:

```go
// Single line
data-on-click="$count++; $message = 'Updated'; @post('/api/update')"

// Multi-line (semicolons required)
data-on-submit="
    $submitting = true;
    $errors = {};
    @post('/api/form')
"
```

**Event Context**: Access browser events with `evt`:

```go
data-on-input="$value = evt.target.value"
data-on-keydown="evt.key === 'Enter' && @post('/search')"
data-on-click="evt.preventDefault(); $modal = true"
```

**Namespaced Signals**: Only leaf nodes are signals:

```go
// ✅ Correct - both are valid signals
data-signals-user.name="'John'" data-signals-count="0"
data-text="$user.name"  // Valid
data-text="$count"      // Valid

// ❌ Wrong - namespace is not a signal
data-text="$user"       // Invalid if user.name exists
```

**Key Rules**:

- **Semicolons required**: Line breaks alone don't separate statements
- **Implicit return**: Last statement is automatically returned
- **No global $**: JavaScript globals starting with `$` are not accessible
- **Signal-first**: `$signalName` is converted before JavaScript evaluation

**Common Patterns**:

```go
// Conditional actions
data-on-click="$isValid && @post('/submit')"

// State updates with feedback
data-on-click="$loading = true; @post('/api').then(() => $loading = false)"

// Form validation
data-on-blur="$touched.email = true; $errors.email = $email ? '' : 'Required'"

// Toggle patterns
data-on-click="$isOpen = !$isOpen"

// Complex conditions
data-show="$user && $user.role === 'admin' && !$loading"
```

#### Consistent ID Patterns and Signal Naming

**Critical for Component Interactivity**: Datastar signals are global on the page, so consistent ID patterns and signal naming are essential for components to work correctly. We leverage Datastar's namespacing capabilities to create organized, collision-free signal structures.

#### Datastar Namespacing

Datastar supports namespaced signals using dot notation. Only leaf nodes are actual signals:

```go
// ✅ Correct - bar is the signal, foo is the namespace
data-signals-foo.bar="1"        // Creates signal $foo.bar
data-text="$foo.bar"            // Valid - accessing the signal

// ❌ Wrong - foo is not a signal when foo.bar exists
data-text="$foo"                // Invalid - foo is just a namespace
```

#### Recommended ID and Signal Patterns

**ID Format**: Use hyphens for readability, avoid dots in IDs themselves:

```go
// ✅ Good ID patterns
"basic-checkbox"
"user-profile-form"
"navigation-menu"
"product-card-123"
```

**DatastarUI Component Signal Pattern**: Use the component ID as the namespace for all component signals:

```go
// Standard pattern: {[componentID]: {signal1: value, signal2: value}}
// Examples:
$terms.checked              // Checkbox component with ID "terms"
$user-menu.open             // Dropdown component with ID "user-menu"
$contact-form.submitting    // Form component with ID "contact-form"
$product-123.selected       // Product card with ID "product-123"
```

**✅ DatastarUI Component Pattern - ID as Namespace:**

```go
templ Checkbox(props CheckboxProps) {
    {{
        // Create nested signal structure: {[id]: {checked: false}}
        // This allows the ID to namespace all signals for this component instance
        initialValue := "false"
        if props.Checked {
            initialValue = "true"
        }
        signals := "{\"" + props.ID + "\": {\"checked\": " + initialValue + "}}"
        signalRef := "$" + props.ID + ".checked"
        toggleExpr := signalRef + " = !" + signalRef
    }}
    <div data-signals={ signals }>
        <button
            id={ props.ID }
            data-on-click={ toggleExpr }
            data-attr-aria-checked={ signalRef + " ? 'true' : 'false'" }
            data-show={ signalRef }
        >
            <!-- component content -->
        </button>
    </div>
}

templ CheckboxExampleUsage() {
    <div>
        @checkbox.Checkbox(checkbox.CheckboxProps{
            ID: "terms",  // Creates signals: $terms.checked
        })
        <!-- Access signal from anywhere on page -->
        <span data-text="$terms.checked ? 'Accepted' : 'Not accepted'"></span>
    </div>
}
```

**Multi-Signal Component Example:**

```go
templ Dropdown(props DropdownProps) {
    {{
        // Multiple signals namespaced under component ID
        signals := "{\"" + props.ID + "\": {\"open\": false, \"selected\": null, \"loading\": false}}"
    }}
    <div data-signals={ signals }>
        <button data-on-click="$" + props.ID + ".open = !$" + props.ID + ".open">
            Toggle Dropdown
        </button>
        <div data-show="$" + props.ID + ".open">
            <!-- dropdown content -->
        </div>
        <!-- Loading state -->
        <div data-show="$" + props.ID + ".loading">Loading...</div>
        <!-- Selected value -->
        <span data-text="$" + props.ID + ".selected || 'None selected'"></span>
    </div>
}

templ DropdownExampleUsage() {
    <div>
        @dropdown.Dropdown(dropdown.DropdownProps{
            ID: "country-selector",  // Creates: $country-selector.open, $country-selector.selected, $country-selector.loading
        })

        <!-- Access any signal from anywhere on page -->
        <div data-show="$country-selector.open">Dropdown is open!</div>
        <div data-text="'Selected: ' + ($country-selector.selected || 'None')"></div>
    </div>
}
```

**Benefits of ID-as-Namespace Pattern:**

1. **Collision Prevention**: Each component instance has its own signal namespace
2. **Scalability**: Easy to add more signals to a component without conflicts
3. **Predictable Structure**: `$componentID.signalName` is always the pattern
4. **Global Access**: Any component's state can be accessed from anywhere on the page
5. **Clean Organization**: Related signals are grouped under the component ID
6. **Future-Proof**: Works for simple components (1 signal) and complex ones (many signals)

**Page-Level Signal Organization**: For coordinating multiple components, use logical page-level namespaces:

```go
// Page-level coordination
data-signals="{
    forms: {
        contact: {name: '', email: '', terms: false},
        newsletter: {email: '', frequency: 'weekly'}
    },
    ui: {
        modal: {visible: false},
        loading: {submit: false}
    }
}"

// Component signals remain separate and don't conflict
// $contact-form.submitting (component)
// $forms.contact.terms (page-level)
```

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

#### Common templ Gotchas and Solutions

**URL Safety**: Always wrap href attributes with `templ.SafeURL()` for security:

```go
// ❌ Wrong - will cause compilation error
<a href={ props.Href }>

// ✅ Correct - prevents XSS attacks
<a href={ templ.SafeURL(props.Href) }>
```

**Template Structure**: All content must be inside the layout wrapper:

```go
// ❌ Wrong - content outside @l.Root() will cause "expected nodes" error
templ MyPage() {
    <div>Outside content</div>
    @l.Root("page") {
        <div>Inside content</div>
    }
}

// ✅ Correct - all content inside @l.Root()
templ MyPage() {
    @l.Root("page") {
        <div>All content here</div>
    }
}
```

**Component Composition**: Use proper nesting for multi-part components:

```go
// ✅ Correct breadcrumb structure
templ ComponentPageBreadcrumbs(currentPage string) {
	<div>
		@breadcrumb.Breadcrumb(breadcrumb.BreadcrumbProps{}) {
			@breadcrumb.BreadcrumbList(breadcrumb.BreadcrumbListProps{}) {
				@breadcrumb.BreadcrumbItem(breadcrumb.BreadcrumbItemProps{}) {
					@breadcrumb.BreadcrumbLink(breadcrumb.BreadcrumbLinkProps{Href: "/docs"}) {
						Docs
					}
				}
				@breadcrumb.BreadcrumbSeparator(breadcrumb.BreadcrumbSeparatorProps{})
				@breadcrumb.BreadcrumbItem(breadcrumb.BreadcrumbItemProps{}) {
					@breadcrumb.BreadcrumbPage(breadcrumb.BreadcrumbPageProps{}) {
						{ currentPage }
					}
				}
			}
		}
	</div>
}
```

### Tailwind CSS Watch Process

**IMPORTANT**: Tailwind CSS is running in watch mode during development. The developer has `tailwindcss --watch` running automatically, which means:

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

### Step 4: Create Component Demo Page

1. **Create demo page directory**: `pages/components/[component-name]/`
2. **Create demo page**: `pages/components/[component-name]/[component-name]_page.templ`
3. **Follow the button component pattern** for consistency
4. **Include comprehensive examples** showcasing all variants and features
5. **Add descriptive sections** explaining each example
6. **Add route to main.go** and **sidebar.go** to link to the new component demo page

### Step 6: Update Routing Structure

Update the main.go file to use the new folder-based imports:

```go
import (
    "github.com/coreycole/datastarui/pages/components/button"
    "github.com/coreycole/datastarui/pages/components/tabs"
    // ... other component imports
)

// Route handlers
e.GET("/components/button", func(c echo.Context) error {
    component := button.ButtonPage()
    return component.Render(c.Request().Context(), c.Response().Writer)
})
```

## Debugging and HTML Comparison

When implementing components, it's crucial to compare the actual HTML output between DatastarUI and shadcn/ui to ensure pixel-perfect accuracy.

### HTML Comparison Workflow with shad-diffs/

All HTML comparison files should be stored in `/mnt/d/cdev/shad-diffs/` to keep the root directory clean.

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

## The Datastar Way - Best Practices

### Stop Overcomplicating It

Most of the time, if you run into issues when using Datastar, you are probably **overcomplicating it™**.

Datastar is a **hypermedia framework**, not a JavaScript framework. If you approach it like a JavaScript framework, you are likely to run into complications.

### The Hypermedia Approach

Between attribute plugins and action plugins, Datastar provides everything you need to build hypermedia-driven applications. Using this approach:

- **Backend drives state** to the frontend and acts as the single source of truth
- **Server determines** what actions the user can take next
- **State flows down through props, events flow up** - always encapsulate state and send props down, events up

### When You Need Additional JavaScript

Any additional JavaScript functionality that doesn't work via `data-*` attributes and `datastar-execute-script` SSE events should be extracted into:

1. **External scripts** (preferred for simple functions)
2. **Web components** (preferred for reusable, encapsulated elements)

### External Scripts Pattern

When using external scripts, pass data via arguments and return results or dispatch custom events:

```go
// Synchronous function call
<div data-signals-result="''">
    <input data-bind-foo
           data-on-input="$result = myfunction($foo)"
    >
    <span data-text="$result"></span>
</div>
```

```javascript
function myfunction(data) {
  return `You entered ${data}`;
}
```

**Asynchronous functions** must dispatch custom events (Datastar won't await them):

```go
// Asynchronous function with custom event
<div data-signals-result="''"
     data-on-mycustomevent__window="$result = evt.detail.value"
>
    <input data-bind-foo
           data-on-input="myfunction($foo)"
    >
    <span data-text="$result"></span>
</div>
```

```javascript
async function myfunction(data) {
  const value = await new Promise((resolve) => {
    setTimeout(() => resolve(`You entered ${data}`), 1000);
  });
  window.dispatchEvent(new CustomEvent("mycustomevent", { detail: { value } }));
}
```

### Web Components Pattern

Web components create reusable, encapsulated custom elements with no external dependencies:

```go
// Using web component with attribute binding
<div data-signals-result="''">
    <input data-bind-foo />
    <my-component
        data-attr-src="$foo"
        data-on-mycustomevent="$result = evt.detail.value"
    ></my-component>
    <span data-text="$result"></span>
</div>
```

```javascript
class MyComponent extends HTMLElement {
  static get observedAttributes() {
    return ["src"];
  }

  attributeChangedCallback(name, oldValue, newValue) {
    const value = `You entered ${newValue}`;
    this.dispatchEvent(new CustomEvent("mycustomevent", { detail: { value } }));
  }
}

customElements.define("my-component", MyComponent);
```

**Alternative: Using `data-bind` with web components:**

```go
// Binding directly to web component value
<input data-bind-foo />
<my-component
    data-attr-src="$foo"
    data-bind-result
></my-component>
<span data-text="$result"></span>
```

```javascript
class MyComponent extends HTMLElement {
  static get observedAttributes() {
    return ["src"];
  }

  attributeChangedCallback(name, oldValue, newValue) {
    this.value = `You entered ${newValue}`;
    this.dispatchEvent(new Event("change")); // Required for data-bind
  }
}

customElements.define("my-component", MyComponent);
```

### Key Principles

1. **Think hypermedia first** - let the server drive state and available actions
2. **Use data-\* attributes** for all reactive behavior when possible
3. **Extract complex logic** into external scripts or web components
4. **Follow props down, events up** - encapsulate functionality and communicate via well-defined interfaces
5. **Avoid JavaScript framework patterns** - resist the urge to manage state in the frontend
6. **Keep it simple** - if it feels complicated, you're probably overengineering it

### DatastarUI Component Guidelines

When building DatastarUI components:

- **Prefer server-driven state** over complex client-side logic
- **Use Datastar expressions** for simple reactive behavior
- **Extract complex interactions** into web components when needed
- **Follow the hypermedia mindset** - components should be declarative and server-controlled
- **Test with minimal JavaScript** - the goal is to eliminate JavaScript dependencies while maintaining full functionality
