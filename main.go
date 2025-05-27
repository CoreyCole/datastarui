package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	datastar "github.com/starfederation/datastar/sdk/go"

	p "github.com/coreycole/datastarui/pages"
	"github.com/coreycole/datastarui/pages/components/breadcrumb"
	"github.com/coreycole/datastarui/pages/components/button"
	"github.com/coreycole/datastarui/pages/components/card"
	"github.com/coreycole/datastarui/pages/components/dropdown"
	"github.com/coreycole/datastarui/pages/components/form"
	"github.com/coreycole/datastarui/pages/components/tabs"
)

const port = "4242"

// generateErrorHTML creates a single error div with field-specific errors
func generateErrorHTML(formID string, errors map[string][]string) string {
	if len(errors) == 0 {
		return ""
	}

	errorHTML := `<div id="` + formID + `-errors" class="p-4 bg-red-50 border border-red-200 rounded-md mb-4"><div class="text-sm text-red-800">`
	for field, fieldErrors := range errors {
		errorHTML += `<div class="mb-2"><strong>` + field + `:</strong><ul class="ml-4">`
		for _, err := range fieldErrors {
			errorHTML += `<li>• ` + err + `</li>`
		}
		errorHTML += `</ul></div>`
	}
	errorHTML += `</div></div>`
	return errorHTML
}

