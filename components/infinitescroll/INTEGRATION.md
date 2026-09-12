# InfiniteScroll Integration Guide

This document provides integration patterns and backend implementation guidance for the InfiniteScroll component.

## Overview

The InfiniteScroll component implements **Pattern A**: intersection-based progressive loading using Datastar's `data-on:intersect` attribute. When a sentinel element scrolls into view, it triggers a backend request that returns new content via Server-Sent Events (SSE).

## Key Principles

### 1. Host Stability
The Host container is **never remorphed**. Only the Items container and Loading indicators are patch targets. This ensures scroll position and container state remain stable during content updates.

### 2. Items = Stable History Only
The Items container holds **loaded content history**. Live, working, or latest content that shouldn't be part of the scrollable history should stay OUTSIDE the Items container as siblings.

**Sibling order is app-owned.** Place live content before or after Items based on your UX needs.

For example, in a chat application:
```go
// Live content as bottom anchor (after Items)
@infinitescroll.Host(...) {
  @infinitescroll.Items(...) {
    // Historical messages here (stable history)
  }
  @infinitescroll.Sentinel(...)
  
  <!-- Live content OUTSIDE Items -->
  <div id="chat-latest" class="p-4 border-t">
    <span>User is typing...</span>
  </div>
}
```

Or place live content at the top:
```go
// Live content at top (before Items)
@infinitescroll.Host(...) {
  <div id="chat-latest" class="p-4 border-b">
    <span>User is typing...</span>
  </div>
  
  @infinitescroll.Items(...) {
    // Historical messages here (stable history)
  }
  @infinitescroll.Sentinel(...)
}
```

Or beside Host:
```go
<div>
  <div id="live-status">Current status...</div>
  @infinitescroll.Host(...) {
    @infinitescroll.Items(...) { /* history */ }
  }
</div>
```

### 3. Chunk Patches Without View Transitions
Backend SSE responses that patch new content chunks should have View Transitions **OFF**. Use the datastar-go SDK's patch methods without view transition flags.

### 4. Loading State as Primary SoT
Loading indicators use **same-id DOM replace** on the sentinel ID. When loading starts, backend patches a Loading component onto the sentinel's ID. When content arrives, patch new items + a new sentinel (or remove if exhausted).

First paint: Sentinel only. Loading templ is for SSE patches from the backend.

## Backend Implementation

### Basic SSE Handler Pattern

```go
import (
	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"
	"github.com/coreycole/datastarui/components/infinitescroll"
)

func HandleLoadMore(c echo.Context) error {
	w := c.Response().Writer
	r := c.Request()
	sse := datastar.NewSSE(w, r)

	// Parse cursor/offset from query params
	cursor := c.QueryParam("cursor")
	offset := c.QueryParam("offset")

	// 1. Show loading indicator (replace sentinel with loading)
	loadingComponent := infinitescroll.Loading(infinitescroll.LoadingArgs{
		ID:        "my_list-sentinel-below",  // Same ID as sentinel
		Direction: infinitescroll.DirectionBelow,
	})
	sse.PatchElementTempl(loadingComponent)

	// 2. Fetch new items from database
	items, hasMore := fetchItems(cursor, offset, 20)

	// 3. Build items HTML
	var itemsHTML strings.Builder
	for i, item := range items {
		// Render each item with stable ID
		itemComponent := MyItemComponent(item, i)
		itemComponent.Render(context.Background(), &itemsHTML)
	}

	// 4. Append items to Items container
	sse.PatchElement(datastar.PatchElementOptions{
		Selector: "#my_list-items",
		Mode:     datastar.ModeAppend,
		Fragment: itemsHTML.String(),
	})

	// 5a. If more content available, add new sentinel
	if hasMore {
		newSentinel := infinitescroll.Sentinel(infinitescroll.SentinelArgs{
			ID:        "my_list-sentinel-below",
			Direction: infinitescroll.DirectionBelow,
			PatchExpr: "@get('/api/items/more?cursor=" + items[len(items)-1].ID + "')",
		})
		return sse.PatchElementTempl(newSentinel)
	}

	// 5b. If exhausted, remove the loading indicator (no new sentinel)
	return sse.PatchElement(datastar.PatchElementOptions{
		Selector: "#my_list-sentinel-below",
		Mode:     datastar.ModeRemove,
	})
}
```

### Above Edge (Inverse/Scrollback) Pattern

For loading earlier content (e.g., scrolling up in a chat):

```go
func HandleLoadBefore(c echo.Context) error {
	w := c.Response().Writer
	r := c.Request()
	sse := datastar.NewSSE(w, r)

	cursor := c.QueryParam("cursor")

	// 1. Show loading indicator
	loadingComponent := infinitescroll.Loading(infinitescroll.LoadingArgs{
		ID:        "my_list-sentinel-above",
		Direction: infinitescroll.DirectionAbove,
	})
	sse.PatchElementTempl(loadingComponent)

	// 2. Fetch earlier items
	items, hasMore := fetchItemsBefore(cursor, 20)

	// 3. Build items HTML
	var itemsHTML strings.Builder
	for i, item := range items {
		itemComponent := MyItemComponent(item, i)
		itemComponent.Render(context.Background(), &itemsHTML)
	}

	// 4. PREPEND items to Items container (note: Prepend mode)
	sse.PatchElement(datastar.PatchElementOptions{
		Selector: "#my_list-items",
		Mode:     datastar.ModePrepend,
		Fragment: itemsHTML.String(),
	})

	// 5. Add new sentinel or remove if exhausted
	if hasMore {
		newSentinel := infinitescroll.Sentinel(infinitescroll.SentinelArgs{
			ID:        "my_list-sentinel-above",
			Direction: infinitescroll.DirectionAbove,
			PatchExpr: "@get('/api/items/before?cursor=" + items[0].ID + "')",
		})
		return sse.PatchElementTempl(newSentinel)
	}

	return sse.PatchElement(datastar.PatchElementOptions{
		Selector: "#my_list-sentinel-above",
		Mode:     datastar.ModeRemove,
	})
}
```

