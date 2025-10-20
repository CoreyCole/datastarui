---
date: 2025-10-19T00:00:00Z
researcher: Claude
git_commit: 9772de32de0b5f5dbe7d16f7c085b141f3edb19f
branch: main
repository: datastarui
topic: "Adding flexible error messages per field for proto validation in login form"
tags: [research, codebase, proto-validation, form-errors, datastar, connect-rpc, http-wrapper]
status: complete
last_updated: 2025-10-19
last_updated_by: Claude
last_updated_note: "Added follow-up research on HTTP wrapper solution for returning HTML from Connect RPC handlers"
---

# Research: Adding Flexible Error Messages Per Field for Proto Validation

**Date**: 2025-10-19T00:00:00Z
**Researcher**: Claude
**Git Commit**: 9772de32de0b5f5dbe7d16f7c085b141f3edb19f
**Branch**: main
**Repository**: datastarui

## Research Question

How to add flexible error messages per field on the frontend for the recently added login form with proto validation? The requirement is to:
1. Get structured error responses from proto validate that can be parsed
2. Slot error messages under the correct form field
3. Re-render the form with user content and error messages so Datastar can use morphing to merge in just the differences

## Summary

The codebase has a well-structured proto validation system using buf.validate and Connect RPC that generates detailed field-level validation errors in the `details.debug.violations` array. However, there is currently **no implementation for parsing these violations or displaying field-specific errors** in the UI.

The good news is that:
- Proto validation already generates structured violations with field names and messages
- Form components have built-in support for displaying field errors (`FormMessage`, `FormLabel.HasError`)
- Datastar's signal system with `data-bind` preserves form input values without needing morphing
- The error data structure is well-documented and contains all needed information

**Key recommendation**: Implement a server-side violation parser that extracts field errors from the violations array and returns them as part of a re-rendered form fragment with preserved user input values.

## Detailed Findings

### Current Proto Validation Implementation

The login form uses a three-layer architecture:

1. **Proto Definition with Validation Rules** (`proto/com/datastarui/v1/forms/login/login.proto:10-20`)
   - Email field: `required: true`, `string: {email: true}`
   - Password field: `required: true`, `string: {min_len: 8}`

2. **Validation Interceptor** (`api/interceptors/validation.go:15-40`)
   ```go
   if validationErr := new(protovalidate.ValidationError); errors.As(err, &validationErr) {
       if detail, err := connect.NewErrorDetail(validationErr.ToProto()); err == nil {
           connectErr.AddDetail(detail)
       }
   }
   ```
   - Validates all incoming RPC requests
   - Attaches detailed violations to error responses
   - Violations are in `details` array with type `buf.validate.Violations`

3. **Error Response Structure** (from provided example)
   ```json
   {
     "code": "invalid_argument",
     "message": "validation error:\n - password: value length must be at least 8 characters [string.min_len]",
     "details": [{
       "type": "buf.validate.Violations",
       "value": "Cm0SDnN0cmluZy5taW5fbGVuGip2YWx1ZSBsZW5ndGggbXVzdCBiZSBhdCBsZWBzdCA4IGNoYXJhY3RlcnMqEAoOCAISCHBhc3N3b3JkGAkyHQoMCA4SBnN0cmluZxgLCg0IAhIHbWluX2xlbhgE",
       "debug": {
         "violations": [
           {
             "field": {
               "elements": [
                 {"fieldNumber": 2, "fieldName": "password", "fieldType": "TYPE_STRING"}
               ]
             },
             "rule": {
               "elements": [
                 {"fieldNumber": 14, "fieldName": "string", "fieldType": "TYPE_MESSAGE"},
                 {"fieldNumber": 2, "fieldName": "min_len", "fieldType": "TYPE_UINT64"}
               ]
             },
             "ruleId": "string.min_len",
             "message": "value length must be at least 8 characters"
           }
         ]
       }
     }]
   }
   ```

### Current Gaps

