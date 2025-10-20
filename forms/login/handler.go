package login

import (
	"connectrpc.com/connect"
	"github.com/labstack/echo/v4"

	authconnect "github.com/coreycole/datastarui/pkg/proto/com/datastarui/v1/auth/authconnect"
	loginv1 "github.com/coreycole/datastarui/pkg/proto/com/datastarui/v1/forms/login"
	"github.com/coreycole/datastarui/utils"
)

type Handler struct {
	authService authconnect.AuthServiceClient
}

func NewHandler(authService authconnect.AuthServiceClient) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) HandleLoginForm(c echo.Context) error {
	// Parse form data - data comes nested under FormID ("login")
	var requestData struct {
		Login struct {
			Email    string `json:"email" form:"email"`
			Password string `json:"password" form:"password"`
		} `json:"login"`
	}

	if err := c.Bind(&requestData); err != nil {
		c.Logger().Errorf("Failed to bind form data: %v", err)
		return err
	}

	formData := requestData.Login
	c.Logger().Errorf("Bound form data - Email: '%s', Password: '%s'", formData.Email, formData.Password)

	// Call Connect RPC service internally
	_, err := h.authService.Login(
		c.Request().Context(),
		connect.NewRequest(&loginv1.LoginRequest{
			Email:    formData.Email,
			Password: formData.Password,
		}),
	)

	if err != nil {
		// Parse validation violations from Connect error
		violations := utils.ParseConnectViolations(err)

		// Debug: Log the error and violations
		c.Logger().Errorf("Login validation error: %v", err)
		c.Logger().Errorf("Parsed violations map: %+v", violations)
		c.Logger().Errorf("Email error: '%s', Password error: '%s'", violations["email"], violations["password"])
		c.Logger().Errorf("Form data - Email: '%s', Password: '%s'", formData.Email, formData.Password)

		// Re-render form component with errors, preserving user input
		args := LoginFormArgs{
			ID:            "login",
			Email:         formData.Email,    // Preserve user input
			Password:      formData.Password, // Preserve user input
			EmailError:    violations["email"],
			PasswordError: violations["password"],
		}
		c.Logger().Errorf("LoginFormArgs: %+v", args)
		component := LoginForm(args)

		// Set Datastar response headers for patching
		c.Response().Header().Set("Content-Type", "text/html")
		c.Response().Header().Set("datastar-selector", "#login-container")
		c.Response().Header().Set("datastar-mode", "outer")

		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	// Success case - redirect to dashboard
	c.Response().Header().Set("HX-Redirect", "/dashboard")
	return c.NoContent(200)
}
