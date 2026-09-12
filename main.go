package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	l "github.com/coreycole/datastarui/layouts"
	p "github.com/coreycole/datastarui/pages"
	"github.com/coreycole/datastarui/pages/components/breadcrumbpage"
	"github.com/coreycole/datastarui/pages/components/buttonpage"
	"github.com/coreycole/datastarui/pages/components/calendarpage"
	"github.com/coreycole/datastarui/pages/components/cardpage"
	"github.com/coreycole/datastarui/pages/components/checkboxpage"
	"github.com/coreycole/datastarui/pages/components/componentspage"
	"github.com/coreycole/datastarui/pages/components/datepickerpage"
	"github.com/coreycole/datastarui/pages/components/dialogpage"
	"github.com/coreycole/datastarui/pages/components/dropdownpage"
	"github.com/coreycole/datastarui/pages/components/formpage"
	"github.com/coreycole/datastarui/pages/components/infinitescrollpage"
	"github.com/coreycole/datastarui/pages/components/popoverpage"
	"github.com/coreycole/datastarui/pages/components/selectpage"
	"github.com/coreycole/datastarui/pages/components/sheetpage"
	"github.com/coreycole/datastarui/pages/components/sidebarpage"
	"github.com/coreycole/datastarui/pages/components/tabspage"
	"github.com/coreycole/datastarui/pages/components/toastpage"
	"github.com/coreycole/datastarui/pages/components/tooltippage"
	"github.com/coreycole/datastarui/pages/login"
	"github.com/coreycole/datastarui/utils"

	"github.com/coreycole/datastarui/api/handler"
	apimw "github.com/coreycole/datastarui/api/middleware"
	"github.com/coreycole/datastarui/api/services/auth"
	loginform "github.com/coreycole/datastarui/forms/login"
	authconnect "github.com/coreycole/datastarui/pkg/proto/com/datastarui/v1/auth/authconnect"

	"github.com/starfederation/datastar-go/datastar"
	"github.com/coreycole/datastarui/components/infinitescroll"
)

// Config holds environment configuration
type Config struct {
	Port                     string `envconfig:"PORT" default:"4242"`
	DatastarInspectorEnabled bool   `envconfig:"DATASTAR_INSPECTOR_ENABLED" default:"false"`
	DatastarProAvailable     bool
}

// Helper function to create RootArgs for component pages
func componentRootArgs(path string, cfg Config) l.RootArgs {
	return l.RootArgs{
		CurrentPage:          "components",
		CurrentPath:          path,
		InspectorEnabled:     cfg.DatastarInspectorEnabled,
		DatastarProAvailable: cfg.DatastarProAvailable,
	}
}

