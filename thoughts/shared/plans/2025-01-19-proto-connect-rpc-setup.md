# Proto Connect RPC Setup for DatastarUI

## Overview

Set up a simplified Connect RPC API infrastructure for DatastarUI that automatically generates API handlers from proto service definitions. The implementation will support authentication/user profile forms, handle form submissions, and enable real-time data updates via Datastar SSE connections. This is a simplified version of the monorepo's Connect RPC setup, optimized for DatastarUI's server-side rendering model with Datastar reactivity.

## Current State Analysis

**DatastarUI Today:**
- Pure server-side rendering with templ templates
- No API endpoints exist (only page rendering routes)
- Echo web framework with basic middleware (logger, recover, CORS, mobile detection)
- All interactivity handled client-side via Datastar signals
- Components follow 3-4 file pattern (template, args, variants, expressions)

**What's Missing:**
- No proto definitions or code generation infrastructure
- No Connect RPC handlers or service pattern
- No API request/response validation
- No form submission endpoints
- No real-time server-to-client communication

**Monorepo Pattern to Follow:**
From `/mnt/d/cdev/monorepo/thoughts/shared/research/2025-01-10-connect-rpc-setup.md`, the monorepo uses:
- Proto-first development with buf toolchain
- ServiceHandler interface pattern for consistent registration
- Connect RPC with interceptor chains
- Proto-validate for automatic request validation
- Generated code in `pkg/proto/`

## Desired End State

**Architecture:**
```
datastarui/
├── buf.yaml                     # Buf workspace config (root)
├── buf.gen.yaml                 # Code generation config (root)
│
├── proto/
│   └── com/
│       └── datastarui/
│           └── v1/
│               ├── auth/
│               │   └── auth.proto          # Service definitions
│               └── forms/
│                   └── login/
│                       └── login.proto     # Form message definitions
│
├── forms/                       # Templ form components
│   └── login/
│       └── login_form.templ
│
├── pkg/proto/                   # Generated code (COMMITTED to git)
│   └── com/
│       └── datastarui/
│           └── v1/
│               ├── auth/
│               │   ├── auth.pb.go
│               │   └── authconnect/
│               │       └── auth.connect.go
│               └── forms/
│                   └── login/
│                       └── login.pb.go
│
├── api/
│   ├── handler/
│   │   └── handler.go           # ServiceHandler interface
│   ├── interceptors/
│   │   └── validation.go        # Proto validation interceptor
│   └── services/
│       └── auth/
│           └── service.go       # Auth service implementation
│
└── main.go                      # Echo server + Connect RPC mux
```

**API Endpoints:**
- `/connect/com.datastarui.v1.auth.AuthService/Login` - Login form submission
- `/connect/com.datastarui.v1.auth.AuthService/Signup` - Signup form submission

**ServiceHandler Pattern:**
```go
type ServiceHandler interface {
    Handler() (string, http.Handler)
}

// Implementation:
func (s *AuthService) Handler() (string, http.Handler) {
    return authconnect.NewAuthServiceHandler(s,
        connect.WithInterceptors(
            interceptors.ValidationInterceptor(),
        ))
}

// Registration:
connectMux := http.NewServeMux()
connectMux.Handle(authService.Handler())
e.Any("/connect/*", echo.WrapHandler(http.StripPrefix("/connect", connectMux)))
```

**Proto Organization Pattern:**
- Form messages: `proto/com/datastarui/v1/forms/{formname}/{formname}.proto`
- Services: `proto/com/datastarui/v1/{service}/{service}.proto`
- Services import form messages: `import "com/datastarui/v1/forms/login/login.proto"`
- Templ forms: `forms/{formname}/{formname}_form.templ`

### Verification Criteria

**Proto Generation Works:**
- `just build` runs `buf generate` successfully
- Generated code appears in `pkg/proto/com/datastarui/v1/`
- No compilation errors after generation

**Connect RPC Endpoints Respond:**
- `curl -X POST http://localhost:4242/connect/com.datastarui.v1.auth.AuthService/Login` returns valid response
- Invalid requests return validation errors with field details

