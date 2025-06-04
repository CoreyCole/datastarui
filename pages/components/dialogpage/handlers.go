package dialogpage

import (
	"encoding/json"
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

// RegisterDialogPageHandlers registers all dialog demo route handlers
func RegisterDialogPageHandlers(e *echo.Echo) {
	// Form Dialog Handler
	e.POST("/dialog/dialog-page/form-submit", func(c echo.Context) error {
		// Parse multipart form data
		err := c.Request().ParseMultipartForm(32 << 20) // 32 MB max memory
		if err != nil {
			log.Printf("Error parsing form: %v", err)
			sse := datastar.NewSSE(c.Response().Writer, c.Request())
			sse.MergeFragments(`<div id="form_dialog_errors" class="p-4 bg-red-50 border border-red-200 rounded-md mb-4">
				<p class="text-sm text-red-800">Error processing form</p>
			</div>`)
			return nil
		}

		// Get form data
		name := c.FormValue("name")
		email := c.FormValue("email")

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

		sse := datastar.NewSSE(c.Response().Writer, c.Request())

		// If there are validation errors, add the error div
		if len(errors) > 0 {
			errorHTML := generateErrorHTML("form_dialog", errors)
			sse.MergeFragments(errorHTML)
			return nil
		}

		// Show success message in the main page (not in the dialog)
		successHTML := `<div id="form_success_message" class="p-4 bg-green-50 border border-green-200 rounded-md mb-4">
			<p class="text-sm text-green-800">✓ Form submitted successfully! Name: ` + name + `, Email: ` + email + `</p>
		</div>`

		// Clear any existing errors in the dialog
		clearErrorsHTML := `<div id="form_dialog_errors"></div>`

		// Merge both success message and clear errors
		sse.MergeFragments(successHTML + clearErrorsHTML)

		// Update signals in a single merged call to ensure reliable state updates
		allSignalsJSON, _ := json.Marshal(map[string]interface{}{
			"form_dialog": map[string]interface{}{
				"name":      name,
				"email":     email,
				"submitted": true,
			},
			"form_demo": map[string]interface{}{
				"open": false,
			},
		})
		sse.MergeSignals(allSignalsJSON)

		// Log the submission (in a real app, you'd save to database)
		log.Printf("Dialog form submitted - Name: %s, Email: %s", name, email)

		return nil
	})
}
