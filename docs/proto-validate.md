# Proto Validation with Datastar Forms

This document explains how DatastarUI handles form submissions with Connect RPC and Protocol Buffers validation, including the automatic extraction of nested form data.

## Overview

When using Datastar with Connect RPC, forms send data as nested signal objects. DatastarUI automatically extracts the relevant data and validates it against protobuf schemas using buf/validate.

## The Challenge

Datastar organizes signals in namespaces for isolation and organization:

```javascript
// Frontend signals
{
  "theme": "dark",
  "login": {
    "email": "user@example.com",
    "password": "password123",
    "error": "",
    "loading": false
  },
  "main_sidebar": {...}
}
```

But Connect RPC expects flat protobuf structures:

```protobuf
message LoginRequest {
  string email = 1;
  string password = 2;
}
```

**The mismatch:** Datastar sends `{"login": {"email": "...", "password": "..."}}` but protobuf expects `{"email": "...", "password": "..."}`.

## The Solution: Three-Layer Architecture with Explicit Form Metadata

DatastarUI uses an explicit `formId` signal to reliably extract form data from nested signal structures. The `formId` is initialized as a **root-level signal** alongside the form's namespaced signals:

```json
{
  "formId": "login",           // ← Root-level signal: identifies the form
  "login": {                   // ← Namespaced form signals
    "email": "user@example.com",
    "password": "password123",
    "error": "",
    "loading": false
  }
}
```

**Key Features:**
- ✅ `formId` is a root-level signal (not nested inside the form namespace)
- ✅ Form data signals remain namespaced under their form ID (`login.*`)
- ✅ Filter pattern includes both: `filterSignals: {include: /^(formId|login\.)/}`
- ✅ Backend reads `formId` to know which key to extract from the payload
- ✅ More reliable than heuristic-based detection (e.g., single-key guessing)
- ✅ Works correctly even with multiple forms submitting simultaneously

## The Solution: Three-Layer Architecture

### 1. Frontend: Signal Initialization and Filtering

#### Step 1: Initialize Signals with formId

Use the `form.SignalsWithFormId()` helper to create signals with `formId` at root level:

**`forms/login/login_form.templ`:**
```go
signals := form.SignalsWithFormId(formID, LoginFormSignals{
    Email:    "",
    Password: "",
    Error:    "",
    Loading:  false,
})
```

This creates the signal structure:
```json
{
  "formId": "login",
  "login": {
    "email": "",
    "password": "",
    "error": "",
    "loading": false
  }
}
```

**Helper Implementation** (`components/form/form.templ`):
```go
func SignalsWithFormId(formID string, signalsStruct any) *utils.SignalManager {
    sanitizedID := strings.ReplaceAll(formID, "-", "_")

    signalsMap := map[string]any{
        "formId":    sanitizedID,  // Root-level signal
        sanitizedID: signalsStruct, // Namespaced signals
    }

    signalsJSON, _ := json.Marshal(signalsMap)

    return &utils.SignalManager{
        ID:          sanitizedID,
        Signals:     signalsStruct,
        DataSignals: string(signalsJSON),
    }
}
```

This returns a `SignalManager` so you can use helper methods like `signals.Signal("email")` which generates `$login.email`.

#### Step 2: Configure Form Component Filtering

**`components/form/form.templ`:**
```go
if args.FormDataFields != nil && len(args.FormDataFields) > 0 {
    // Build regex pattern to match:
    // - formId (exact match)
    // - login.email, login.password, etc. (formID.*)
    filterPattern := fmt.Sprintf("^(formId|%s\\.)", args.ID)
    filterSignals := fmt.Sprintf("filterSignals: {include: /%s/}", filterPattern)

    submitHandler = fmt.Sprintf(
        "@post('%s', {contentType: 'json', %s})",
        args.Action,
        filterSignals
    )
}
```

**Generated HTML:**
```html
<form
    id="login"
    data-on-submit="@post('/connect/...AuthService/Login', {
        contentType: 'json',
        filterSignals: {include: /^(formId|login\.)/}
    })"
>
```

**What this does:**
- `filterSignals: {include: /^(formId|login\.)/}` matches:
  - `formId` - Root-level signal (exact match)
  - `login.email`, `login.password`, etc. - Namespaced signals (with dot)