**Form Submission Works:**
- Login form on `/login` submits via Datastar to Connect RPC endpoint
- Validation errors display in form UI
- Successful login stores token/session

**Datastar Integration:**
- Forms use Form component with `ContentType: "json"` for Connect RPC
- Datastar serializes signals to JSON automatically (no custom code needed)
- Connect RPC unmarshals JSON to proto via protojson (built-in, no reflection)
- Server responses can update Datastar signals
- Real-time updates flow via SSE (future phase)
- **No runtime reflection** - pure type-safe JSON serialization

## What We're NOT Doing

- **NO TypeScript client generation** - Full-stack Go + Datastar only
- **NO production authentication** - Simple placeholder for sandbox/demo
- **NO authorization system** - No authz decider or permission checks initially
- **NO multi-tenancy** - No tenant ID headers or isolation
- **NO Vanguard transcoding** - No REST/gRPC compatibility layer
- **NO OpenTelemetry tracing** - Keep interceptors minimal initially
- **NO Sentry integration** - Simple error handling only
- **NO custom authorization annotations** - Proto-validate only

## Implementation Approach

Follow the monorepo's proven patterns but simplified:
1. Use buf for proto compilation (matching monorepo's buf.gen.yaml)
2. Implement ServiceHandler interface (simplified, no authz.Decider)
3. Use proto-validate for automatic request validation
4. Mount Connect RPC to Echo via http.ServeMux
5. Integrate with Datastar forms via JSON serialization (contentType: 'json')
6. Match form protos with templ components via naming convention

**Important: Avoid Runtime Reflection**
- Connect RPC already handles JSON→Proto via protojson (no custom reflection needed)
- Datastar serializes signals to JSON when using contentType: 'json'
- No form-data-to-proto conversion interceptors (would require runtime reflection)
- Keep the integration simple and type-safe

---

## Phase 1: Proto Infrastructure Setup

### Overview
Create proto directory structure, configure buf toolchain at root, and define initial form messages and authentication service. This establishes the foundation for all future proto-driven development.

### Changes Required

#### 1. Create Root Buf Configuration
**File**: `buf.yaml` (root)

```yaml
version: v2
modules:
  - path: proto
deps:
  - buf.build/bufbuild/protovalidate
lint:
  use:
    - STANDARD
  except:
    - FIELD_NOT_REQUIRED
breaking:
  use:
    - WIRE_JSON
```

**File**: `buf.gen.yaml` (root)

```yaml
version: v2
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/coreycole/datastarui/pkg/proto
plugins:
  # Generate Go protobuf code
  - remote: buf.build/protocolbuffers/go:v1.36.6
    out: pkg/proto
    opt: paths=source_relative

  # Generate Connect RPC Go code
  - remote: buf.build/connectrpc/go:v1.18.1
    out: pkg/proto
    opt: paths=source_relative
```

#### 2. Create Login Form Proto
**File**: `proto/com/datastarui/v1/forms/login/login.proto`

```proto
syntax = "proto3";

package com.datastarui.v1.forms.login;

import "buf/validate/validate.proto";

option go_package = "github.com/coreycole/datastarui/pkg/proto/com/datastarui/v1/forms/login;loginv1";

// LoginRequest contains user credentials for authentication
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

// LoginResponse contains the authentication result
message LoginResponse {
  // Session token for authenticated requests
  string token = 1;

  // User information
  User user = 2;
}

// SignupRequest contains new user registration data
message SignupRequest {
  string email = 1 [(buf.validate.field) = {
    required: true,
    string: {email: true}
  }];

  string password = 2 [(buf.validate.field) = {
    required: true,
    string: {
      min_len: 8,
      max_len: 128
    }
  }];

  string name = 3 [(buf.validate.field) = {
    required: true,
    string: {min_len: 1}
  }];
}

// SignupResponse contains the registration result
message SignupResponse {
  string token = 1;
  User user = 2;
}

// User represents a user account
message User {
  string id = 1;
  string email = 2;
  string name = 3;
}
```

#### 3. Create Auth Service Proto
**File**: `proto/com/datastarui/v1/auth/auth.proto`

```proto
syntax = "proto3";

package com.datastarui.v1.auth;

import "com/datastarui/v1/forms/login/login.proto";

option go_package = "github.com/coreycole/datastarui/pkg/proto/com/datastarui/v1/auth;authv1";

// AuthService handles user authentication
service AuthService {
  // Login authenticates a user with email and password
  rpc Login(com.datastarui.v1.forms.login.LoginRequest)
      returns (com.datastarui.v1.forms.login.LoginResponse);

  // Signup creates a new user account
  rpc Signup(com.datastarui.v1.forms.login.SignupRequest)
      returns (com.datastarui.v1.forms.login.SignupResponse);
}
```

#### 4. Update justfile
**File**: `justfile`
**Changes**: Add proto generation to build command

```bash
build:
	buf generate           # Add this line (runs from root)
	templ generate
	go build -o datastarui
	just build-tailwind

# New command for proto-only generation
proto:
	buf generate
```

#### 5. Install buf CLI
**File**: `justfile`
**Changes**: Add buf to install command

```bash
install:
	pnpm install
	go install github.com/air-verse/air@latest
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/bufbuild/buf/cmd/buf@latest    # Add this
	go get ./...
	go mod tidy
	go mod download
```

#### 6. Add Connect RPC Dependencies
**File**: `go.mod`
**Changes**: Add dependencies (run via `go get`)

```bash
go get connectrpc.com/connect@latest
go get buf.build/go/protovalidate@latest
go get google.golang.org/protobuf@latest
```

### Success Criteria

#### Automated Verification:
- [ ] Buf CLI installed: `buf --version` returns version
- [ ] Dependencies fetched: `go mod tidy` completes without errors
- [ ] Proto linting passes: `buf lint`
- [ ] Breaking change detection works: `buf breaking --against .git#branch=main`
- [ ] Code generation succeeds: `buf generate` creates files in `pkg/proto/`
- [ ] Generated code compiles: `go build ./pkg/proto/...` succeeds

#### Manual Verification:
- [ ] Directory structure matches plan
- [ ] `buf.yaml` and `buf.gen.yaml` are at root
- [ ] Proto validation annotations are correctly formatted
- [ ] Generated files appear in `pkg/proto/com/datastarui/v1/forms/login/`
- [ ] `login.pb.go` contains User, LoginRequest, LoginResponse structs
- [ ] Generated files appear in `pkg/proto/com/datastarui/v1/auth/`
- [ ] `authconnect/auth.connect.go` contains AuthServiceHandler interface
- [ ] Import path in auth.proto correctly references login.proto

**Implementation Note**: After completing this phase and all automated verification passes, pause here for manual confirmation from the human that the manual testing was successful before proceeding to the next phase.

---

## Phase 2: API Infrastructure Foundation

### Overview
Create the service handler interface, validation interceptor, and basic API infrastructure. This establishes the pattern that all services will follow.

### Changes Required

#### 1. Create ServiceHandler Interface
**File**: `api/handler/handler.go`

```go
package handler

import "net/http"

// ServiceHandler defines the interface that all Connect RPC services must implement.
// This enables automatic registration and consistent middleware application.
type ServiceHandler interface {
	// Handler returns the service path and configured http.Handler.
	// The returned string is the service path (e.g., "/com.datastarui.v1.auth.AuthService/").
	// The returned http.Handler includes all interceptors.
	Handler() (string, http.Handler)
}
```

#### 2. Create Validation Interceptor
**File**: `api/interceptors/validation.go`

```go
package interceptors

import (
	"context"
	"errors"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
)

// ValidationInterceptor returns a Connect interceptor that validates all
// incoming requests using proto-validate rules defined in proto files.
// Invalid requests are rejected with CodeInvalidArgument and detailed error messages.
func ValidationInterceptor() connect.UnaryInterceptorFunc {
	validator, err := protovalidate.New()
	if err != nil {
		panic("failed to initialize validator: " + err.Error())
	}

	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Validate the request message
			if err := validator.Validate(req.Any().(proto.Message)); err != nil {
				connectErr := connect.NewError(connect.CodeInvalidArgument, err)

				// Add detailed validation violations to error
				if validationErr := new(protovalidate.ValidationError); errors.As(err, &validationErr) {
					if detail, err := connect.NewErrorDetail(validationErr.ToProto()); err == nil {
						connectErr.AddDetail(detail)
					}
				}

				return nil, connectErr
			}

			return next(ctx, req)
		}
	}
}
```

#### 3. Create Logging Interceptor
**File**: `api/interceptors/logging.go`

```go
package interceptors

import (
	"context"
	"log"
	"time"

	"connectrpc.com/connect"
)

// LoggingInterceptor logs all RPC calls with method name and duration.
func LoggingInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			method := req.Spec().Procedure

			log.Printf("RPC: %s started", method)

			resp, err := next(ctx, req)

			duration := time.Since(start)
			if err != nil {
				log.Printf("RPC: %s failed in %v: %v", method, duration, err)
			} else {
				log.Printf("RPC: %s succeeded in %v", method, duration)
			}

			return resp, err
		}
	}
}
```

### Success Criteria

#### Automated Verification:
- [ ] Package compiles: `go build ./api/...` succeeds
- [ ] No linting errors: `golangci-lint run ./api/...` (if available)
- [ ] Imports resolve correctly: `go list -m all | grep protovalidate` shows dependency

#### Manual Verification:
- [ ] `handler.Handler` interface is well-documented
- [ ] `ValidationInterceptor` correctly wraps Connect unary functions
- [ ] Error messages include field violation details
- [ ] Logging interceptor format is readable and useful

**Implementation Note**: After completing this phase and all automated verification passes, pause here for manual confirmation from the human that the manual testing was successful before proceeding to the next phase.

---

## Phase 3: Implement Auth Service

### Overview
Create the authentication service implementation with placeholder logic for login and signup. This demonstrates the service pattern and validates the infrastructure.

### Changes Required

#### 1. Create Auth Service Implementation
**File**: `api/services/auth/service.go`

```go
package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"

	"connectrpc.com/connect"

	authv1 "github.com/coreycole/datastarui/pkg/proto/com/datastarui/v1/auth"
	"github.com/coreycole/datastarui/pkg/proto/com/datastarui/v1/auth/authconnect"
	loginv1 "github.com/coreycole/datastarui/pkg/proto/com/datastarui/v1/forms/login"

	"github.com/coreycole/datastarui/api/handler"
	"github.com/coreycole/datastarui/api/interceptors"
)

// Compile-time interface checks
var (
	_ authconnect.AuthServiceHandler = (*Service)(nil)
	_ handler.ServiceHandler          = (*Service)(nil)
)

// Service implements the AuthService RPC handlers.
// This is a placeholder implementation for demonstration purposes.
type Service struct {
	authconnect.UnimplementedAuthServiceHandler

	// In-memory user storage (placeholder)
	mu    sync.RWMutex
	users map[string]*loginv1.User // email -> user
}

// NewService creates a new auth service instance.
func NewService() *Service {
	return &Service{
		users: make(map[string]*loginv1.User),
	}
}

// Handler returns the service path and HTTP handler with interceptors.
func (s *Service) Handler() (string, http.Handler) {
	return authconnect.NewAuthServiceHandler(s,
		connect.WithInterceptors(
			interceptors.LoggingInterceptor(),
			interceptors.ValidationInterceptor(),
		))
}

// Login authenticates a user and returns a session token.
// This is a placeholder implementation that accepts any email/password
// as long as the user exists.
func (s *Service) Login(
	ctx context.Context,
	req *connect.Request[loginv1.LoginRequest],
) (*connect.Response[loginv1.LoginResponse], error) {
	email := req.Msg.Email
	password := req.Msg.Password

	s.mu.RLock()
	user, exists := s.users[email]
	s.mu.RUnlock()

	if !exists {
		return nil, connect.NewError(
			connect.CodeUnauthenticated,
			fmt.Errorf("user not found: %s", email),
		)
	}

	// Placeholder: accept any password
	// In production, verify hashed password
	_ = password

	// Generate session token
	token, err := generateToken()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&loginv1.LoginResponse{
		Token: token,
		User:  user,
	}), nil
}

// Signup creates a new user account.
// This is a placeholder implementation that stores users in memory.
func (s *Service) Signup(
	ctx context.Context,
	req *connect.Request[loginv1.SignupRequest],
) (*connect.Response[loginv1.SignupResponse], error) {
	email := req.Msg.Email
	password := req.Msg.Password
	name := req.Msg.Name

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if user already exists
	if _, exists := s.users[email]; exists {
		return nil, connect.NewError(
			connect.CodeAlreadyExists,
			fmt.Errorf("user already exists: %s", email),
		)
	}

	// Generate user ID
	userID, err := generateToken()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Create user (placeholder: not hashing password)
	_ = password // In production, hash the password

	user := &loginv1.User{
		Id:    userID,
		Email: email,
		Name:  name,
	}

	s.users[email] = user

	// Generate session token
	token, err := generateToken()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&loginv1.SignupResponse{
		Token: token,
		User:  user,
	}), nil
}

// generateToken creates a random session token.
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
```

### Success Criteria

#### Automated Verification:
- [ ] Service compiles: `go build ./api/services/auth` succeeds
- [ ] Implements all proto service methods
- [ ] No unused imports or variables
- [ ] go vet passes: `go vet ./api/services/auth`

#### Manual Verification:
- [ ] Service correctly implements `ServiceHandler` interface
- [ ] Login rejects non-existent users
- [ ] Signup creates users in memory map
- [ ] Tokens are generated securely (crypto/rand)
- [ ] Error codes match Connect conventions
- [ ] Imports use form message types from `loginv1` package

**Implementation Note**: After completing this phase and all automated verification passes, pause here for manual confirmation from the human that the manual testing was successful before proceeding to the next phase.

---

## Phase 4: Integrate Connect RPC with Echo

### Overview
Mount the Connect RPC service multiplexer to Echo and verify API endpoints respond correctly. This integrates the proto-generated APIs with the existing web server.

### Changes Required

#### 1. Update main.go
**File**: `main.go`
**Changes**: Add Connect RPC mux and mount to Echo

```go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// ... existing imports ...

	// Add new imports
	"github.com/coreycole/datastarui/api/handler"
	"github.com/coreycole/datastarui/api/services/auth"
)

// ... existing code ...

func main() {
	// ... existing config loading ...

	e := echo.New()

	// Existing middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isMobile := utils.IsMobile(c)
			ctx := utils.WithMobile(c.Request().Context(), isMobile)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	})

	// NEW: Setup Connect RPC services
	connectMux := http.NewServeMux()

	// Register services
	services := []handler.ServiceHandler{
		auth.NewService(),
	}

	for _, svc := range services {
		path, handler := svc.Handler()
		connectMux.Handle(path, handler)
		log.Printf("Registered Connect RPC service: %s", path)
	}

	// Mount Connect RPC at /connect/*
	e.Any("/connect/*", echo.WrapHandler(http.StripPrefix("/connect", connectMux)))

	// ... existing routes ...

	// Start server
	log.Printf("Starting server on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
```

### Success Criteria

#### Automated Verification:
- [ ] Application compiles: `go build` succeeds
- [ ] Application starts: `go run main.go` runs without panics
- [ ] Service registration logs appear on startup
- [ ] Server listens on port 4242

#### Manual Verification:
- [ ] Test login endpoint with valid request:
```bash
curl -X POST http://localhost:4242/connect/com.datastarui.v1.auth.AuthService/Login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}'
```
Expected: 401 Unauthenticated (user not found)

- [ ] Test signup endpoint:
```bash
curl -X POST http://localhost:4242/connect/com.datastarui.v1.auth.AuthService/Signup \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123", "name": "Test User"}'
```
Expected: 200 OK with token and user object

- [ ] Test validation with invalid email:
```bash
curl -X POST http://localhost:4242/connect/com.datastarui.v1.auth.AuthService/Login \
  -H "Content-Type: application/json" \
  -d '{"email": "notanemail", "password": "short"}'
```
Expected: 400 Invalid Argument with field violation details

- [ ] Test login after signup succeeds
- [ ] Check server logs show RPC method calls

**Implementation Note**: After completing this phase and all automated verification passes, pause here for manual confirmation from the human that the manual testing was successful before proceeding to the next phase.

---

## Phase 5: Create Login Form Component

### Overview
Build a reusable login form component using templ and Datastar that submits to the Connect RPC endpoint. This form will be usable across multiple pages (login page, modal, etc.).

### Changes Required

#### 1. Create Forms Directory
**Directory**: `forms/login/`

This directory will contain the templ form component. The proto definition is at `proto/com/datastarui/v1/forms/login/login.proto` (created in Phase 1).

#### 2. Create Login Form Component
**File**: `forms/login/login_form.templ`

```go
package login

import (
	"github.com/coreycole/datastarui/components/button"
	"github.com/coreycole/datastarui/components/form"
	"github.com/coreycole/datastarui/components/input"
	"github.com/coreycole/datastarui/utils"
)

// LoginFormSignals defines the client-side state for the login form
type LoginFormSignals struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Error    string `json:"error"`
	Loading  bool   `json:"loading"`
}

// LoginFormArgs contains configuration for the login form
type LoginFormArgs struct {
	ID          string // Form ID for signal namespacing
	RedirectURL string // Where to redirect on success
}

templ LoginForm(args LoginFormArgs) {
	{{
		formID := args.ID
		if formID == "" {
			formID = "login"
		}

		signals := utils.Signals(formID, LoginFormSignals{
			Email:    "",
			Password: "",
			Error:    "",
			Loading:  false,
		})
	}}
	<div data-signals={ signals.DataSignals }>
		<!-- Error Message -->
		<div
			data-show={ signals.Signal("error") + " !== ''" }
			data-text={ signals.Signal("error") }
			class="mb-4 rounded-md bg-destructive/10 p-3 text-sm text-destructive"
		></div>

		@form.Form(form.FormArgs{
			ID:          formID,
			Action:      "/connect/com.datastarui.v1.auth.AuthService/Login",
			ContentType: "json",
			Class:       "space-y-4",
		}) {
			<!-- Email Field -->
			@form.FormItem(form.FormItemArgs{}) {
				@form.FormLabel(form.FormLabelArgs{For: "email"}) {
					Email
				}
				@input.Input(input.InputArgs{
					Type:        "email",
					ID:          "email",
					Name:        "email",
					FormID:      formID,
					Placeholder: "you@example.com",
					Required:    true,
				})
			}

			<!-- Password Field -->
			@form.FormItem(form.FormItemArgs{}) {
				@form.FormLabel(form.FormLabelArgs{For: "password"}) {
					Password
				}
				@input.Input(input.InputArgs{
					Type:        "password",
					ID:          "password",
					Name:        "password",
					FormID:      formID,
					Placeholder: "••••••••",
					Required:    true,
				})
			}

			<!-- Submit Button -->
			<div class="pt-2">
				@button.Button(button.ButtonArgs{
					Type:    "submit",
					Variant: "default",
					Class:   "w-full",
					Attributes: templ.Attributes{
						"data-attr-disabled": "$fetching",
					},
				}) {
					<span data-show={ "!$fetching" }>Sign In</span>
					<span data-show={ "$fetching" }>Signing in...</span>
				}
			</div>
		}
	</div>
}
```

**Key patterns:**
- Signals defined with `utils.Signals()` and wrapped in outer div
- Form component with `ContentType: "json"` for Connect RPC
- FormItem/FormLabel for consistent structure and styling
- Input uses `FormID` which auto-generates `data-bind="login.email"`
- Datastar's `$fetching` signal tracks form submission state
- Error display uses signals for reactive updates

#### 3. Create Login Page Using Form Component
**File**: `pages/login/login_page.templ`

```go
package login

import (
	"github.com/coreycole/datastarui/layouts"
	loginform "github.com/coreycole/datastarui/forms/login"
)

templ LoginPage(rootArgs layouts.RootArgs) {
	@layouts.Root(rootArgs) {
		<div class="flex min-h-screen items-center justify-center">
			<div class="w-full max-w-md space-y-8 rounded-lg border bg-card p-8 shadow-lg">
				<div class="space-y-2 text-center">
					<h1 class="text-3xl font-bold">Welcome back</h1>
					<p class="text-muted-foreground">Enter your credentials to sign in</p>
				</div>
				@loginform.LoginForm(loginform.LoginFormArgs{
					ID:          "login",
					RedirectURL: "/dashboard",
				})
				<!-- Signup Link -->
				<p class="text-center text-sm text-muted-foreground">
					Don't have an account?
					<a href="/signup" class="font-medium text-primary hover:underline">
						Sign up
					</a>
				</p>
			</div>
		</div>
	}
}
```

#### 4. Add Login Route to main.go
**File**: `main.go`
**Changes**: Add route for /login

```go
import (
	// ... existing imports ...
	loginpage "github.com/coreycole/datastarui/pages/login"
)

// Add to main() function after existing routes

e.GET("/login", func(c echo.Context) error {
	rootArgs := l.RootArgs{
		CurrentPage:          "login",
		CurrentPath:          c.Request().URL.Path,
		InspectorEnabled:     cfg.DatastarInspectorEnabled,
		DatastarProAvailable: datastarProAvailable,
	}
	component := loginpage.LoginPage(rootArgs)
	return component.Render(c.Request().Context(), c.Response().Writer)
})
```

#### 5. Enhance Form Component for Connect RPC

**File**: `components/form/args.go`
**Changes**: Add ContentType field

```go
type FormArgs struct {
	ID          string           // Form ID
	Action      string           // Form action URL
	ContentType string           // "form" (default) | "json" for Connect RPC
	Class       string           // Additional CSS classes
	Attributes  templ.Attributes // Additional HTML attributes
}
```

**File**: `components/form/form.templ`
**Changes**: Support JSON content type for Connect RPC

```go
if args.Action != "" {
	// Check if there's a target specified in attributes
	target := ""
	if targetAttr, exists := args.Attributes["data-target"]; exists {
		target = fmt.Sprintf(", target: '%s'", targetAttr)
		delete(formAttrs, "data-target")
	}

	// Determine content type (default to 'form')
	contentType := "form"
	if args.ContentType == "json" {
		contentType = "json"
	}

	formAttrs["data-on-submit"] = fmt.Sprintf("@post('%s', {contentType: '%s'%s})",
		templ.SafeURL(args.Action), contentType, target)
}
```

**How it works with Connect RPC:**
- Set `ContentType: "json"` in FormArgs
- Datastar serializes signals to JSON: `{"email": "...", "password": "..."}`
- Connect RPC receives `application/json` and unmarshals via protojson
- **No runtime reflection needed** - Connect handles JSON→Proto natively

### Success Criteria

#### Automated Verification:
- [ ] Templates compile: `templ generate` succeeds
- [ ] Application compiles: `go build` succeeds
- [ ] Login page route registered
- [ ] No template syntax errors

#### Manual Verification:
- [ ] Visit http://localhost:4242/login
- [ ] Login form renders correctly with Form/FormItem/FormLabel components
- [ ] Email and password fields bind to signals (check Datastar inspector: `$login.email`, `$login.password`)
- [ ] Submit button shows loading state when clicked (uses `$fetching` signal)
- [ ] Invalid email shows validation error from server (proto-validate)
- [ ] Short password (<8 chars) shows validation error (proto-validate min_len rule)
- [ ] Valid credentials after signup show success response
- [ ] Error message displays reactively when login fails
- [ ] Network tab shows:
  - POST to /connect/com.datastarui.v1.auth.AuthService/Login
  - Content-Type: application/json
  - Request body is JSON: `{"email": "...", "password": "..."}`
  - Response is JSON (Connect RPC JSON format)
- [ ] Form component can be reused (import in another page)
- [ ] **Verify no runtime reflection**: Check server logs show no reflection warnings

**Implementation Note**: After completing this phase and all automated verification passes, pause here for manual confirmation from the human that the manual testing was successful before proceeding to the next phase.

---

## Testing Strategy

### Unit Tests

**Test proto validation:**
```go
// api/interceptors/validation_test.go
func TestValidationInterceptor_InvalidEmail(t *testing.T) {
	req := &loginv1.LoginRequest{
		Email:    "not-an-email",
		Password: "password123",
	}
	// Assert validation error
}
```

**Test service handlers:**
```go
// api/services/auth/service_test.go
func TestAuthService_Signup(t *testing.T) {
	svc := NewService()
	req := &loginv1.SignupRequest{
		Email:    "test@example.com",
		Password: "securepass123",
		Name:     "Test User",
	}
	// Assert user created
}
```

### Integration Tests

**Test full Connect RPC flow:**
```go
func TestLoginFlow(t *testing.T) {
	// 1. Start test server
	// 2. Call Signup endpoint
	// 3. Verify user created
	// 4. Call Login endpoint
	// 5. Verify token returned
}
```

### Manual Testing Steps

1. **Proto Generation**:
   - Run `just build`
   - Verify `pkg/proto/` contains generated code
   - Check for compilation errors
   - Verify generated code is committed

2. **API Endpoints**:
   - Test signup with valid data
   - Test login after signup
   - Test validation errors (invalid email, short password)
   - Verify error messages include field details

3. **Datastar Integration**:
   - Visit /login page
   - Open Datastar inspector (if enabled)
   - Watch signals update on form input
   - Submit form and verify loading state
   - Check network requests use correct format

4. **Form Reusability**:
   - Import LoginForm in another page
   - Verify form works with different IDs
   - Test multiple forms on same page

## Performance Considerations

**Proto Validation**:
- Validator instance created once per interceptor (not per request)
- Validation overhead is minimal (<1ms per request typically)
- Error messages are detailed but not verbose

**In-Memory Storage**:
- Current placeholder uses sync.RWMutex for thread safety
- Suitable for development/demo only
- Replace with database in production

**Connect RPC Protocol**:
- Uses HTTP/2 binary protobuf by default (efficient)
- Smaller payload sizes than JSON
- Multiplexing support for concurrent requests

## Migration Notes

**Moving to Production Authentication**:
1. Replace in-memory user storage with database (e.g., PostgreSQL)
2. Add password hashing with bcrypt
3. Implement JWT token generation and validation
4. Add session management
5. Create authentication interceptor for protected endpoints

**Adding Authorization**:
1. Define roles in proto (e.g., admin, user)
2. Create authorization interceptor
3. Add role checks in service methods
4. Implement permission system

**Adding SSE for Real-Time Updates**:
1. Create `/sse/*` endpoints separate from Connect RPC
2. Use Datastar's SSE support for server-push
3. Send proto messages serialized as JSON over SSE
4. Example: Real-time profile updates when changed by admin

**Adding More Forms**:
1. Create proto at `proto/com/datastarui/v1/forms/{formname}/{formname}.proto`
2. Create templ at `forms/{formname}/{formname}_form.templ`
3. Import form messages in service proto
4. Implement RPC handlers using form messages

## References

- Original monorepo research: `/mnt/d/cdev/monorepo/thoughts/shared/research/2025-01-10-connect-rpc-setup.md`
- Connect RPC docs: https://connectrpc.com/docs/
- Proto-validate docs: https://buf.build/bufbuild/protovalidate
- Datastar docs: https://data-star.dev/
- Buf CLI docs: https://buf.build/docs/

## Next Steps (Future Work)

After this implementation is complete, consider:

1. **Database Integration**: Replace in-memory storage with sqlc + PostgreSQL
2. **Session Management**: Add Redis for session storage
3. **Real-Time Updates**: Implement SSE endpoints for server-push notifications
4. **Additional Forms**: Create signup, profile update, password reset forms
5. **Rate Limiting**: Add rate limiting interceptor
6. **Metrics**: Add Prometheus metrics for RPC calls
7. **API Documentation**: Generate OpenAPI/Swagger docs from proto
8. **Form Validation UI**: Better client-side validation error display