// generateSuccessHTML creates a success message div
func generateSuccessHTML(formID string, message string) string {
	return `<div id="` + formID + `-errors" class="p-4 bg-green-50 border border-green-200 rounded-md mb-4">
		<p class="text-sm text-green-800">` + message + `</p>
	</div>`
}

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
		return button.ButtonPage().Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/breadcrumb", func(c echo.Context) error {
		return breadcrumb.BreadcrumbPage().Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/dropdown", func(c echo.Context) error {
		return dropdown.DropdownPage().Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/form", func(c echo.Context) error {
		return form.FormPage().Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/tabs", func(c echo.Context) error {
		return tabs.TabsPage().Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/components/card", func(c echo.Context) error {
		return card.CardPage().Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve the docs page
	e.GET("/docs", func(c echo.Context) error {
		return p.DocsPage().Render(c.Request().Context(), c.Response().Writer)
	})

	// Serve the examples page
	e.GET("/examples", func(c echo.Context) error {
		return p.ExamplesPage().Render(c.Request().Context(), c.Response().Writer)
	})

	// Form submission handlers
	e.POST("/form/form-page/basic-form", func(c echo.Context) error {
		// Parse multipart form data
		err := c.Request().ParseMultipartForm(32 << 20) // 32 MB max memory
		if err != nil {
			log.Printf("Error parsing form: %v", err)
			sse := datastar.NewSSE(c.Response().Writer, c.Request())
			sse.MergeFragments(`<div id="basic-form-errors" class="p-4 bg-red-50 border border-red-200 rounded-md mb-4">
				<p class="text-sm text-red-800">Error processing form</p>
			</div>`)
			return nil
		}

		// Get form data
		username := c.FormValue("username")

		// Validate form data using field-specific error map
		errors := make(map[string][]string)

		if username == "" {
			errors["username"] = append(errors["username"], "Username is required")
		}
		if len(username) > 0 && len(username) < 3 {
			errors["username"] = append(errors["username"], "Username must be at least 3 characters")
		}

		sse := datastar.NewSSE(c.Response().Writer, c.Request())

		// If there are validation errors, add the error div
		if len(errors) > 0 {
			errorHTML := generateErrorHTML("basic-form", errors)
			sse.MergeFragments(errorHTML)
			return nil
		}

		// Show success message
		successHTML := generateSuccessHTML("basic-form", "Profile updated successfully!")
		sse.MergeFragments(successHTML)

		// Log the submission (in a real app, you'd save to database)
		log.Printf("Basic form submitted - Username: %s", username)

		return nil
	})

	e.POST("/form/form-page/validation-form", func(c echo.Context) error {
		// Parse multipart form data
		err := c.Request().ParseMultipartForm(32 << 20)
		if err != nil {
			log.Printf("Error parsing form: %v", err)
			sse := datastar.NewSSE(c.Response().Writer, c.Request())
			sse.MergeFragments(`<div id="validation-form-errors" class="p-4 bg-red-50 border border-red-200 rounded-md mb-4">
				<p class="text-sm text-red-800">Error processing form</p>
			</div>`)
			return nil
		}

		// Get form data
		email := c.FormValue("email")
		password := c.FormValue("password")

		// Validate form data using field-specific error map
		errors := make(map[string][]string)

		if email == "" {
			errors["email"] = append(errors["email"], "Email is required")
		} else {
			if !strings.Contains(email, "@") {
				errors["email"] = append(errors["email"], "Please enter a valid email")
			}
			if !strings.Contains(email, ".") {
				errors["email"] = append(errors["email"], "Email must contain a domain")
			}
		}

		if password == "" {
			errors["password"] = append(errors["password"], "Password is required")
		} else {
			if len(password) < 8 {
				errors["password"] = append(errors["password"], "Password must be at least 8 characters")
			}
			if !strings.ContainsAny(password, "0123456789") {
				errors["password"] = append(errors["password"], "Password must contain at least one number")
			}
		}

		sse := datastar.NewSSE(c.Response().Writer, c.Request())

		// If there are validation errors, merge them
		if len(errors) > 0 {
			errorHTML := generateErrorHTML("validation-form", errors)
			sse.MergeFragments(errorHTML)
			return nil
		}

		// Log the submission
		log.Printf("Validation form submitted - Email: %s, Password length: %d", email, len(password))

		// Show success message
		successHTML := generateSuccessHTML("validation-form", "Account created successfully!")
		sse.MergeFragments(successHTML)
		return nil
	})

	e.POST("/form/form-page/contact-form", func(c echo.Context) error {
		// Parse multipart form data
		err := c.Request().ParseMultipartForm(32 << 20)
		if err != nil {
			log.Printf("Error parsing form: %v", err)
			sse := datastar.NewSSE(c.Response().Writer, c.Request())
			sse.MergeFragments(`<div id="contact-form-errors" class="p-4 bg-red-50 border border-red-200 rounded-md mb-4">
				<p class="text-sm text-red-800">Error processing form</p>
			</div>`)
			return nil
		}

		// Get form data
		name := c.FormValue("name")
		email := c.FormValue("email")
		subject := c.FormValue("subject")
		message := c.FormValue("message")

		// Validate form data using field-specific error map
		errors := make(map[string][]string)

		if name == "" {
			errors["name"] = append(errors["name"], "Name is required")
		} else if len(name) < 2 {
			errors["name"] = append(errors["name"], "Name must be at least 2 characters")
		}

		if email == "" {
			errors["email"] = append(errors["email"], "Email is required")
		} else {
			if !strings.Contains(email, "@") {
				errors["email"] = append(errors["email"], "Please enter a valid email")
			}
			if !strings.Contains(email, ".") {
				errors["email"] = append(errors["email"], "Email must contain a domain")
			}
		}

		if subject == "" {
			errors["subject"] = append(errors["subject"], "Subject is required")
		} else if len(subject) < 5 {
			errors["subject"] = append(errors["subject"], "Subject must be at least 5 characters")
		}

		if message == "" {
			errors["message"] = append(errors["message"], "Message is required")
		} else if len(message) < 10 {
			errors["message"] = append(errors["message"], "Message must be at least 10 characters")
		}

		sse := datastar.NewSSE(c.Response().Writer, c.Request())

		// If there are validation errors, merge them
		if len(errors) > 0 {
			errorHTML := generateErrorHTML("contact-form", errors)
			sse.MergeFragments(errorHTML)
			return nil
		}

		// Log the submission
		log.Printf("Contact form submitted - Name: %s, Email: %s, Subject: %s, Message length: %d", name, email, subject, len(message))

		// Show success message
		successHTML := generateSuccessHTML("contact-form", "Thank you! Your message has been sent successfully.")
		sse.MergeFragments(successHTML)
		return nil
	})

	// Serve static files
	e.Static("/", "static/")

	// Start the server
	if err := e.Start(":" + port); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
