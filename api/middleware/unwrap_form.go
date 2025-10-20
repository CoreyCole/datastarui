package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// UnwrapFormData is HTTP middleware that unwraps nested form data from Datastar signals.
// When a form sends {formId: "login", login: {field1: "value1", field2: "value2"}},
// this middleware extracts the nested object to {field1: "value1", field2: "value2"}.
func UnwrapFormData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only process POST requests with JSON content
		if r.Method != "POST" {
			next.ServeHTTP(w, r)
			return
	}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" && contentType != "application/connect+json" {
			next.ServeHTTP(w, r)
			return
		}

		// Read the request body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		r.Body.Close()

		// Try to parse as JSON
		var dataMap map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &dataMap); err != nil {
			// Not valid JSON or can't parse, pass through unchanged
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			next.ServeHTTP(w, r)
			return
		}

		// Check if request contains formId metadata
		if formIdValue, hasFormId := dataMap["formId"]; hasFormId {
			if formId, ok := formIdValue.(string); ok && formId != "" {
				// Look for the nested object with the formId key
				if nestedMap, hasNested := dataMap[formId]; hasNested {
					if formData, ok := nestedMap.(map[string]interface{}); ok {
						// Extract the form data
						unwrappedBytes, err := json.Marshal(formData)
						if err != nil {
							// Failed to marshal, use original
							r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
							next.ServeHTTP(w, r)
							return
						}

						// Replace request body with unwrapped data
						r.Body = io.NopCloser(bytes.NewBuffer(unwrappedBytes))
						r.ContentLength = int64(len(unwrappedBytes))
						next.ServeHTTP(w, r)
						return
					}
				}
			}
		}

		// No unwrapping needed, pass through original
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		next.ServeHTTP(w, r)
	})
}