- Datastar collects and sends: `{"formId": "login", "login": {"email": "...", "password": "...", "error": "", "loading": false}}`
- Excludes global signals like `theme`, `fetching`, `main_sidebar`

### 2. Backend: HTTP Middleware Unwrapping

Before Connect RPC processes the request, HTTP middleware uses an explicit `formId` field to reliably extract the form data:

**`api/middleware/unwrap_form.go`:**
```go
func UnwrapFormData(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Read request body
        bodyBytes, _ := io.ReadAll(r.Body)
        r.Body.Close()

        // Parse as JSON
        var dataMap map[string]interface{}
        json.Unmarshal(bodyBytes, &dataMap)

        // Check if request contains formId metadata
        if formIdValue, hasFormId := dataMap["formId"]; hasFormId {
            if formId, ok := formIdValue.(string); ok && formId != "" {
                // Look for the nested object with the formId key
                if nestedMap, hasNested := dataMap[formId]; hasNested {
                    if formData, ok := nestedMap.(map[string]interface{}); ok {
                        // Extract the form data
                        unwrappedBytes, _ := json.Marshal(formData)

                        // Replace request body with unwrapped data
                        r.Body = io.NopCloser(bytes.NewBuffer(unwrappedBytes))
                        r.ContentLength = int64(len(unwrappedBytes))
                        next.ServeHTTP(w, r)
                        return
                    }
                }
            }
        }

        // No unwrapping needed, pass through original
        r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
        next.ServeHTTP(w, r)
    })
}
```

**Applied in `main.go`:**
```go
// Mount Connect RPC at /connect/* with unwrap middleware
unwrappedHandler := apimw.UnwrapFormData(connectMux)
e.Any("/connect/*", echo.WrapHandler(http.StripPrefix("/connect", unwrappedHandler)))
```

**Transformation:**
- **Input:** `{"formId": "login", "login": {"email": "user@example.com", "password": "password123", "error": "", "loading": false}}`
- **Output:** `{"email": "user@example.com", "password": "password123", "error": "", "loading": false}`

**Why this approach is reliable:**
- ✅ **Explicit metadata**: The `formId` field explicitly tells the backend which key to extract
- ✅ **No guessing**: Unlike single-key heuristics, this works with multiple forms on a page
- ✅ **Future-proof**: Works even if you add more metadata fields to the request (e.g., `requestId`, `timestamp`)
- ✅ **Multiple forms**: Can handle multiple forms submitting simultaneously without conflicts
- ✅ **Debugging**: Easy to see which form sent the request by checking the `formId` value

### 3. Connect RPC: Protobuf Deserialization & Validation