1. **No Violation Parsing**: No utilities exist to extract field names and messages from the violations array
2. **No Field Error Display**: Forms only show general errors, not field-specific ones (`forms/login/login_form.templ:36-40`)
3. **No Server-Side Re-Rendering**: No pattern for returning form HTML with errors after validation failure
4. **No Field-to-Signal Mapping**: No mechanism to map proto field names to Datastar signals

### Existing Form Error Components (Ready to Use)

The codebase already has components for displaying field-specific errors:

1. **FormMessage Component** (`components/form/form.templ:149-170`)
   ```go
   @form.FormMessage(form.FormMessageArgs{
       Message: "This field is required",  // Can be populated with violation message
   })
   ```

2. **FormLabel with Error State** (`components/form/form.templ:116-130`)
   ```go
   @form.FormLabel(form.FormLabelArgs{
       For:      "email",
       HasError: true,  // Changes label to destructive color
   })
   ```

3. **Input with ARIA Invalid** (`components/input/variants.go:5-15`)
   - Built-in styling for `aria-invalid="true"`
   - Shows destructive border and ring

### How Datastar Preserves Form Values

**Important discovery**: Datastar doesn't use morphing for form value preservation. Instead:

1. **Two-Way Binding** (`components/input/input.templ:45-54`)
   ```go
   inputAttrs["data-bind"] = "login.email"  // Auto-generated from FormID and Name
   ```
   - Creates reactive binding between input and signal
   - Input value always reflects signal value
   - Signal persists across DOM updates

2. **Signal-Based State** (`forms/login/login_form.templ:29-32`)
   ```go
   signals := form.SignalsWithFormId(formID, LoginFormSignals{
       Email:    "",
       Password: "",
   })
   ```

3. **No Morphing Needed**
   - Form values live in Datastar signals (client-side)
   - `data-bind` maintains values automatically
   - Server can return new HTML, values preserved by signals

## Code References

### Proto Validation
- `api/interceptors/validation.go:15-40` - Validation interceptor that attaches violations
- `proto/com/datastarui/v1/forms/login/login.proto:10-20` - Login validation rules
- `api/services/auth/service.go:69-101` - Login handler

### Form Components
- `components/form/form.templ:149-170` - FormMessage component for field errors
- `components/form/form.templ:116-130` - FormLabel with error state support
- `components/form/variants.go:53-63` - Error message styling
- `components/input/input.templ:45-54` - Auto data-bind generation

### Current Login Form
- `forms/login/login_form.templ:21-92` - Login form implementation
- `forms/login/login_form.templ:36-40` - Current error display (global only)
- `forms/login/login_form.templ:29-32` - Signal initialization

### Data Flow
- `api/middleware/unwrap_form.go:44-67` - Form data unwrapping
- `components/form/form.templ:64-76` - Form submission with signal filtering
- `components/form/form.templ:11-33` - SignalsWithFormId helper

## Architecture Insights

### Pattern 1: Server-Side Form Re-Rendering with Errors

Based on the research, here's the recommended approach:

1. **Parse Violations on Server**
   ```go
   type FieldErrors map[string]string  // fieldName -> error message

   func ParseViolations(err error) FieldErrors {
       // Extract violations from Connect error details
       // Map field.elements[0].fieldName to message
   }
   ```

2. **Re-Render Form with Errors**
   ```go
   func (s *Service) Login(ctx context.Context, req *connect.Request[loginv1.LoginRequest]) (*connect.Response[loginv1.LoginResponse], error) {
       // Validate request...
       if err != nil {
           fieldErrors := ParseViolations(err)

           // Return HTML response with form re-rendered
           return RenderLoginFormWithErrors(req.Msg, fieldErrors)
       }
   }
   ```

3. **Preserve User Input in Re-Rendered Form**
   ```go
   signals := form.SignalsWithFormId(formID, LoginFormSignals{
       Email:    req.Email,     // Preserve user input
       Password: req.Password,  // Preserve user input
       EmailError: fieldErrors["email"],
       PasswordError: fieldErrors["password"],
   })
   ```

