---
date: 2025-10-19T15:52:33-07:00
researcher: Claude
git_commit: 547b968e967ccee234c4e52c2e6a5c8d10bd0967
branch: develop
repository: monorepo
topic: "Connect RPC Setup with Proto Validate Integration"
tags: [research, codebase, connect-rpc, protobuf, proto-validate, api, code-generation]
status: complete
last_updated: 2025-01-10
last_updated_by: Claude
---

# Research: Connect RPC Setup with Proto Validate Integration

**Date**: 2025-10-19T15:52:33-07:00
**Researcher**: Claude
**Git Commit**: 547b968e967ccee234c4e52c2e6a5c8d10bd0967
**Branch**: develop
**Repository**: monorepo

## Research Question

How does this project set up Connect RPC in the api directory and how does that connect to the proto directory with proto validate?

## Summary

The monorepo uses Connect RPC (connectrpc.com) as a modern RPC framework built on HTTP/2, providing type-safe communication between Go backend services and TypeScript frontends. Protocol Buffers in `/proto/` define service contracts with embedded validation rules using `buf.validate`. The Buf toolchain generates Go code in `/pkg/proto/` which the API services in `/api/internal/` consume. The architecture features comprehensive middleware chains for authentication, authorization, validation, tracing, and error handling, supporting both internal Connect RPC services and external REST/gRPC APIs via Vanguard transcoding.

## Detailed Findings

### Connect RPC Architecture Overview

The API server provides two distinct API surfaces:
1. **Internal Connect RPC API** at `/connect` - For internal clients using Connect protocol
2. **External REST/gRPC API** at `/api` - For external clients via Vanguard transcoder

Key architectural components ([api/cmd/cmd.go:153-192]()):
- `connectMux *http.ServeMux` - Multiplexer for internal Connect RPC services
- `vanguardHandler http.HandlerFunc` - External API transcoder supporting REST/gRPC
- `streamMux *http.ServeMux` - Dedicated HTTP streaming endpoints
- `router router.Router` - Chi-based main router for HTTP middleware

### Proto Organization and Structure

#### Directory Hierarchy ([proto/com/chestnutfi/]()):
```
proto/com/chestnutfi/
├── actions/v1/          - Action system definitions
├── api/v1/              - Public API contracts
├── authz/v1/            - Authorization annotations & rules
├── entities/v1/         - Producer/entity management
├── workflows/v1-v3/     - Workflow engine (multiple versions)
└── [40+ domain packages]
```

#### Proto Validation Patterns

**Field-Level Validation** ([proto/com/chestnutfi/apikeys/v1/keys.proto:24-33]()):
```protobuf
message CreateAPIKeyRequest {
  string name = 1 [
    (buf.validate.field).required = true,
    (buf.validate.field).string.min_len = 1
  ];
  string version = 2 [
    (buf.validate.field).required = true,
    (buf.validate.field).string.len = 10
  ];
}
```

**Message-Level CEL Validation** ([proto/com/chestnutfi/entities/v1/entities.proto:726-730]()):
```protobuf
option (buf.validate.message).cel = {
  id: "request.business"
  message: "business required if agency"
  expression: "has(this.business) && this.business.npn != '' && this.business.ein != '' || this.type != 1"
};
```

**Authorization Annotations** ([proto/com/chestnutfi/entities/v1/entities.proto:32-46]()):
```protobuf
rpc UpdateEntity(UpdateEntityRequest) returns (UpdateEntityResponse) {
  option (com.chestnutfi.authz.v1.method) = {
    object: "entity"
    action: "write"
    all: true
  };
}
```

### Code Generation Pipeline

#### Configuration ([buf.gen.yaml]()):
```yaml
plugins:
  # Go Protocol Buffers
  - remote: buf.build/protocolbuffers/go:v1.36.6
    out: pkg/proto
    opt: paths=source_relative

  # Connect RPC for Go
  - remote: buf.build/connectrpc/go:v1.18.1
    out: pkg/proto
    opt: paths=source_relative

  # TypeScript with Connect Query
  - remote: buf.build/bufbuild/es:v2.2.5
    out: frontend/packages/proto
  - remote: buf.build/connectrpc/query-es:v2.0.1
    out: frontend/packages/proto
```

#### Generated Code Structure:
- **Go**: `/pkg/proto/com/chestnutfi/{service}/{version}/`
  - `*.pb.go` - Protobuf message types
  - `*connect/*.connect.go` - Connect RPC service interfaces