Connect RPC's protojson deserializer automatically:
1. Maps JSON fields to protobuf fields by name (case-sensitive)
2. **Ignores unknown fields** (like `error` and `loading` which aren't in the proto)
3. Unmarshals to the typed protobuf message

**Proto definition (`proto/com/datastarui/v1/forms/login/login.proto`):**
```protobuf
message LoginRequest {
  string email = 1 [(buf.validate.field) = {
    required: true,
    string: {email: true}
  }];

  string password = 2 [(buf.validate.field) = {
    required: true,
    string: {min_len: 8}
  }];
}
```

**Validation interceptor (`api/interceptors/validation.go`):**
```go
func ValidationInterceptor() connect.UnaryInterceptorFunc {
    validator, _ := protovalidate.New()

    return func(next connect.UnaryFunc) connect.UnaryFunc {
        return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
            // Validate the protobuf message
            if err := validator.Validate(req.Any().(proto.Message)); err != nil {
                return nil, connect.NewError(connect.CodeInvalidArgument, err)
            }
            return next(ctx, req)
        }
    }
}
```

## Complete Data Flow

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. User Input                                                   │
│    Email: user@example.com                                      │
│    Password: password123                                        │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. Datastar Signals (data-bind)                                 │
│    $login.email = "user@example.com"                            │
│    $login.password = "password123"                              │
│    $login.error = ""                                            │
│    $login.loading = false                                       │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. Form Submit (includes formId metadata)                       │
│    Sends nested JSON with explicit formId:                      │
│    {                                                            │
│      "formId": "login",                                         │
│      "login": {                                                 │
│        "email": "user@example.com",                             │
│        "password": "password123",                               │
│        "error": "",                                             │
│        "loading": false                                         │
│      }                                                          │
│    }                                                            │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. HTTP Middleware (UnwrapFormData)                             │
│    Reads formId field: "login"                                  │
│    Extracts dataMap["login"] object:                            │
│    {                                                            │
│      "email": "user@example.com",                               │
│      "password": "password123",                                 │
│      "error": "",                                               │
│      "loading": false                                           │
│    }                                                            │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│ 5. Connect RPC Deserialization                                  │
│    Unmarshals JSON to LoginRequest protobuf                     │
│    Maps: email → LoginRequest.Email                             │
│    Maps: password → LoginRequest.Password                       │
│    Ignores: error (not in proto)                                │
│    Ignores: loading (not in proto)                              │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│ 6. Validation Interceptor                                       │
│    ✓ Email format valid (has @)                                 │
│    ✓ Password length >= 8 characters                            │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│ 7. Service Handler                                              │
│    func Login(req *LoginRequest) (*LoginResponse, error) {      │
│        email := req.Email  // "user@example.com"                │
│        password := req.Password  // "password123"               │
│        // ... authenticate user                                 │
│    }                                                            │
└─────────────────────────────────────────────────────────────────┘
```

## Example: Login Form Implementation

### Step 1: Define Proto Message

**`proto/com/datastarui/v1/forms/login/login.proto`:**
```protobuf
syntax = "proto3";

package com.datastarui.v1.forms.login;

import "buf/validate/validate.proto";

message LoginRequest {
  string email = 1 [(buf.validate.field) = {
    required: true,
    string: {email: true}
  }];

  string password = 2 [(buf.validate.field) = {
    required: true,
    string: {min_len: 8}
  }];
}

message LoginResponse {
  string token = 1;
  User user = 2;
}

message User {
  string id = 1;
  string email = 2;
  string name = 3;
}
```

### Step 2: Create Form Component

**`forms/login/login_form.templ`:**
```go
package login

type LoginFormSignals struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Error    string `json:"error"`
    Loading  bool   `json:"loading"`
}

templ LoginForm(args LoginFormArgs) {
    {{
        formID := "login"
        signals := utils.Signals(formID, LoginFormSignals{
            Email:    "",
            Password: "",
            Error:    "",
            Loading:  false,
        })
    }}
    <div data-signals={ signals.DataSignals }>
        <!-- Error Display -->
        <div
            data-show={ signals.Signal("error") + " !== ''" }
            data-text={ signals.Signal("error") }
            class="error-message"
        ></div>

        @form.Form(form.FormArgs{
            ID:             formID,
            Action:         "/connect/com.datastarui.v1.auth.AuthService/Login",
            ContentType:    "json",
            FormDataFields: []string{"email", "password"}, // Only send these fields
            Class:          "space-y-4",
        }) {
            @input.Input(input.InputArgs{
                Type:        "email",
                ID:          "email",
                Name:        "email",
                FormID:      formID,  // Auto-generates data-bind="login.email"
                Required:    true,
            })

            @input.Input(input.InputArgs{
                Type:        "password",
                ID:          "password",
                Name:        "password",
                FormID:      formID,  // Auto-generates data-bind="login.password"
                Required:    true,
            })

            @button.Button(button.ButtonArgs{
                Type: "submit",
                Attributes: templ.Attributes{
                    "data-attr-disabled": "$fetching",
                },
            }) {
                <span data-show="!$fetching">Sign In</span>
                <span data-show="$fetching">Signing in...</span>
            }
        }
    </div>
}
```

### Step 3: Register Service with Interceptors

**`api/services/auth/service.go`:**
```go
func (s *Service) Handler() (string, http.Handler) {
    return authconnect.NewAuthServiceHandler(s,
        connect.WithInterceptors(
            interceptors.LoggingInterceptor(),
            interceptors.ValidationInterceptor(),  // Validates proto constraints
        ))
}

