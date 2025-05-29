# DatastarUI Component Development Guide

## Project Overview

DatastarUI is a Go/templ port of shadcn/ui components that maintains pixel-perfect visual and behavioral parity while eliminating JavaScript dependencies (except for the 15KB Datastar library for reactivity).

### Key Features

- **Pixel-perfect shadcn/ui components** ported to Go/templ
- **Datastar signals** for interactive components

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
│   │   ├── buttonpage/
│   │   │   └── button_page.templ # Button demo page
│   │   ├── tabspage/
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

## IMPORTANT

The developer is running a live reload server wathcing for file changes. Do not try to run the compiled binary as it is already running. Templ files will be automatically generated, but feel free to run `templ generate` to check for errors.

The root directory for datastarui is `/d/cdev/datastarui`.

### Tailwind CSS Watch Process

**IMPORTANT**: Tailwind CSS is running in watch mode during development. The developer has `tailwindcss --watch` running automatically, which means:

- CSS changes in `static/css/index.css` are automatically compiled to `static/css/build.css`
- The watch process monitors all `.templ` files for class changes
- Only run `tailwind` or similar commands if there are CSS compilation issues that need debugging

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

# Go/templ Patterns

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

**Namespaced Signals**: Only leaf nodes are signals:

```go
// ✅ Correct - both are valid signals
data-signals-user.name="'John'" data-signals-count="0"
data-text="$user.name"  // Valid
data-text="$count"      // Valid

// ❌ Wrong - namespace is not a signal
data-text="$user"       // Invalid if user.name exists
```

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

#### Datastar Fetching Signal

The data-indicator attribute accepts the name of a signal whose value is set to true when a fetch request initiated from the same element is in progress, otherwise false. If the signal does not exist in the signals, it will be added.

Note: If you use the data-indicator attribute, you MUST also make sure to have a unique id attribute on the element that is making the fetch request. The is because the element might not exist otherwise nor be stable when the fetch request is completed.

```html
<style>
  .indicator {
    opacity: 0;
    transition: opacity 300ms ease-out;
  }
  .indicator.loading {
    opacity: 1;
    transition: opacity 300ms ease-in;
  }
</style>
<button
  id="greetingBtn"
  data-indicator-fetching
  data-on-click="@get('/examples/fetch_indicator/greet')"
  data-attr-disabled="$fetching"
>
  Click me for a greeting
</button>
<div class="indicator" data-class="{loading: $fetching}">Loading Indicator</div>
<div id="greeting"></div>
```

### Datastar Forms

Setting the contentType option to form tells the @get() action to look for the closest form, perform validation on it, and send all form elements within it to the backend. A selector option can be provided to specify a form element. No signals are sent to the backend in this type of request.

```html
<form>
  <input name="foo" required placeholder="Type foo contents" />
  <button data-on-click="@get('/endpoint', {contentType: 'form'})">
    Submit GET request
  </button>
  <button data-on-click="@post('/endpoint', {contentType: 'form'})">
    Submit POST request
  </button>
</form>

<button
  data-on-click="@get('/endpoint', {contentType: 'form', selector: '#myform'})"
>
  Submit GET request from outside the form
</button>
```

#### DatastarUI Forms

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

### Datastar Expressions

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