4. **Display Field Errors**
   ```html
   @form.FormItem() {
       @form.FormLabel(form.FormLabelArgs{
           For: "email",
           HasError: signals.Signal("emailError") != "",
       })
       @input.Input(...)
       @form.FormMessage(form.FormMessageArgs{
           Message: signals.Signal("emailError"),
       })
   }
   ```

### Pattern 2: Datastar Morphing Strategy

Since Datastar uses `data-bind` for form values, the morphing strategy is simplified:

1. **Server returns entire form HTML** with error messages
2. **Datastar merges the HTML** using its internal morphing algorithm
3. **Input values are preserved** via `data-bind` attributes (not morphing)
4. **Error messages appear** as new DOM elements

The key insight: **Morphing updates the DOM structure (adding error divs), while data-bind preserves input values**.

### Pattern 3: Connect RPC Error Detail Extraction

The violations are in the error details as base64-encoded protobuf:

```go
func ExtractViolations(connectErr *connect.Error) ([]*validate.Violation, error) {
    for _, detail := range connectErr.Details() {
        if detail.Type() == "buf.validate.Violations" {
            // Decode detail.Value() from base64
            // Unmarshal as Violations protobuf message
            // Return violations array
        }
    }
    return nil, nil
}
```

## Related Research

- `/mnt/d/cdev/datastarui/docs/proto-validate.md:501-525` - Complete error structure documentation
- No existing research documents found on form validation patterns

## Open Questions

1. **Response Format**: Should validation errors return HTML (for morphing) or JSON (for client-side rendering)?
   - HTML approach fits hypermedia pattern better
   - JSON requires client-side logic but is more flexible

2. **Error Signal Structure**: Should errors be separate signals or nested under fields?
   - Option A: `$login.emailError`, `$login.passwordError`
   - Option B: `$login.errors.email`, `$login.errors.password`

3. **Generic vs Specific Handlers**: Should each form have its own error handler or use a generic pattern?
   - Generic interceptor could handle all form validation errors
   - Specific handlers allow custom error formatting

4. **Partial Updates**: Should the server return the entire form or just error fragments?
   - Full form is simpler but larger payload
   - Fragments require more complex client-side logic

## Follow-up Research: Returning HTML from Connect RPC Handlers

### The Challenge

Connect RPC handlers are designed to return protobuf/JSON responses, not HTML. This creates a fundamental architectural challenge:

- **Current**: Forms POST to `/connect/...` → Connect RPC returns JSON
- **Desired**: Forms need HTML response to re-render with errors → Datastar morphs DOM

### Solution: HTTP Form Handler Wrapper (Recommended)

Create a separate HTTP endpoint layer that wraps Connect RPC calls and returns HTML. This approach:
- ✅ Maintains pure Connect RPC for API consumers
- ✅ Provides HTML responses for form submissions (hypermedia pattern)
- ✅ Full control over error display and re-rendering
- ✅ Leverages Datastar's HTML patching capabilities

### Architecture

```
Browser Form (Datastar)
    ↓ POST /forms/login (HTML expected)
    ↓
HTTP Form Handler
    ↓ Calls internally
Connect RPC Service (pure JSON/proto)
    ↓ Returns error
HTTP Handler parses violations
    ↓ Re-renders form with errors
Returns HTML Fragment
    ↓
Datastar patches DOM
```

### Implementation Example

#### 1. HTTP Form Handler (`handlers/forms/login.go`)

