---
date: 2025-10-25T23:44:08+0000
researcher: Claude (claude-sonnet-4-5)
git_commit: 597740fade14ec2fe67f4c37ccd051342cb20801
branch: main
repository: datastarui
topic: "How does Tailwind CSS build and caching work in datastarui?"
tags: [research, codebase, tailwind, css, http-caching, cache-busting]
status: complete
last_updated: 2025-10-25
last_updated_by: Claude
---

# Research: Tailwind CSS Build Process and HTTP Caching in DatastarUI

**Date**: 2025-10-25T23:44:08+0000
**Researcher**: Claude (claude-sonnet-4-5)
**Git Commit**: 597740fade14ec2fe67f4c37ccd051342cb20801
**Branch**: main
**Repository**: datastarui

## Research Question

How does Tailwind CSS build work in datastarui? Why are there 304 Not Modified responses to built files with unique hashed names, suggesting the browser cache is not being invalidated properly?

## Summary

DatastarUI implements a **build-time CSS fingerprinting system** that generates SHA-256 hashed filenames (e.g., `out.bfd99bc8.css`) for cache busting. However, the system has a critical flaw: **HTTP 304 responses are still sent even for new hashed filenames** because:

1. The Echo server uses Go's default `http.FileServer` which implements ETag/Last-Modified HTTP caching
2. The CSS build process copies files, preserving their modification timestamps (ModTime)
3. Go's `http.FileServer` uses ModTime to generate ETags and check `If-Modified-Since` headers
4. Even with a new filename, if the ModTime is the same, the server sends 304 Not Modified
5. The browser serves stale content from cache despite the filename change

**The Solution**: Set `Cache-Control: no-store` headers for static assets during development (documented in `.cursor/rules/templ.mdc:1237-1258` but not implemented).

## Detailed Findings

### 1. Tailwind CSS Build Configuration

**Build Command** (`justfile:10-20`):

```bash
build-tailwind:
    @echo "🎨 Building Tailwind CSS..."
    @pnpm exec tailwindcss -i static/css/index.css -o static/css/out.css \
        --content "./components/**/*" \
        --content "./pages/**/*" \
        --content "./layouts/**/*"
    @if [ -f static/css/out.css ]; then \
        echo "📝 Generating CSS hash..."; \
        HASH=$(sha256sum static/css/out.css | cut -d' ' -f1 | head -c8); \
        echo "🔖 Hash: $HASH"; \
        rm -f static/css/out.*.css; \
        cp static/css/out.css static/css/out.$HASH.css; \
        echo "✅ Created static/css/out.$HASH.css"; \
    fi
```

**How it works**:
- Tailwind CLI compiles `static/css/index.css` → `static/css/out.css`
- Generates SHA-256 hash of compiled CSS file (first 8 characters)
- Removes any previously hashed files matching `out.*.css`
- **Copies** (not moves) the compiled CSS to `static/css/out.$HASH.css`

**Key files**:
- Input: `static/css/index.css`
- Temp output: `static/css/out.css`
- Final output: `static/css/out.bfd99bc8.css` (hashed)

**Configuration**:
- `tailwind.config.js` - Tailwind CSS configuration
- `package.json` - Contains Tailwind dependencies (`@tailwindcss/cli@^4.1.12`)
- `Dockerfile:22-32` - Installs pnpm and Tailwind dependencies

### 2. CSS File Serving and Hashing Mechanism

**Runtime Discovery** (`layouts/root.templ:32-40`):

```go
templ Root(args RootArgs) {
	{{
		// Find the hashed CSS file
		cssFile := "/css/out.css" // fallback
		if files, err := filepath.Glob("static/css/out.*.css"); err == nil && len(files) > 0 {
			// Extract just the filename from the full path
			fileName := filepath.Base(files[0])
			cssFile = "/css/" + fileName
		}
	}}
```

The template dynamically discovers the hashed CSS file at runtime using `filepath.Glob("static/css/out.*.css")` and constructs the URL path.

**HTML Reference** (`layouts/root.templ:62`):

```html
<link rel="stylesheet" href={ cssFile }/>
```

