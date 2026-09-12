package main

import (
	"log"
	"net/http"
	"os"

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

	"context"
	"strings"
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

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	sse := datastar.NewSSE(w, r)

	// 1. Show loading indicator (same-id DOM replace on sentinel)
	loadingComponent := infinitescroll.Loading(infinitescroll.LoadingArgs{
		ID:        "diff_viewer-sentinel-below",
		Direction: infinitescroll.DirectionBelow,
	})
	sse.PatchElementTempl(loadingComponent)

	// 2. Simulate fetching more files (in real app, query database with cursor)
	moreFiles := []infinitescrollpage.DiffFile{
		{
			Path:     "pages/components/infinitescrollpage/infinitescroll_page.templ",
			OldLines: 0,
			NewLines: 180,
			Hunks: []infinitescrollpage.DiffHunk{
				{
					Header: "@@ -0,0 +1,180 @@",
					Lines: []infinitescrollpage.DiffLine{
						{Type: "add", Content: "package infinitescrollpage", Number: 1},
						{Type: "add", Content: "", Number: 2},
						{Type: "add", Content: "import (", Number: 3},
						{Type: "add", Content: "\t\"github.com/coreycole/datastarui/components/infinitescroll\"", Number: 4},
						{Type: "context", Content: "\t...", Number: 5},
						{Type: "add", Content: ")", Number: 6},
					},
				},
			},
		},
		{
			Path:     "main.go",
			OldLines: 222,
			NewLines: 230,
			Hunks: []infinitescrollpage.DiffHunk{
				{
					Header: "@@ -222,6 +222,14 @@",
					Lines: []infinitescrollpage.DiffLine{
						{Type: "context", Content: "\t})", Number: 222},
						{Type: "context", Content: "", Number: 223},
						{Type: "add", Content: "\t// API handlers for component demos", Number: 224},
						{Type: "add", Content: "\te.GET(\"/api/infinitescroll/more\", func(c echo.Context) error {", Number: 225},
						{Type: "add", Content: "\t\treturn handleInfiniteScrollMore(c)", Number: 226},
						{Type: "add", Content: "\t})", Number: 227},
						{Type: "add", Content: "", Number: 228},
						{Type: "context", Content: "\t// Serve static files", Number: 229},
						{Type: "context", Content: "\te.Static(\"/\", \"static/\")", Number: 230},
					},
				},
			},
		},
	}

	// 3. Build HTML for new file cards
	var itemsHTML strings.Builder
	for i, file := range moreFiles {
		fileCard := infinitescrollpage.DiffFileCard(file, 3+i)
		if err := fileCard.Render(context.Background(), &itemsHTML); err != nil {
			return err
		}
	}

	// 4. Append new items to Items container (View Transitions OFF - default for PatchElements)
	sse.PatchElements(itemsHTML.String(),
		datastar.WithSelectorID("diff_viewer-items"),
		datastar.WithModeAppend(),
		datastar.WithoutViewTransitions(),
	)

	// 5. Exhausted - remove sentinel (in real app, check if hasMore from database)
	// For demo, we'll remove it after this batch
	return sse.RemoveElementByID("diff_viewer-sentinel-below")

	// 5. Alternative: If more content available, remint sentinel
	// newSentinel := infinitescroll.Sentinel(infinitescroll.SentinelArgs{
	// 	ID:        "diff_viewer-sentinel-below",
	// 	Direction: infinitescroll.DirectionBelow,
	// 	PatchExpr: "@get('/api/infinitescroll/more?cursor=xyz')",
	// })
	// return sse.PatchElementTempl(newSentinel)
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
