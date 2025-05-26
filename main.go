package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	p "github.com/coreycole/datastarui/pages"
	pc "github.com/coreycole/datastarui/pages/components"
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
		return c.Redirect(http.StatusFound, "/components/breadcrumb")
	})
	e.GET("/components/button", func(c echo.Context) error {
		component := pc.ButtonPage()
		return component.Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/breadcrumb", func(c echo.Context) error {
		component := pc.BreadcrumbPage()
		return component.Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/dropdown", func(c echo.Context) error {
		component := pc.DropdownPage()
		return component.Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve the docs page
	e.GET("/docs", func(c echo.Context) error {
		component := p.DocsPage()
		return component.Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve the examples page
	e.GET("/examples", func(c echo.Context) error {
		component := p.ExamplesPage()
		return component.Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve static files
	e.Static("/", "static/")

	// Start the server
	if err := e.Start(":" + port); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