func (s *Service) Login(
    ctx context.Context,
    req *connect.Request[loginv1.LoginRequest],
) (*connect.Response[loginv1.LoginResponse], error) {
    // At this point, req.Msg is fully validated and typed
    email := req.Msg.Email       // ✓ Valid email format
    password := req.Msg.Password // ✓ At least 8 characters

    // ... authenticate user
}
```

## Key Benefits

### 1. **Type Safety**
- Protobuf provides compile-time type checking
- Frontend signals match proto field names
- Validation rules defined once in proto schema

### 2. **Automatic Validation**
- Email format validation: `string: {email: true}`
- Length constraints: `string: {min_len: 8}`
- Required fields: `required: true`
- No manual validation code needed

### 3. **Clean Separation**
- **UI State** (`error`, `loading`): Signals for frontend reactivity
- **Form Data** (`email`, `password`): Sent to backend, validated by proto
- Protobuf automatically ignores UI-only fields

### 4. **Error Handling**
- Validation errors are detailed and structured
- Frontend receives field-level error messages
- Can display inline validation feedback

## Validation Error Example

**Request (as sent from frontend):**
```json
{
  "formId": "login",
  "login": {
    "email": "invalid-email",
    "password": "short",
    "error": "",
    "loading": false
  }
}
```

**Request (after middleware unwrapping):**
```json
{
  "email": "invalid-email",
  "password": "short",
  "error": "",
  "loading": false
}
```

**Response:**
```json
{
  "code": "invalid_argument",
  "message": "validation error:\n - email: value must be a valid email [string.email]\n - password: value length must be at least 8 characters [string.min_len]",
  "details": [{
    "type": "buf.validate.Violations",
    "value": "...",
    "debug": {
      "violations": [
        {
          "field": {"elements": [{"fieldName": "email"}]},
          "ruleId": "string.email",
          "message": "value must be a valid email"
        },
        {
          "field": {"elements": [{"fieldName": "password"}]},
          "ruleId": "string.min_len",
          "message": "value length must be at least 8 characters"
        }
      ]
    }
  }]
}
```

## Best Practices

### 1. **Signal Naming Convention**
- Use lowercase with underscores for form IDs: `user_profile`, `checkout_form`
- Never use hyphens or dots in form IDs
- Form ID becomes signal namespace: `$user_profile.email`

### 2. **Separate UI State from Form Data**
Define signals that split concerns:
```go
type LoginFormSignals struct {
    // Form data (sent to backend)
    Email    string `json:"email"`
    Password string `json:"password"`

    // UI state (frontend only, ignored by protobuf)
    Error    string `json:"error"`
    Loading  bool   `json:"loading"`
}
```

### 3. **Explicit Field Mapping**
Always specify `FormDataFields` to document what's being sent:
```go
FormDataFields: []string{"email", "password"}  // Clear intent
```

### 4. **Proto-First Design**
1. Define your proto message with validation rules
2. Generate Go code with `buf generate`
3. Create matching signal struct in templ
4. Protobuf validation is your source of truth

### 5. **Consistent Error Handling**
Display validation errors in the UI:
```html
<div data-show="$login.error !== ''" data-text="$login.error"></div>
```

## Debugging

### Check Request Payload
Open browser DevTools → Network tab → Find POST request:
```json
// Should see nested structure with formId metadata
{
  "formId": "login",
  "login": {
    "email": "user@example.com",
    "password": "password123",
    "error": "",
    "loading": false
  }
}
```

### Check Server Logs
```
2025/10/20 03:31:40 RPC: /com.datastarui.v1.auth.AuthService/Login started
2025/10/20 03:31:40 RPC: /com.datastarui.v1.auth.AuthService/Login succeeded in 15ms
```

### Common Issues

**Empty request body `{}`:**
- Check that `FormDataFields` is set in your form component
- Verify `$$signals(/^formID\./)` regex pattern matches your signal paths
- Should be `/^formID\./` (with dot) not `/^formID$/`
- Check browser console for JavaScript errors

**Missing formId in request:**
- Verify the form generates: `Object.assign({formId: 'login'}, $$signals(...))`
- Check DevTools → Network tab to see actual request payload
- Should include `"formId": "login"` at the top level

**Middleware not unwrapping:**
- Check that `formId` field is present in request
- Verify middleware is applied in correct order (before Connect RPC)
- Check server logs for middleware errors

**Validation errors on valid data:**
- Check proto field names match JSON keys (case-sensitive)
- Verify `FormDataFields` includes all required proto fields
- Remember: middleware extracts `dataMap[formId]`, not the whole body

**Extra fields causing errors:**
- Protobuf ignores unknown fields by default (like `error`, `loading`)
- Only an issue if using strict unmarshal options

## References

- [Protocol Buffers](https://protobuf.dev/)
- [buf/validate](https://github.com/bufbuild/protovalidate)
- [Connect RPC](https://connectrpc.com/)
- [Datastar Documentation](https://data-star.dev/)
