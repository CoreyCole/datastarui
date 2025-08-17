package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	p "github.com/coreycole/datastarui/pages"
	l "github.com/coreycole/datastarui/layouts"
	"github.com/coreycole/datastarui/pages/components/breadcrumbpage"
	"github.com/coreycole/datastarui/pages/components/buttonpage"
	"github.com/coreycole/datastarui/pages/components/calendarpage"
	"github.com/coreycole/datastarui/pages/components/cardpage"
	"github.com/coreycole/datastarui/pages/components/checkboxpage"
	"github.com/coreycole/datastarui/pages/components/datepickerpage"
	"github.com/coreycole/datastarui/pages/components/dialogpage"
	"github.com/coreycole/datastarui/pages/components/dropdownpage"
	"github.com/coreycole/datastarui/pages/components/sheetpage"
	"github.com/coreycole/datastarui/pages/components/sidebarpage"
	"github.com/coreycole/datastarui/pages/components/formpage"
	"github.com/coreycole/datastarui/pages/components/popoverpage"
	"github.com/coreycole/datastarui/pages/components/selectpage"
	"github.com/coreycole/datastarui/pages/components/tabspage"
	"github.com/coreycole/datastarui/pages/components/componentspage"
	"github.com/coreycole/datastarui/utils"
)

const port = "4242"

// Config holds environment configuration
type Config struct {
	DatastarInspectorEnabled bool `envconfig:"DATASTAR_INSPECTOR_ENABLED" default:"false"`
}

// Helper function to create RootArgs for component pages
func componentRootArgs(path string, cfg Config, datastarProAvailable bool) l.RootArgs {
	return l.RootArgs{
		CurrentPage:          "components",
		CurrentPath:          path,
		InspectorEnabled:     cfg.DatastarInspectorEnabled,
		DatastarProAvailable: datastarProAvailable,
	}
}

func main() {
	// Load configuration from environment
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	// Check if datastar pro file exists
	datastarProAvailable := false
	if _, err := os.Stat("static/js/datastar.js"); err == nil {
		datastarProAvailable = true
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

	// Serve the home page at the root route
	e.GET("/", func(c echo.Context) error {
		rootArgs := l.RootArgs{
			CurrentPage:          "home",
			CurrentPath:          c.Request().URL.Path,
			InspectorEnabled:     cfg.DatastarInspectorEnabled,
			DatastarProAvailable: datastarProAvailable,
		}
		component := p.HomePage(rootArgs)
		return component.Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve the components page
	e.GET("/components", func(c echo.Context) error {
		return componentspage.ComponentsPage(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/button", func(c echo.Context) error {
		return buttonpage.ButtonPage(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/breadcrumb", func(c echo.Context) error {
		return breadcrumbpage.BreadcrumbPage(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/dropdown", func(c echo.Context) error {
		return dropdownpage.DropdownPage(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/calendar", func(c echo.Context) error {
		return calendarpage.Page(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/form", func(c echo.Context) error {
		return formpage.FormPage(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/popover", func(c echo.Context) error {
		return popoverpage.PopoverPage(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/tabs", func(c echo.Context) error {
		return tabspage.TabsPage(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/card", func(c echo.Context) error {
		return cardpage.CardPage(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/checkbox", func(c echo.Context) error {
		return checkboxpage.CheckboxPage(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/dialog", func(c echo.Context) error {
		return dialogpage.DialogPage(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/select", func(c echo.Context) error {
		return selectpage.SelectPage(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/datepicker", func(c echo.Context) error {
		return datepickerpage.Page(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/sheet", func(c echo.Context) error {
		return sheetpage.SheetPage(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/sidebar", func(c echo.Context) error {
		return sidebarpage.SidebarPage(componentRootArgs(c.Request().URL.Path, cfg, datastarProAvailable)).Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve the docs page
	e.GET("/docs", func(c echo.Context) error {
		rootArgs := l.RootArgs{
			CurrentPage:          "docs",
			CurrentPath:          c.Request().URL.Path,
			InspectorEnabled:     cfg.DatastarInspectorEnabled,
			DatastarProAvailable: datastarProAvailable,
		}
		return p.DocsPage(rootArgs).Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve the examples page
	e.GET("/examples", func(c echo.Context) error {
		rootArgs := l.RootArgs{
			CurrentPage:          "examples",
			CurrentPath:          c.Request().URL.Path,
			InspectorEnabled:     cfg.DatastarInspectorEnabled,
			DatastarProAvailable: datastarProAvailable,
		}
		return p.ExamplesPage(rootArgs).Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve static files
	e.Static("/", "static/")

	// Start the server
	if err := e.Start(":" + port); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