// handleInfiniteScrollMore demonstrates Pattern A infinite scroll with SSE patches.
// Backend patches: (1) Loading replace on sentinel ID, (2) append items to Items,
// (3) new sentinel or remove when exhausted. View Transitions OFF for chunk patches.
func handleInfiniteScrollMore(c echo.Context) error {
	w := c.Response().Writer
	r := c.Request()

	// Parse cursor for pagination (0 = page 1, 1 = page 2, etc.)
	cursor := c.QueryParam("cursor")
	if cursor == "" {
		cursor = "0" // Default to page 1
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	sse := datastar.NewSSE(w, r)

	// 1. Show loading indicator (same-id DOM replace on sentinel) with VT OFF
	loadingComponent := infinitescroll.Loading(infinitescroll.LoadingArgs{
		ID:        "diff_viewer-sentinel-below",
		Direction: infinitescroll.DirectionBelow,
	})
	sse.PatchElementTempl(loadingComponent, datastar.WithoutViewTransitions())

	// 2. Add delay for visible loading (250-400ms)
	time.Sleep(300 * time.Millisecond)

	// 3. Fetch more files based on cursor
	var moreFiles []infinitescrollpage.DiffFile
	var nextCursor string
	var hasMore bool

	switch cursor {
	case "0": // Page 1
		moreFiles = []infinitescrollpage.DiffFile{
			{
				Path:     "components/infinitescroll/args.go",
				OldLines: 0,
				NewLines: 77,
				Hunks: []infinitescrollpage.DiffHunk{
					{
						Header: "@@ -0,0 +1,77 @@",
						Lines: []infinitescrollpage.DiffLine{
							{Type: "add", Content: "package infinitescroll", Number: 1},
							{Type: "add", Content: "", Number: 2},
							{Type: "add", Content: "import \"github.com/a-h/templ\"", Number: 3},
							{Type: "add", Content: "", Number: 4},
							{Type: "add", Content: "// InfiniteScrollArgs defines the properties for the InfiniteScroll container", Number: 5},
							{Type: "add", Content: "type InfiniteScrollArgs struct {", Number: 6},
							{Type: "add", Content: "\tID             string // Required: snake_case identifier", Number: 7},
							{Type: "add", Content: "\tPatchAboveExpr string // Optional: expression for loading above/scrollback", Number: 8},
							{Type: "add", Content: "\tPatchBelowExpr string // Optional: expression for loading below/forward", Number: 9},
							{Type: "add", Content: "\t// MorphMap overrides - empty = derived from ID", Number: 10},
							{Type: "add", Content: "\tHostID            string", Number: 11},
							{Type: "add", Content: "\tItemsID           string", Number: 12},
							{Type: "add", Content: "\tSentinelAboveID   string", Number: 13},
							{Type: "add", Content: "\tSentinelBelowID   string", Number: 14},
							{Type: "add", Content: "\tLoadingAboveID    string", Number: 15},
							{Type: "add", Content: "\tLoadingBelowID    string", Number: 16},
							{Type: "add", Content: "\tClass             string", Number: 17},
							{Type: "add", Content: "\tAttributes        templ.Attributes", Number: 18},
							{Type: "add", Content: "}", Number: 19},
							{Type: "add", Content: "", Number: 20},
							{Type: "add", Content: "// Direction represents the scroll direction", Number: 21},
							{Type: "add", Content: "type Direction string", Number: 22},
							{Type: "add", Content: "", Number: 23},
							{Type: "add", Content: "const (", Number: 24},
							{Type: "add", Content: "\tDirectionAbove Direction = \"above\"", Number: 25},
							{Type: "add", Content: "\tDirectionBelow Direction = \"below\"", Number: 26},
							{Type: "add", Content: ")", Number: 27},
							{Type: "add", Content: "", Number: 28},
							{Type: "add", Content: "// HostArgs defines properties for the Host container", Number: 29},
							{Type: "add", Content: "type HostArgs struct {", Number: 30},
							{Type: "add", Content: "\tID         string", Number: 31},
							{Type: "add", Content: "\tClass      string", Number: 32},
							{Type: "add", Content: "\tAttributes templ.Attributes", Number: 33},
							{Type: "add", Content: "}", Number: 34},
							{Type: "add", Content: "", Number: 35},
							{Type: "add", Content: "// ItemsArgs defines properties for the Items container", Number: 36},
							{Type: "add", Content: "type ItemsArgs struct {", Number: 37},
							{Type: "add", Content: "\tID         string", Number: 38},
							{Type: "add", Content: "\tClass      string", Number: 39},
							{Type: "add", Content: "\tAttributes templ.Attributes", Number: 40},
							{Type: "add", Content: "}", Number: 41},
						},
					},
				},
			},
		}
		nextCursor = "1"
		hasMore = true
	case "1": // Page 2
		moreFiles = []infinitescrollpage.DiffFile{
			{
				Path:     "components/infinitescroll/variants.go",
				OldLines: 0,
				NewLines: 51,
				Hunks: []infinitescrollpage.DiffHunk{
					{
						Header: "@@ -0,0 +1,51 @@",
						Lines: []infinitescrollpage.DiffLine{
							{Type: "add", Content: "package infinitescroll", Number: 1},
							{Type: "add", Content: "", Number: 2},
							{Type: "add", Content: "import \"github.com/coreycole/datastarui/utils\"", Number: 3},
							{Type: "add", Content: "", Number: 4},
							{Type: "add", Content: "// InfiniteScrollVariants generates CSS classes for InfiniteScroll", Number: 5},
							{Type: "add", Content: "func InfiniteScrollVariants(args InfiniteScrollArgs) string {", Number: 6},
							{Type: "add", Content: "\tbaseClasses := \"relative\"", Number: 7},
							{Type: "add", Content: "\tif args.Class != \"\" {", Number: 8},
							{Type: "add", Content: "\t\treturn utils.TwMerge(baseClasses, args.Class)", Number: 9},
							{Type: "add", Content: "\t}", Number: 10},
							{Type: "add", Content: "\treturn baseClasses", Number: 11},
							{Type: "add", Content: "}", Number: 12},
							{Type: "add", Content: "", Number: 13},
							{Type: "add", Content: "// HostVariants generates CSS classes for Host container", Number: 14},
							{Type: "add", Content: "func HostVariants(args HostArgs) string {", Number: 15},
							{Type: "add", Content: "\tbaseClasses := \"overflow-auto\"", Number: 16},
							{Type: "add", Content: "\tif args.Class != \"\" {", Number: 17},
							{Type: "add", Content: "\t\treturn utils.TwMerge(baseClasses, args.Class)", Number: 18},
							{Type: "add", Content: "\t}", Number: 19},
							{Type: "add", Content: "\treturn baseClasses", Number: 20},
							{Type: "add", Content: "}", Number: 21},
							{Type: "add", Content: "", Number: 22},
							{Type: "add", Content: "// ItemsVariants generates CSS classes for Items container", Number: 23},
							{Type: "add", Content: "func ItemsVariants(args ItemsArgs) string {", Number: 24},
							{Type: "add", Content: "\tbaseClasses := \"\"", Number: 25},
							{Type: "add", Content: "\tif args.Class != \"\" {", Number: 26},
							{Type: "add", Content: "\t\treturn utils.TwMerge(baseClasses, args.Class)", Number: 27},
							{Type: "add", Content: "\t}", Number: 28},
							{Type: "add", Content: "\treturn baseClasses", Number: 29},
							{Type: "add", Content: "}", Number: 30},
						},
					},
				},
			},
			{
				Path:     "pages/components/infinitescrollpage/infinitescroll_page.templ",
				OldLines: 0,
				NewLines: 65,
				Hunks: []infinitescrollpage.DiffHunk{
					{
						Header: "@@ -0,0 +1,65 @@",
						Lines: []infinitescrollpage.DiffLine{
							{Type: "add", Content: "package infinitescrollpage", Number: 1},
							{Type: "add", Content: "", Number: 2},
							{Type: "add", Content: "import (", Number: 3},
							{Type: "add", Content: "\t\"fmt\"", Number: 4},
							{Type: "add", Content: "\t\"github.com/coreycole/datastarui/components/card\"", Number: 5},
							{Type: "add", Content: "\t\"github.com/coreycole/datastarui/components/infinitescroll\"", Number: 6},
							{Type: "add", Content: "\tl \"github.com/coreycole/datastarui/layouts\"", Number: 7},
							{Type: "add", Content: ")", Number: 8},
							{Type: "add", Content: "", Number: 9},
							{Type: "add", Content: "templ InfiniteScrollPage(rootArgs l.RootArgs) {", Number: 10},
							{Type: "add", Content: "\t{{", Number: 11},
							{Type: "add", Content: "\t\tpatchMoreExpr := \"@get('/api/infinitescroll/more')\"", Number: 12},
							{Type: "add", Content: "\t}}", Number: 13},
							{Type: "add", Content: "\t@l.Root(rootArgs) {", Number: 14},
							{Type: "add", Content: "\t\t<div class=\"space-y-8\">", Number: 15},
							{Type: "add", Content: "\t\t\t@l.ComponentPageBreadcrumbs(\"Infinite Scroll\")", Number: 16},
							{Type: "add", Content: "\t\t\t<div class=\"space-y-2\">", Number: 17},
							{Type: "add", Content: "\t\t\t\t<h1 class=\"text-3xl font-bold tracking-tight\">Infinite Scroll</h1>", Number: 18},
							{Type: "add", Content: "\t\t\t\t<p class=\"text-lg text-muted-foreground\">", Number: 19},
							{Type: "add", Content: "\t\t\t\t\tLoad content progressively as the user scrolls", Number: 20},
							{Type: "add", Content: "\t\t\t\t</p>", Number: 21},
							{Type: "add", Content: "\t\t\t</div>", Number: 22},
							{Type: "add", Content: "\t\t\t<section class=\"space-y-6\">", Number: 23},
							{Type: "add", Content: "\t\t\t\t<h2 class=\"text-2xl font-semibold tracking-tight\">Demo</h2>", Number: 24},
							{Type: "add", Content: "\t\t\t\t@card.Card(card.CardArgs{}) {", Number: 25},
							{Type: "add", Content: "\t\t\t\t\t@infinitescroll.InfiniteScroll(infinitescroll.InfiniteScrollArgs{", Number: 26},
							{Type: "add", Content: "\t\t\t\t\t\tID:             \"diff_viewer\",", Number: 27},
							{Type: "add", Content: "\t\t\t\t\t\tPatchBelowExpr: patchMoreExpr,", Number: 28},
							{Type: "add", Content: "\t\t\t\t\t}) {", Number: 29},
							{Type: "add", Content: "\t\t\t\t\t\t<!-- Initial content -->", Number: 30},
							{Type: "add", Content: "\t\t\t\t\t}", Number: 31},
							{Type: "add", Content: "\t\t\t\t}", Number: 32},
							{Type: "add", Content: "\t\t\t</section>", Number: 33},
							{Type: "add", Content: "\t\t</div>", Number: 34},
							{Type: "add", Content: "\t}", Number: 35},
							{Type: "add", Content: "}", Number: 36},
						},
					},
				},
			},
		}
		nextCursor = "2"
		hasMore = true
	default: // Page 3+ - exhausted
		hasMore = false
	}

	// 4. Build HTML for new file cards
	var itemsHTML strings.Builder
	startIndex := 1 // First page already shows file 0
	if cursor == "1" {
		startIndex = 2 // Page 2 starts after page 1
	} else if cursor != "0" {
		startIndex = 3 // Page 3+
	}

	for i, file := range moreFiles {
		fileCard := infinitescrollpage.DiffFileCard(file, startIndex+i)
		if err := fileCard.Render(context.Background(), &itemsHTML); err != nil {
			return err
		}
	}

	// 5. Append new items to Items container (View Transitions OFF)
	sse.PatchElements(itemsHTML.String(),
		datastar.WithSelectorID("diff_viewer-items"),
		datastar.WithModeAppend(),
		datastar.WithoutViewTransitions(),
	)

	// 6. Either remint sentinel with next cursor or exhaust
	if hasMore {
		// Remint sentinel with next cursor
		newSentinel := infinitescroll.Sentinel(infinitescroll.SentinelArgs{
			ID:        "diff_viewer-sentinel-below",
			Direction: infinitescroll.DirectionBelow,
			PatchExpr: fmt.Sprintf("@get('/api/infinitescroll/more?cursor=%s')", nextCursor),
		})
		return sse.PatchElementTempl(newSentinel)
	} else {
		// Exhausted - remove sentinel
		return sse.RemoveElementByID("diff_viewer-sentinel-below")
	}
}

func main() {
	// Load configuration from environment
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	// Check if datastar pro file exists. The layout still falls back to the
	// public CDN if the licensed local asset is absent.
	if _, err := os.Stat("static/js/datastar-pro-v1.js"); err == nil {
		cfg.DatastarProAvailable = true
	}
	if _, err := os.Stat("static/js/datastar-inspector.js"); err == nil {
		cfg.DatastarInspectorEnabled = true
	}

	// Create a new Echo instance
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Mobile detection middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isMobile := utils.IsMobile(c)
			// Add mobile detection to context
			ctx := utils.WithMobile(c.Request().Context(), isMobile)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	})

	// Setup Connect RPC services
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

	// Mount Connect RPC at /connect/* with unwrap middleware
	unwrappedHandler := apimw.UnwrapFormData(connectMux)
	e.Any("/connect/*", echo.WrapHandler(http.StripPrefix("/connect", unwrappedHandler)))

	// Create Connect RPC client for internal use by HTTP form handlers
	authClient := authconnect.NewAuthServiceClient(
		http.DefaultClient,
		"http://localhost:"+cfg.Port+"/connect",
	)

	// Create HTTP form handlers
	loginHandler := loginform.NewHandler(authClient)

	// Register HTTP form routes (return HTML for browser interactions)
	e.POST("/forms/login", loginHandler.HandleLoginForm)

	// Serve the home page at the root route
	e.GET("/", func(c echo.Context) error {
		rootArgs := l.RootArgs{
			CurrentPage:          "home",
			CurrentPath:          c.Request().URL.Path,
			InspectorEnabled:     cfg.DatastarInspectorEnabled,
			DatastarProAvailable: cfg.DatastarProAvailable,
		}
		component := p.HomePage(rootArgs)
		return component.Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve the components page
	e.GET("/components", func(c echo.Context) error {
		return componentspage.ComponentsPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/button", func(c echo.Context) error {
		return buttonpage.ButtonPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/breadcrumb", func(c echo.Context) error {
		return breadcrumbpage.BreadcrumbPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/dropdown", func(c echo.Context) error {
		return dropdownpage.DropdownPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/calendar", func(c echo.Context) error {
		return calendarpage.Page(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/form", func(c echo.Context) error {
		return formpage.FormPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/infinitescroll", func(c echo.Context) error {
		return infinitescrollpage.InfiniteScrollPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/popover", func(c echo.Context) error {
		return popoverpage.PopoverPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/tabs", func(c echo.Context) error {
		return tabspage.TabsPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/card", func(c echo.Context) error {
		return cardpage.CardPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/checkbox", func(c echo.Context) error {
		return checkboxpage.CheckboxPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/dialog", func(c echo.Context) error {
		return dialogpage.DialogPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/select", func(c echo.Context) error {
		return selectpage.SelectPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/datepicker", func(c echo.Context) error {
		return datepickerpage.Page(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/sheet", func(c echo.Context) error {
		return sheetpage.SheetPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/sidebar", func(c echo.Context) error {
		return sidebarpage.SidebarPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/tooltip", func(c echo.Context) error {
		return tooltippage.TooltipPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/toast", func(c echo.Context) error {
		return toastpage.ToastPage(componentRootArgs(c.Request().URL.Path, cfg)).Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve the docs page
	e.GET("/docs", func(c echo.Context) error {
		rootArgs := l.RootArgs{
			CurrentPage:          "docs",
			CurrentPath:          c.Request().URL.Path,
			InspectorEnabled:     cfg.DatastarInspectorEnabled,
			DatastarProAvailable: cfg.DatastarProAvailable,
		}
		return p.DocsPage(rootArgs).Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve the examples page
	e.GET("/examples", func(c echo.Context) error {
		rootArgs := l.RootArgs{
			CurrentPage:          "examples",
			CurrentPath:          c.Request().URL.Path,
			InspectorEnabled:     cfg.DatastarInspectorEnabled,
			DatastarProAvailable: cfg.DatastarProAvailable,
		}
		return p.ExamplesPage(rootArgs).Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve the login page
	e.GET("/login", func(c echo.Context) error {
		rootArgs := l.RootArgs{
			CurrentPage:          "login",
			CurrentPath:          c.Request().URL.Path,
			InspectorEnabled:     cfg.DatastarInspectorEnabled,
			DatastarProAvailable: cfg.DatastarProAvailable,
		}
		return login.LoginPage(rootArgs).Render(c.Request().Context(), c.Response().Writer)
	})

	// API handlers for component demos
	e.GET("/api/infinitescroll/more", func(c echo.Context) error {
		return handleInfiniteScrollMore(c)
	})

	// Serve static files
	e.Static("/", "static/")

	// Start the server
	if err := e.Start(":" + cfg.Port); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
