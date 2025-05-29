package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	p "github.com/coreycole/datastarui/pages"
	"github.com/coreycole/datastarui/pages/components/breadcrumbpage"
	"github.com/coreycole/datastarui/pages/components/buttonpage"
	"github.com/coreycole/datastarui/pages/components/cardpage"
	"github.com/coreycole/datastarui/pages/components/checkboxpage"
	"github.com/coreycole/datastarui/pages/components/dropdownpage"
	"github.com/coreycole/datastarui/pages/components/formpage"
	"github.com/coreycole/datastarui/pages/components/popoverpage"
	"github.com/coreycole/datastarui/pages/components/tabspage"
)

const port = "4242"

func main() {
	// Create a new Echo instance
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Serve the home page at the root route
	e.GET("/", func(c echo.Context) error {
		component := p.HomePage()
		return component.Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve the components page
	e.GET("/components", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/components/button")
	})
	e.GET("/components/button", func(c echo.Context) error {
		return buttonpage.ButtonPage().Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/breadcrumb", func(c echo.Context) error {
		return breadcrumbpage.BreadcrumbPage().Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/dropdown", func(c echo.Context) error {
		return dropdownpage.DropdownPage().Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/form", func(c echo.Context) error {
		return formpage.FormPage().Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/popover", func(c echo.Context) error {
		return popoverpage.PopoverPage().Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/tabs", func(c echo.Context) error {
		return tabspage.TabsPage().Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/card", func(c echo.Context) error {
		return cardpage.CardPage().Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/checkbox", func(c echo.Context) error {
		return checkboxpage.CheckboxPage().Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve the docs page
	e.GET("/docs", func(c echo.Context) error {
		return p.DocsPage().Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve the examples page
	e.GET("/examples", func(c echo.Context) error {
		return p.ExamplesPage().Render(c.Request().Context(), c.Response().Writer)
	})

	// Register form demo handlers
	formpage.RegisterFormPageHandlers(e)
	checkboxpage.RegisterCheckboxHandlers(e)

	// Serve static files
	e.Static("/", "static/")

	// Start the server
	if err := e.Start(":" + port); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
