package tabs

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestFormDecoderBasic(t *testing.T) {
	// Create a simple form decoder
	decoder := NewFormDecoder()

	// Create test form data
	formData := url.Values{}
	formData.Set("name", "Test User")
	formData.Set("username", "testuser")

	// Create HTTP request
	req, err := http.NewRequest("POST", "/test", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create an AccountForm instance
	form := &AccountForm{}

	// Test the basic decode functionality
	err = decoder.DecodeForm(req, form)
	if err != nil {
		t.Fatalf("DecodeForm failed: %v", err)
	}

	// Verify the data was parsed correctly
	if form.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", form.Name)
	}
	if form.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", form.Username)
	}

	t.Logf("Successfully decoded form: Name=%s, Username=%s", form.Name, form.Username)
}

func TestFormDecoderValidation(t *testing.T) {
	// Create a form decoder
	decoder := NewFormDecoder()

	// Create test form data with invalid values
	formData := url.Values{}
	formData.Set("name", "A")      // Too short (should be at least 2 chars)
	formData.Set("username", "ab") // Too short (should be at least 3 chars)

	// Create HTTP request
	req, err := http.NewRequest("POST", "/test", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create an AccountForm instance
	form := &AccountForm{}

	// Test the validation functionality
	response, err := decoder.DecodeFormWithValidation(req, form)
	if err != nil {
		t.Fatalf("DecodeFormWithValidation failed: %v", err)
	}

	// Should have validation errors
	if response.Success {
		t.Error("Expected validation to fail for invalid input")
	}

	if len(response.Errors) == 0 {
		t.Error("Expected validation errors but got none")
	}

	t.Logf("Validation correctly failed with errors: %v", response.Errors)
}

func TestPasswordFormDecoding(t *testing.T) {
	// Create a form decoder
	decoder := NewFormDecoder()

	// Create test form data
	formData := url.Values{}
	formData.Set("current_password", "myoldpassword123")
	formData.Set("new_password", "mynewpassword456")

	// Create HTTP request
	req, err := http.NewRequest("POST", "/test", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a PasswordForm instance
	form := &PasswordForm{}

	// Test the decode functionality
	err = decoder.DecodeForm(req, form)
	if err != nil {
		t.Fatalf("DecodeForm failed: %v", err)
	}

	// Verify the data was parsed correctly
	if form.CurrentPassword != "myoldpassword123" {
		t.Errorf("Expected current_password 'myoldpassword123', got '%s'", form.CurrentPassword)
	}
	if form.NewPassword != "mynewpassword456" {
		t.Errorf("Expected new_password 'mynewpassword456', got '%s'", form.NewPassword)
	}

	t.Logf("Successfully decoded password form: CurrentPassword=%s, NewPassword=%s",
		form.CurrentPassword, form.NewPassword)
}