Results in: `<link rel="stylesheet" href="/css/out.bfd99bc8.css"/>`

**Static File Serving** (`main.go:218`):

```go
// Serve static files
e.Static("/", "static/")
```

Echo's Static middleware maps:
- URL: `/css/out.bfd99bc8.css` → Filesystem: `static/css/out.bfd99bc8.css`

### 3. HTTP Caching Behavior (The Problem)

**Default Caching** (`main.go:218`):

The single line `e.Static("/", "static/")` uses Echo's default static file serving, which delegates to Go's standard library `http.FileServer`. This automatically implements HTTP caching via:
- **Last-Modified** header (file's ModTime)
- **ETag** header (generated from ModTime and file size)
- **304 Not Modified** responses when `If-Modified-Since` or `If-None-Match` matches

**Why 304 Responses Occur**:

```
1. Build process: cp static/css/out.css static/css/out.bfd99bc8.css
   → Both files have the SAME modification timestamp

2. Browser requests: /css/out.bfd99bc8.css
   → Server reads file ModTime
   → Generates ETag from ModTime + size
   → Sends Last-Modified: <timestamp>

3. Browser caches file with ETag and Last-Modified

4. Next request (even with new hash):
   → Browser sends: If-Modified-Since: <timestamp>
   → Server checks file ModTime
   → Same ModTime → Sends 304 Not Modified
   → Browser serves stale content from cache
```

**The Root Cause**: File copy (`cp`) preserves the original modification timestamp, so even though the filename changes, the ModTime doesn't.

**No Cache Control Headers**:
- Application doesn't set `Cache-Control` headers
- No development mode detection (`Config` struct in `main.go:42-45` only has `DatastarInspectorEnabled`)
- No middleware to disable caching in development

### 4. Development Environment Build Triggers

**Docker Setup** (`docker-compose.yml:18`):

```yaml
command: sh -c "just watch"
```

Runs `just watch` → starts Air live reload

**Air Configuration** (`.air.toml:8,16-17`):

```toml
[build]
  cmd = "just"  # Runs default justfile recipe
  include_dir = ["api", "forms", "pkg", "components", "pages", "layouts", "static"]
  include_ext = ["go", "templ", "html", "css", "tsx", "ts", "jsx", "js"]
  exclude_file = ["static/css/out*"]  # Prevent rebuild loops
```

**Build Pipeline** (`justfile:1-4`):

```bash
build:
  @templ generate
  @go build -o datastarui main.go
  @just build-tailwind
```

**Flow**:
1. Developer edits `.templ` file with Tailwind classes
2. Air detects file change
3. Runs `just` → `templ generate` → `go build` → `just build-tailwind`
4. CSS is rebuilt and rehashed
5. Server restarts
6. **Problem**: Browser receives 304 for new hashed file due to identical ModTime

**Alternative: Manual Watch Mode** (`justfile:22-23`):

```bash
tailwind:
  @pnpm exec tailwindcss -i static/css/index.css -o static/css/out.css \
      --watch \
      --content "./components/**/*" \
      --content "./pages/**/*" \
      --content "./layouts/**/*"
```

This is **NOT** started automatically by Docker. Developer must run `just tailwind` manually for continuous CSS rebuilds.

## Code References

- `justfile:10-20` - Tailwind build and hash generation
- `justfile:1-4` - Full build pipeline (templ + go + css)
- `layouts/root.templ:32-40` - Dynamic CSS file discovery using filepath.Glob
- `layouts/root.templ:62` - HTML link tag for hashed CSS file
- `main.go:218` - Static file serving configuration
- `main.go:42-45` - Config struct (no dev mode flag)
- `docker-compose.yml:18` - Container startup command
- `.air.toml:8,11,16-17` - Air live reload configuration
- `Dockerfile:22-32` - pnpm and Tailwind installation

## Architecture Insights

### Design Patterns

1. **Build-Time Fingerprinting**:
   - Hash generated during build, not at server startup
   - SHA-256 ensures content-based cache invalidation
   - 8-character hash provides uniqueness while keeping URLs short

2. **Runtime Discovery**:
   - Template uses `filepath.Glob()` to find hashed file dynamically
   - No hardcoded hash values in code
   - Fallback mechanism (`/css/out.css`) ensures graceful degradation

3. **Cache Busting Strategy**:
   - Intended: Every CSS change produces new hash → new URL → browser fetches fresh content
   - Reality: HTTP 304 responses bypass the cache busting due to ModTime preservation

### Key Conventions

- **Glob Pattern Limitation**: `filepath.Glob("static/css/out.*.css")` returns all matching files; uses `files[0]`. The build process cleans up old files to prevent multiple matches.
- **Air Exclusion**: `exclude_file = ["static/css/out*"]` prevents infinite rebuild loops when CSS output changes
- **No Separate Tailwind Watch**: Despite documentation claiming Tailwind runs in watch mode automatically, it doesn't - CSS is rebuilt on every Air reload

### The Documentation-Implementation Gap

**Documented Solution** (`.cursor/rules/templ.mdc:1237-1258`):

```go
var dev = true

func disableCacheInDevMode(next http.Handler) http.Handler {
	if !dev {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
```

**Status**: This middleware pattern is documented in project rules but **NOT implemented** in the codebase.

## Related Research

- None found in `thoughts/shared/research/` - this is the first research document for datastarui.

## Open Questions

1. **Should development mode detection be added?**
   - Add environment variable to `Config` struct (`main.go:42-45`)
   - Example: `DevMode bool envconfig:"DEV_MODE" default:"true"`

2. **Where should cache-disabling middleware be applied?**
   - Option A: Wrap `e.Static()` with custom middleware
   - Option B: Use Echo's middleware chain before static file handler
   - Option C: Configure `middleware.StaticConfig` with custom `Filesystem` handler

3. **Should production use different cache headers?**
   - Development: `Cache-Control: no-store` (no caching)
   - Production: `Cache-Control: public, max-age=31536000, immutable` (aggressive caching for hashed files)

4. **Should the build process touch files to update ModTime?**
   - Alternative to cache headers: Use `touch` command after copying to update ModTime
   - Pros: Simpler, no middleware needed
   - Cons: Doesn't address root cause, still relies on ModTime-based caching

## Recommended Solutions

### Solution 1: Implement Cache-Control Middleware (Recommended)

```go
// In main.go
func (s *Server) setupMiddleware() {
	if s.config.DevMode {
		s.echo.Use(disableCacheInDevMode)
	}
}

func disableCacheInDevMode(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-store")
		return next(c)
	}
}
```

**Pros**:
- Addresses root cause of caching issue
- Documented best practice
- Works for all static assets
- Easy to toggle based on environment

**Cons**:
- Requires adding `DevMode` to config

### Solution 2: Update ModTime After Copy

```bash
# In justfile:18
cp static/css/out.css static/css/out.$HASH.css
touch static/css/out.$HASH.css  # Update modification time
```

**Pros**:
- Simple one-line fix
- No code changes needed

**Cons**:
- Band-aid solution
- Doesn't prevent 304s between subsequent requests
- Doesn't address production caching strategy

### Solution 3: Use Aggressive Production Caching

For production, leverage the hash-based cache busting by setting long-lived cache headers:

```go
func productionCacheHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Only for hashed static files
		if strings.Contains(c.Request().URL.Path, "/css/out.") {
			c.Response().Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		return next(c)
	}
}
```

**Pros**:
- Optimal performance in production
- Safe because hash changes with content
- Reduces server load

**Cons**:
- Requires environment detection
- More complex middleware logic

## Conclusion

The DatastarUI Tailwind build system implements content-based cache busting through SHA-256 filename hashing, but HTTP caching defeats this mechanism during development. The server sends 304 Not Modified responses because:

1. File copying preserves modification timestamps
2. Go's `http.FileServer` uses ModTime for ETags and Last-Modified headers
3. No cache-control headers are set to disable browser caching

The solution is documented (`.cursor/rules/templ.mdc`) but not implemented. Adding a simple middleware to set `Cache-Control: no-store` during development would resolve the issue.
