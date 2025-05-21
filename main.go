package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	p "github.com/coreycole/datastarui/pages"
)

func main() {
	// Create a new Echo instance
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Serve the templ component at the root route
	e.GET("/", func(c echo.Context) error {
		component := p.DropdownExample()
		return component.Render(c.Request().Context(), c.Response().Writer)
	})

	e.Static("/", "static/")

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
