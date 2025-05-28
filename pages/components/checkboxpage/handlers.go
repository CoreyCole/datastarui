package checkboxpage

import (
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	datastar "github.com/starfederation/datastar/sdk/go"
)

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

// RegisterCheckboxHandlers registers the checkbox demo form handlers
func RegisterCheckboxHandlers(e *echo.Echo) {
	e.POST("/forms/checkbox-demo", handleCheckboxDemoForm)
}

func handleCheckboxDemoForm(c echo.Context) error {
	// Parse multipart form data
	err := c.Request().ParseMultipartForm(32 << 20) // 32 MB max memory
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		sse.MergeFragments(`<div id="checkbox-form-errors" class="p-4 bg-red-50 border border-red-200 rounded-md mb-4">
			<p class="text-sm text-red-800">Error processing form</p>
		</div>`)
		return nil
	}

	// Get form data
	name := c.FormValue("name")
	email := c.FormValue("email")
	terms := c.FormValue("terms-form")

	// Debug logging to see what we're receiving
	log.Printf("DEBUG - Received form data: name='%s', email='%s', terms='%s'", name, email, terms)

	// Validate form data using field-specific error map
	errors := make(map[string][]string)

	if strings.TrimSpace(name) == "" {
		errors["name"] = append(errors["name"], "Name is required")
	} else if len(strings.TrimSpace(name)) < 2 {
		errors["name"] = append(errors["name"], "Name must be at least 2 characters")
	}

	if strings.TrimSpace(email) == "" {
		errors["email"] = append(errors["email"], "Email is required")
	} else {
		if !strings.Contains(email, "@") {
			errors["email"] = append(errors["email"], "Please enter a valid email")
		}
		if !strings.Contains(email, ".") {
			errors["email"] = append(errors["email"], "Email must contain a domain")
		}
	}

	// Check for checkbox - handle boolean string values
	if terms != "true" && terms != "on" && terms != "" {
		log.Printf("DEBUG - Unexpected terms value: '%s'", terms)
	}

	// Accept "true" (boolean true as string) or "on" (standard checkbox value)
	if terms != "true" && terms != "on" {
		errors["terms"] = append(errors["terms"], "You must accept the terms and conditions")
	}

	sse := datastar.NewSSE(c.Response().Writer, c.Request())

	// If there are validation errors, merge them
	if len(errors) > 0 {
		errorHTML := generateErrorHTML("checkbox-form", errors)
		sse.MergeFragments(errorHTML)
		return nil
	}

	// Log the submission (in a real app, you'd save to database)
	log.Printf("Checkbox demo form submitted - Name: %s, Email: %s, Terms: %s", name, email, terms)

	// Show success message
	successHTML := generateSuccessHTML("checkbox-form", "Account created successfully!")
	sse.MergeFragments(successHTML)

	return nil
}
