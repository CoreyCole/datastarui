package handler

import "net/http"

// ServiceHandler defines the interface that all Connect RPC services must implement.
// This enables automatic registration and consistent middleware application.
type ServiceHandler interface {
	// Handler returns the service path and configured http.Handler.
	// The returned string is the service path (e.g., "/com.datastarui.v1.auth.AuthService/").
	// The returned http.Handler includes all interceptors.
	Handler() (string, http.Handler)
}
