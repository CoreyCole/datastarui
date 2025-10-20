---
date: 2025-10-19T17:49:35-07:00
researcher: Claude
git_commit: addb91a65cb927f8da9e29ead2d1a617a4e723ab
branch: cc/proto-forms-planned
repository: datastarui
topic: "Proto Connect RPC Setup for DatastarUI"
tags: [implementation, proto, connect-rpc, datastar, forms, authentication]
status: in_progress
last_updated: 2025-10-19
last_updated_by: Claude
type: implementation_strategy
---

# Handoff: Proto Connect RPC Setup

## Task(s)

**Implementation Plan**: `/mnt/d/cdev/datastarui/thoughts/shared/plans/2025-01-19-proto-connect-rpc-setup.md`

### Status by Phase:
- ✅ **Phase 1: Proto Infrastructure Setup** - COMPLETED
  - Created buf.yaml and buf.gen.yaml configuration
  - Created proto definitions for login forms and auth service
  - Added Connect RPC dependencies
  - Successfully generated proto code in pkg/proto/

- ✅ **Phase 2: API Infrastructure Foundation** - COMPLETED
  - Created ServiceHandler interface
  - Created validation and logging interceptors

- ✅ **Phase 3: Implement Auth Service** - COMPLETED
  - Created placeholder auth service with in-memory user storage
  - Implements Login and Signup RPCs