- **TypeScript**: `/frontend/packages/proto/`
  - `*_pb.ts` - Message types
  - `*_connectquery.ts` - React Query hooks

### Service Implementation Pattern

#### Service Handler Structure ([api/internal/actions/connect.go:15-43]()):
```go
type ServiceHandler struct {
    actionsv1connect.UnimplementedActionsServiceHandler
}

func (s ServiceHandler) Handler(decider authz.Decider) (string, http.Handler) {
    return actionsv1connect.NewActionsServiceHandler(s, connect.WithInterceptors(
        interceptors.OtelInterceptor(),         // Tracing
        interceptors.TenantIDInterceptor(),     // Multi-tenancy
        interceptors.UserIDInterceptor(),       // User context
        interceptors.ValidationInterceptor(),   // Proto validation
        interceptors.SentryInterceptor(),       // Error tracking
        interceptors.RoleInterceptor(),         // Role extraction
        interceptors.RoleIDInterceptor(),       // Role ID extraction
        decider.Interceptor(),                  // Authorization
    ))
}
```

#### RPC Method Implementation ([api/internal/apikeys/service.go:53-102]()):
```go
func (s ServiceHandler) CreateAPIKey(
    ctx context.Context,
    req *connect.Request[v1.CreateAPIKeyRequest],
) (*connect.Response[v1.CreateAPIKeyResponse], error) {
    tenant := interceptors.TenantID(ctx)  // Extract from context

    err := s.repo.InTx(ctx, func(tx Repository) error {
        key, err := tx.CreateAPIKey(ctx, db.CreateAPIKeyParams{
            Name:     req.Msg.Name,      // Access proto fields
            Version:  req.Msg.Version,
            TenantID: tenant,
        })
        // Audit logging...
        return err
    })

    return connect.NewResponse(&v1.CreateAPIKeyResponse{
        ApiKey: encodeKey(key, latestVersion),  // DB to proto
    }), nil
}
```

### Validation Enforcement

#### Runtime Validation Interceptor ([api/pkg/interceptors/validation.go:12-34]()):
```go
func ValidationInterceptor() connect.UnaryInterceptorFunc {
    return func(next connect.UnaryFunc) connect.UnaryFunc {
        return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
            err := protovalidate.Validate(req.Any().(proto.Message))
            if err != nil {
                connectErr := connect.NewError(connect.CodeInvalidArgument, err)
                if validationErr := new(protovalidate.ValidationError); errors.As(err, &validationErr) {
                    if detail, err := connect.NewErrorDetail(validationErr.ToProto()); err == nil {
                        connectErr.AddDetail(detail)  // Attach validation details
                    }
                }
                return nil, connectErr
            }
            return next(ctx, req)
        }
    }
}
```

This interceptor:
1. Validates incoming requests against proto validation rules
2. Returns `CodeInvalidArgument` errors with detailed field violations
3. Prevents invalid requests from reaching service handlers

### Authorization System

#### Authorization Decider ([api/pkg/authz/decider.go:64-79]()):
- LRU cache for role permissions (5 min TTL, 20 entries)
- Singleflight for concurrent permission fetches
- Permission checks at object/field/action granularity

#### Authorization Interceptor ([api/pkg/authz/interceptor.go:150-227]()):
1. **Request Phase**: Check method-level permissions
2. **Response Phase**: Clear unauthorized fields using protorange

### Server Initialization and Routing

#### Main Server Start ([api/cmd/cmd.go:227-283]()):
```go
func (ms *Server) Start() {
    // Initialize tracing
    tracingShutdown := tracing.InitTracer(...)
    defer tracingShutdown(context.Background())

    // Mount services
    ms.router.Mount("/connect", http.StripPrefix("/connect", ms.connectMux))
    ms.router.Mount("/api", http.StripPrefix("/api", ms.vanguardHandler))

    // Start HTTP/2 server with h2c
    h2s := &http2.Server{}
    srv := &http.Server{
        Handler: h2c.NewHandler(ms.router, h2s),
        Addr:    fmt.Sprintf(":%d", ms.port),
    }

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM)
    <-quit
    srv.Shutdown(context.WithTimeout(context.Background(), 5*time.Second))
}
```

### External API with Vanguard

Vanguard transcodes between Connect RPC, gRPC, and REST protocols ([api/cmd/cmd.go:1083-1145]()):