```go
package forms

import (
    "github.com/labstack/echo/v4"
    "connectrpc.com/connect"

    loginv1 "your-project/pkg/proto/com/datastarui/v1/forms/login"
    authconnect "your-project/pkg/proto/com/datastarui/v1/auth/authconnect"
    "your-project/forms/login"
    "your-project/utils"
)

type LoginHandler struct {
    authService authconnect.AuthServiceClient
}

func NewLoginHandler(authService authconnect.AuthServiceClient) *LoginHandler {
    return &LoginHandler{authService: authService}
}

func (h *LoginHandler) HandleLoginForm(c echo.Context) error {
    // Parse form data
    var formData struct {
        Email    string `json:"email" form:"email"`
        Password string `json:"password" form:"password"`
    }

    if err := c.Bind(&formData); err != nil {
        return err
    }

    // Call Connect RPC service internally
    resp, err := h.authService.Login(
        c.Request().Context(),
        connect.NewRequest(&loginv1.LoginRequest{
            Email:    formData.Email,
            Password: formData.Password,
        }),
    )

    if err != nil {
        // Parse validation violations from Connect error
        violations := utils.ParseConnectViolations(err)

        // Re-render form component with errors, preserving user input
        component := login.LoginForm(login.LoginFormArgs{
            Email:         formData.Email,    // Preserve user input
            Password:      formData.Password, // Preserve user input
            EmailError:    violations["email"],
            PasswordError: violations["password"],
        })

        // Render as HTML
        c.Response().Header().Set("Content-Type", "text/html")
        return component.Render(c.Request().Context(), c.Response().Writer)
    }

    // Success case - return success fragment or redirect
    c.Response().Header().Set("Content-Type", "text/html")
    return c.HTML(200, `<div class="text-success">Login successful! Redirecting...</div>`)
}
```

#### 2. Register HTTP Route (`main.go`)

```go
// Register HTTP form handlers (separate from Connect RPC)
formHandlers := forms.NewLoginHandler(authService)
e.POST("/forms/login", formHandlers.HandleLoginForm)

// Connect RPC remains unchanged (for API consumers)
e.Any("/connect/*", echo.WrapHandler(http.StripPrefix("/connect", unwrappedHandler)))
```

#### 3. Update Form Component (`forms/login/login_form.templ`)

```go
type LoginFormArgs struct {
    Email         string
    Password      string
    EmailError    string
    PasswordError string
}

templ LoginForm(args LoginFormArgs) {
    {{
        formID := "login"
        signals := form.SignalsWithFormId(formID, LoginFormSignals{
            Email:    args.Email,    // Preserve from server
            Password: args.Password, // Preserve from server
        })
    }}

    <div
        id="login-form-container"
        data-signals={ signals.DataSignals }
    >
        @form.Form(form.FormArgs{
            ID:     formID,
            Action: "/forms/login",  // HTTP endpoint, not /connect/...
            Attributes: templ.Attributes{
                "data-target": "#login-form-container",  // Datastar patches here
            },
        }) {
            @form.FormItem() {
                @form.FormLabel(form.FormLabelArgs{
                    For:      "email",
                    HasError: args.EmailError != "",
                })
                Email
                }

                @input.Input(input.InputArgs{
                    Type:     "email",
                    ID:       "email",
                    Name:     "email",
                    FormID:   formID,
                    Required: true,
                    Attributes: templ.Attributes{
                        "aria-invalid": templ.Bool(args.EmailError != ""),
                    },
                })

                if args.EmailError != "" {
                    @form.FormMessage(form.FormMessageArgs{
                        Message: args.EmailError,
                    })
                }
            }

            @form.FormItem() {
                @form.FormLabel(form.FormLabelArgs{
                    For:      "password",
                    HasError: args.PasswordError != "",
                })
                Password
                }

                @input.Input(input.InputArgs{
                    Type:     "password",
                    ID:       "password",
                    Name:     "password",
                    FormID:   formID,
                    Required: true,
                    Attributes: templ.Attributes{
                        "aria-invalid": templ.Bool(args.PasswordError != ""),
                    },
                })

                if args.PasswordError != "" {
                    @form.FormMessage(form.FormMessageArgs{
                        Message: args.PasswordError,
                    })
                }
            }

            @button.Button(button.ButtonArgs{
                Type: "submit",
                Attributes: templ.Attributes{
                    "data-attr-disabled": "$fetching",
                },
            }) {
                <span data-show={ "!$fetching" }>Sign In</span>
                <span data-show={ "$fetching" }>Signing in...</span>
            }
        }
    </div>
}
```

#### 4. Violation Parser Utility (`utils/connect_errors.go`)