## Frontend Integration

### Happy Path: Single Root Call

```go
@infinitescroll.InfiniteScroll(infinitescroll.InfiniteScrollArgs{
	ID:             "my_list",
	PatchBelowExpr: "@get('/api/items/more')",
}) {
	// First page of items
	for i, item := range initialItems {
		@MyItemComponent(item, i)
	}
}
```

The component automatically creates:
- Host container: `my_list-host`
- Items container: `my_list-items`
- Below sentinel: `my_list-sentinel-below`

**Note:** `LoadingAboveID` and `LoadingBelowID` default to the corresponding sentinel IDs for same-id DOM replace. Leave them empty (recommended) or explicitly set them equal to the sentinel ID. Setting a different ID breaks the same-id SoT pattern. No separate loading node exists on first paint.

### Advanced Composition with Custom IDs

Override MorphMap IDs when you need custom targeting:

```go
@infinitescroll.InfiniteScroll(infinitescroll.InfiniteScrollArgs{
	ID:             "my_list",
	PatchBelowExpr: "@get('/api/items/more')",
	HostID:         "custom-scroll-host",
	ItemsID:        "custom-items-container",
	SentinelBelowID: "custom-sentinel",
	LoadingBelowID:  "custom-loading",
}) {
	// Initial items
}
```

### Bidirectional Scrolling

Enable both above and below edges:

```go
@infinitescroll.InfiniteScroll(infinitescroll.InfiniteScrollArgs{
	ID:             "chat_history",
	PatchAboveExpr: "@get('/api/chat/before')",
	PatchBelowExpr: "@get('/api/chat/after')",
}) {
	// Middle page of chat messages
	for _, msg := range messages {
		@ChatMessage(msg)
	}
}
```

## Cursor Management

### Query Parameter Approach

Pass cursor state in URL query parameters:

```go
// Sentinel triggers
PatchExpr: "@get('/api/items/more?cursor=" + lastItemID + "&limit=20')"

// Backend extracts
cursor := c.QueryParam("cursor")
limit := c.QueryParam("limit")
```

### Form Data Approach

For POST requests with complex state:

```go
// Sentinel triggers
PatchExpr: "@post('/api/items/more', contentType:'form')"

// Backend extracts from form
type LoadMoreRequest struct {
	Cursor string `form:"cursor"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
}

var req LoadMoreRequest
if err := c.Bind(&req); err != nil {
	return err
}
```

## Common Patterns

### Pattern: Stable Item IDs

Always use stable, unique IDs for each item to enable efficient DOM morphing:

```go
@MyItemComponent(item, index) {
	<div id={ fmt.Sprintf("item-%s", item.ID) } class="...">
		// Item content
	</div>
}
```

### Pattern: Loading Debounce

If rapid scrolling causes multiple requests, add debouncing on the backend:

```go
func HandleLoadMore(c echo.Context) error {
	// Check if request is duplicate within time window
	if isDuplicateRequest(c) {
		return c.NoContent(http.StatusNoContent)
	}
	
	// ... rest of handler
}
```

### Pattern: Error Handling

Show user-friendly error messages when loading fails:

```go
func HandleLoadMore(c echo.Context) error {
	sse := datastar.NewSSE(w, r)
	
	items, err := fetchItems(cursor)
	if err != nil {
		// Show error message in place of loading indicator
		errorComponent := ErrorMessage("Failed to load more items")
		return sse.PatchElementTempl(errorComponent)
	}
	
	// ... success path
}
```

### Pattern: Empty States

Handle the case when there are no items to load:

```go
if len(items) == 0 {
	emptyComponent := EmptyState("No more items to load")
	return sse.PatchElement(datastar.PatchElementOptions{
		Selector: "#my_list-sentinel-below",
		Mode:     datastar.ModeReplace,
		Fragment: renderToString(emptyComponent),
	})
}
```

## Performance Considerations

1. **Batch Size**: Load 10-50 items per request depending on item complexity
2. **Preload**: Consider preloading the next page when user is 80% through current content
3. **Image Lazy Loading**: Use native `loading="lazy"` on images within items
4. **Memory Management**: Consider removing items far from viewport if memory becomes an issue

**Pattern B (out of scope):** DSUI ships Pattern A only. Pattern B tape virtualization (e.g. largediff pixel tape, ~2 live files of 200k lines, POST 204 + warm SSE, fat morph `#app`) is app-owned and not implemented in this component.

## Testing

### Manual Testing Checklist

- [ ] Initial page loads correctly
- [ ] Scrolling to sentinel triggers loading indicator
- [ ] New items append/prepend correctly
- [ ] Sentinel moves to new position after load
- [ ] Loading indicator disappears after content loads
- [ ] Exhausted state handled (sentinel removed)
- [ ] Scroll position maintained during loads
- [ ] Multiple rapid scrolls don't cause issues
- [ ] Network errors handled gracefully
- [ ] Empty states display correctly

### E2E Testing

Use Go Story tests (see `components/infinitescroll/infinitescroll_component_e2e_test.go`) to verify:
- Sentinel element has correct `data-on:intersect` attribute
- Host container remains stable (ID unchanged across patches)
- Items container receives appended/prepended content
- Console remains clean (no JavaScript errors)
