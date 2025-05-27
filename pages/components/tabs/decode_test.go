package tabs

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestDecodeAccountForm(t *testing.T) {
	// Create a test HTTP request with form data
	formData := url.Values{}
	formData.Set("name", "John Doe")
	formData.Set("username", "johndoe")

	req, err := http.NewRequest("POST", "/test", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Test the decoder
	form, response, err := DecodeAccountForm(req)
	if err != nil {
		t.Fatalf("DecodeAccountForm failed: %v", err)
	}

	// Check if decoding was successful
	if !response.Success {
		t.Fatalf("Form validation failed: %s, errors: %v", response.Message, response.Errors)
	}

	// Verify the parsed data
	if form.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", form.Name)
	}
	if form.Username != "johndoe" {
		t.Errorf("Expected username 'johndoe', got '%s'", form.Username)
	}
}

func TestDecodePasswordForm(t *testing.T) {
	// Create a test HTTP request with form data
	formData := url.Values{}
	formData.Set("current_password", "oldpassword123")
	formData.Set("new_password", "newpassword456")

	req, err := http.NewRequest("POST", "/test", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Test the decoder
	form, response, err := DecodePasswordForm(req)
	if err != nil {
		t.Fatalf("DecodePasswordForm failed: %v", err)
	}

	// Check if decoding was successful
	if !response.Success {
		t.Fatalf("Form validation failed: %s, errors: %v", response.Message, response.Errors)
	}

	// Verify the parsed data
	if form.CurrentPassword != "oldpassword123" {
		t.Errorf("Expected current_password 'oldpassword123', got '%s'", form.CurrentPassword)
	}
	if form.NewPassword != "newpassword456" {
		t.Errorf("Expected new_password 'newpassword456', got '%s'", form.NewPassword)
	}
}

func TestFormValidation(t *testing.T) {
	// Test validation failure with invalid username
	formData := url.Values{}
	formData.Set("name", "A")      // Too short
	formData.Set("username", "ab") // Too short

	req, err := http.NewRequest("POST", "/test", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Test the decoder
	_, response, err := DecodeAccountForm(req)
	if err != nil {
		t.Fatalf("DecodeAccountForm failed: %v", err)
	}

	// Check if validation properly failed
	if response.Success {
		t.Error("Expected validation to fail for invalid input")
	}

	// Should have validation errors
	if len(response.Errors) == 0 {
		t.Error("Expected validation errors but got none")
	}
}