```go
package utils

import (
    "encoding/base64"
    "errors"

    "connectrpc.com/connect"
    "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
    "google.golang.org/protobuf/proto"
)

// ParseConnectViolations extracts field-level validation errors from Connect error
func ParseConnectViolations(err error) map[string]string {
    fieldErrors := make(map[string]string)

    var connectErr *connect.Error
    if !errors.As(err, &connectErr) {
        return fieldErrors
    }

    // Iterate through error details
    for _, detail := range connectErr.Details() {
        if detail.Type() != "buf.validate.Violations" {
            continue
        }

        // Decode base64 value
        violationsBytes, err := base64.StdEncoding.DecodeString(detail.Value())
        if err != nil {
            continue
        }

        // Unmarshal as Violations protobuf
        violations := &validate.Violations{}
        if err := proto.Unmarshal(violationsBytes, violations); err != nil {
            continue
        }

        // Extract field name and message from each violation
        for _, violation := range violations.GetViolations() {
            fieldName := extractFieldName(violation)
            if fieldName != "" {
                fieldErrors[fieldName] = violation.GetMessage()
            }
        }
    }

    return fieldErrors
}

// extractFieldName gets the field name from violation.field.elements
func extractFieldName(violation *validate.Violation) string {
    field := violation.GetField()
    if field == nil {
        return ""
    }

    elements := field.GetElements()
    if len(elements) == 0 {
        return ""
    }

    // For simple fields, take first element's field name
    return elements[0].GetFieldName()
}
```

### Key Benefits

1. **Dual Interface**
   - HTTP form endpoints for browser interactions (HTML responses)
   - Connect RPC endpoints for API consumers (JSON/proto responses)

2. **True Hypermedia**
   - Server returns HTML with complete UI state
   - Datastar morphs DOM efficiently
   - No client-side error parsing logic

3. **Separation of Concerns**
   - RPC services remain pure (business logic only)
   - HTTP handlers manage presentation (HTML rendering)
   - Easy to test both layers independently

4. **Input Preservation**
   - Server echoes back user input in re-rendered form
   - `data-bind` maintains values across morphing
   - No lost data on validation errors

### Data Flow Example

```
1. User submits form with invalid data
   Email: "not-an-email"
   Password: "short"

2. Browser POSTs to /forms/login
   {"email": "not-an-email", "password": "short"}

3. HTTP handler calls Connect RPC internally
   authService.Login(ctx, request)

4. Validation interceptor detects violations
   - email: "value must be a valid email"
   - password: "value length must be at least 8 characters"

5. HTTP handler receives Connect error
   ParseConnectViolations(err) → map[email:..., password:...]

6. HTTP handler re-renders form with errors
   LoginForm(args) with EmailError and PasswordError set

7. Returns HTML to browser
   Content-Type: text/html
   <div id="login-form-container">...</div>

8. Datastar receives HTML response
   Patches #login-form-container with new content

9. User sees form with:
   - Email field shows "not-an-email" (preserved)
   - Password field shows "short" (preserved)
   - Error messages under each field
   - Red borders on invalid inputs
```

## Recommended Implementation Steps

1. **Create Violation Parser Utility**
   - Implement `ParseConnectViolations()` to extract field errors
   - Decode base64 violations from Connect error details
   - Map proto field names to form field names
   - Return `map[string]string` of field errors

2. **Create HTTP Form Handler**
   - Add new package `handlers/forms`
   - Implement `HandleLoginForm()` that wraps RPC call
   - Parse violations on error and re-render form
   - Return HTML with `Content-Type: text/html`

3. **Update Form Component Args**
   - Add error fields to `LoginFormArgs` struct
   - Update component to display field-level errors
   - Set `HasError` on labels, `aria-invalid` on inputs
   - Use `FormMessage` for error display

4. **Register HTTP Routes**
   - Add POST `/forms/login` route in `main.go`
   - Keep existing Connect RPC routes for API consumers
   - Update form action to point to HTTP endpoint

5. **Test Implementation**
   - Verify input values are preserved via `data-bind`
   - Confirm error messages appear correctly under fields
   - Test successful submission flow
   - Validate Datastar morphing behavior