- ✅ **Phase 4: Integrate Connect RPC with Echo** - COMPLETED
  - Mounted Connect RPC mux at /connect/*
  - Registered AuthService

- 🔄 **Phase 5: Create Login Form Component** - IN PROGRESS
  - **User rolled back all changes** to forms/, components/, api/ directories
  - Plan document has been updated with correct approach
  - Need to implement from scratch following updated plan

## Critical References

1. **Implementation Plan** (UPDATED): `thoughts/shared/plans/2025-01-19-proto-connect-rpc-setup.md`
   - Contains complete phase-by-phase implementation
   - **Important**: Plan was significantly updated during this session

2. **Existing Form Component**: `components/form/form.templ` and `components/form/args.go`
   - Shows existing form patterns and structure
   - Needs ContentType field added for Connect RPC support

3. **Existing Input Component**: `components/input/input.templ:45-54`
   - Already has FormID pattern that auto-generates data-bind attributes
   - DO NOT add DataBind field - use existing FormID pattern

## Recent changes

**NOTE**: User has undone all implementation changes. Only the plan document was updated.

**Plan Updates Only**:
- `thoughts/shared/plans/2025-01-19-proto-connect-rpc-setup.md` - Complete rewrite of Phase 5 and architecture notes

**Infrastructure (still in place)**:
- `buf.yaml:1-13` - Buf workspace configuration
- `buf.gen.yaml:1-18` - Code generation config with managed mode
- `proto/com/datastarui/v1/forms/login/login.proto:1-65` - Login form messages
- `proto/com/datastarui/v1/auth/auth.proto:1-20` - Auth service definition
- `pkg/proto/com/datastarui/v1/forms/login/login.pb.go` - Generated (via buf generate)
- `pkg/proto/com/datastarui/v1/auth/authconnect/auth.connect.go` - Generated
- `justfile:2` - Added "buf generate" to build command
- `justfile:8-9` - Added "proto" command
- `justfile:33` - Added buf installation to install command
- `go.mod` - Added Connect RPC, protovalidate, protobuf dependencies
- `api/handler/handler.go:1-11` - ServiceHandler interface
- `api/interceptors/validation.go:1-40` - Proto validation interceptor
- `api/interceptors/logging.go:1-31` - Logging interceptor
- `api/services/auth/service.go:1-154` - Auth service implementation
- `main.go:32-33` - Added imports for handler and auth service
- `main.go:88-103` - Connect RPC mux setup and mounting

## Learnings

### 1. Avoid Runtime Reflection
- Connect RPC already handles JSON→Proto via protojson (built-in)
- Datastar serializes signals to JSON when using `contentType: 'json'`
- **Do not** create custom form-data-to-proto interceptors (would require reflection)
- Keep integration simple and type-safe

### 2. Datastar data-bind Syntax
From https://data-star.dev/examples/form_data/:
- `data-bind` does NOT use `$` prefix: `data-bind="login.email"`
- Expressions (data-show, data-text) DO use `$` prefix: `data-show="$login.email"`
- The existing FormID pattern in `components/input/input.templ:45-54` is correct

### 3. Input Component Pattern
- **DO NOT add DataBind field** to InputArgs
- Input component already has FormID pattern that works correctly:
  ```go
  if args.FormID != "" && args.Name != "" {
      signalName := strings.ReplaceAll(args.FormID, "-", "_")
      dataBindValue := signalName + "." + args.Name
      inputAttrs["data-bind"] = dataBindValue
  }
  ```
- Just pass `FormID: "login"` and it auto-generates `data-bind="login.email"`

### 4. Form Component Integration
- Existing Form component uses `contentType: 'form'` by default (form-urlencoded)
- Connect RPC needs JSON, so add `ContentType` field to FormArgs
- When `ContentType: "json"`, Datastar sends JSON directly to Connect RPC
- No custom conversion logic needed

### 5. buf.gen.yaml Managed Mode
- Must disable managed mode for protovalidate module:
  ```yaml
  managed:
    disable:
      - module: buf.build/bufbuild/protovalidate
  ```
- This prevents import path rewriting for validation imports

## Artifacts

**Created during session**:
- `buf.yaml` - Buf workspace configuration
- `buf.gen.yaml` - Code generation configuration
- `proto/com/datastarui/v1/forms/login/login.proto` - Login form messages with validation
- `proto/com/datastarui/v1/auth/auth.proto` - Auth service definition
- `pkg/proto/**/*.go` - Generated proto code (via buf generate)
- `api/handler/handler.go` - ServiceHandler interface
- `api/interceptors/validation.go` - Validation interceptor
- `api/interceptors/logging.go` - Logging interceptor
- `api/services/auth/service.go` - Auth service implementation
- `justfile` updates for proto generation

**Updated during session**:
- `thoughts/shared/plans/2025-01-19-proto-connect-rpc-setup.md` - Major updates to Phase 5 and architecture
- `main.go` - Connect RPC integration
- `go.mod` - Dependencies added

## Action Items & Next Steps

### Immediate Next Steps (Phase 5):

1. **Enhance Form Component** (components/form/)
   - Add `ContentType` field to `FormArgs` in `args.go`
   - Update `form.templ:48` to use configurable content type
   - Default to "form", support "json" for Connect RPC
   - See plan: `thoughts/shared/plans/2025-01-19-proto-connect-rpc-setup.md:966-1008`

2. **Create Login Form Component** (forms/login/)
   - Create `forms/login/login_form.templ`
   - Use Form component with `ContentType: "json"`
   - Use FormItem/FormLabel for structure
   - Use Input with `FormID` pattern (NOT DataBind)
   - See plan: `thoughts/shared/plans/2025-01-19-proto-connect-rpc-setup.md:795-905`

3. **Create Login Page** (pages/login/)
   - Create `pages/login/login_page.templ`
   - Import and use LoginForm component
   - See plan: `thoughts/shared/plans/2025-01-19-proto-connect-rpc-setup.md:907-934`

4. **Add Login Route** (main.go)
   - Add import for login page
   - Add GET /login route handler
   - See plan: `thoughts/shared/plans/2025-01-19-proto-connect-rpc-setup.md:936-958`

5. **Run Verification**
   - `templ generate && go build` - Verify compilation
   - Test with curl commands (see plan manual verification section)
   - Verify JSON content-type in network tab
   - Confirm no runtime reflection warnings in logs

### Follow-up (After Phase 5):
- Run all automated verification checks from plan
- Update plan file with completed checkmarks
- Test complete login flow: signup → login → token validation

## Other Notes

### Key Files to Reference:

**Existing Patterns**:
- `pages/components/formpage/form_page.templ:154-246` - Contact form example showing custom signals + Form component
- `components/input/input.templ:45-54` - FormID pattern implementation
- `components/form/form.templ:40-49` - Current submit handler generation

**Proto Definitions**:
- `proto/com/datastarui/v1/forms/login/login.proto` - LoginRequest, SignupRequest, User messages
- `proto/com/datastarui/v1/auth/auth.proto` - AuthService with Login/Signup RPCs

**Generated Code**:
- `pkg/proto/com/datastarui/v1/forms/login/login.pb.go` - Proto message structs
- `pkg/proto/com/datastarui/v1/auth/authconnect/auth.connect.go` - Connect service interfaces

### Common Commands:
- `buf generate` - Regenerate proto code
- `templ generate` - Regenerate templ templates
- `go build` - Verify compilation
- `just docker-tail app` - View app logs (check for errors)
- `just build` - Full build (runs buf generate, templ generate, go build, tailwind)

### Important Patterns to Follow:
1. **No DataBind field** - Use existing FormID pattern in Input component
2. **ContentType: "json"** - For Connect RPC forms
3. **utils.Signals()** - For initializing form signals
4. **FormItem/FormLabel/FormMessage** - For consistent form structure
5. **$fetching signal** - Datastar built-in for tracking POST/fetch state

### Architecture Decision:
The elegant integration uses:
- Form component with `ContentType: "json"` parameter
- Datastar serializes signals to JSON automatically
- Connect RPC unmarshals JSON via protojson (built-in)
- No custom reflection or conversion code needed
- Type-safe end-to-end