```go
transcoder, err := vanguard.NewTranscoder(
    []*vanguard.Service{
        vanguard.NewService(policyapi.NewService(...).Handler()),
        vanguard.NewService(jobsapi.NewService(...).Handler()),
        // ... more services
    },
    vanguard.WithUnknownHandler(streamHandler),
)
ms.vanguardHandler = transcoder.ServeHTTP
```

External services use different interceptors:
- `ExternalAPIAuth` - API key authentication
- `GuardrailInterceptor` - Rate limiting
- No session-based auth required

## Code References

### Core Files
- `api/main.go:8-13` - Entry point
- `api/cmd/cmd.go:194-223` - Server constructor with builder pattern
- `api/cmd/cmd.go:227-283` - Server start and graceful shutdown
- `api/cmd/server_dependencies.go:38-65` - Router initialization with middleware
- `api/pkg/handler/handler.go:14-17` - ServiceHandler interface definition

### Interceptors
- `api/pkg/interceptors/interceptors.go:44-84` - Base interceptor builder pattern
- `api/pkg/interceptors/validation.go:12-34` - Proto validation enforcement
- `api/pkg/interceptors/tenant_id.go:13-15` - Multi-tenancy support
- `api/pkg/interceptors/otel_interceptor.go:61-82` - OpenTelemetry tracing
- `api/pkg/authz/interceptor.go:150-227` - Authorization with field-level permissions

### Service Examples
- `api/internal/actions/connect.go:15-43` - Service handler registration
- `api/internal/apikeys/service.go:53-102` - Database transaction with proto mapping
- `api/internal/apikeys/encode.go:13-23` - Database to proto encoding

### Proto Definitions
- `proto/com/chestnutfi/authz/v1/authz.proto:58-66` - Custom authorization options
- `proto/com/chestnutfi/type/v1/filter.proto:144-175` - Generic filtering system
- `proto/com/chestnutfi/type/v1/page.proto:7-27` - Pagination patterns

### Configuration
- `buf.yaml` - Buf workspace with dependencies
- `buf.gen.yaml` - Code generation plugins and options
- `justfile:237-238` - Build commands for proto generation

## Architecture Insights

### Key Design Patterns

1. **Service Handler Interface**: All services implement a common interface for consistent registration and middleware application

2. **Interceptor Chaining**: Cross-cutting concerns (auth, validation, tracing) are cleanly separated into composable interceptors

3. **Type-Safe Generics**: Connect RPC uses Go generics (`connect.Request[T]`, `connect.Response[T]`) for compile-time type safety

4. **Builder Pattern**: Server initialization uses method chaining for clear dependency setup

5. **Repository Pattern**: Services use repository interfaces for database operations, enabling easy testing with mocks

### Multi-Protocol Support

The system supports multiple protocols through different paths:
- `/connect/*` - Native Connect RPC protocol (binary protobuf over HTTP/2)
- `/api/*` - REST/gRPC via Vanguard transcoding
- Streaming endpoints for specialized use cases

### Validation Strategy

Three-layer validation approach:
1. **Proto annotations** - Declarative validation rules in `.proto` files
2. **Runtime interceptor** - Automatic validation before handler execution
3. **Manual validation** - Additional business logic validation in handlers

### Error Handling

Structured error handling with Connect error codes:
- `CodeInvalidArgument` - Validation failures
- `CodeInternal` - Server errors
- `CodeNotFound` - Resource not found
- `CodeUnauthenticated` - Auth failures
- `CodePermissionDenied` - Authz failures

Errors include detailed field violations for client-side handling.

## Historical Context (from thoughts/)

The codebase shows evolution from earlier RPC systems to Connect RPC, with versioned services (v1, v2, v3) indicating API evolution. The use of Vanguard for external APIs suggests a need to support legacy clients while modernizing internally with Connect RPC.

## Related Research

- API documentation in `api/README.md` and `api/CLAUDE.md`
- Proto guidelines in `proto/CLAUDE.md`
- Frontend integration patterns in `frontend/apps/web/utils/api/api.ts`

## Open Questions

1. **Performance Tuning**: Are there specific performance optimizations for high-throughput services?
2. **Service Discovery**: How do services discover each other in production deployments?
3. **Circuit Breaking**: Is there resilience patterns for downstream service failures?
4. **Metrics Collection**: Beyond tracing, what metrics are collected for SLI/SLO monitoring?
5. **Schema Evolution**: What's the strategy for backward-compatible proto changes?