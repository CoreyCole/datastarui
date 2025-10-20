package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"

	"connectrpc.com/connect"

	"github.com/coreycole/datastarui/pkg/proto/com/datastarui/v1/auth/authconnect"
	loginv1 "github.com/coreycole/datastarui/pkg/proto/com/datastarui/v1/forms/login"

	"github.com/coreycole/datastarui/api/handler"
	"github.com/coreycole/datastarui/api/interceptors"
)

// Compile-time interface checks
var (
	_ authconnect.AuthServiceHandler = (*Service)(nil)
	_ handler.ServiceHandler          = (*Service)(nil)
)

// Service implements the AuthService RPC handlers.
// This is a placeholder implementation for demonstration purposes.
type Service struct {
	authconnect.UnimplementedAuthServiceHandler

	// In-memory user storage (placeholder)
	mu    sync.RWMutex
	users map[string]*loginv1.User // email -> user
}

// NewService creates a new auth service instance.
func NewService() *Service {
	s := &Service{
		users: make(map[string]*loginv1.User),
	}

	// Add test users for development
	s.users["test@test.com"] = &loginv1.User{
		Id:    "test-user-123",
		Email: "test@test.com",
		Name:  "Test User",
	}
	s.users["t@t.com"] = &loginv1.User{
		Id:    "test-user-456",
		Email: "t@t.com",
		Name:  "T",
	}

	return s
}

// Handler returns the service path and HTTP handler with interceptors.
func (s *Service) Handler() (string, http.Handler) {
	return authconnect.NewAuthServiceHandler(s,
		connect.WithInterceptors(
			interceptors.LoggingInterceptor(),
			interceptors.ValidationInterceptor(),
		))
}

// Login authenticates a user and returns a session token.
// This is a placeholder implementation that accepts any email/password
// as long as the user exists.
func (s *Service) Login(
	ctx context.Context,
	req *connect.Request[loginv1.LoginRequest],
) (*connect.Response[loginv1.LoginResponse], error) {
	email := req.Msg.Email
	password := req.Msg.Password

	s.mu.RLock()
	user, exists := s.users[email]
	s.mu.RUnlock()

	if !exists {
		return nil, connect.NewError(
			connect.CodeUnauthenticated,
			fmt.Errorf("user not found: %s", email),
		)
	}

	// Placeholder: accept any password
	// In production, verify hashed password
	_ = password

	// Generate session token
	token, err := generateToken()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&loginv1.LoginResponse{
		Token: token,
		User:  user,
	}), nil
}

// Signup creates a new user account.
// This is a placeholder implementation that stores users in memory.
func (s *Service) Signup(
	ctx context.Context,
	req *connect.Request[loginv1.SignupRequest],
) (*connect.Response[loginv1.SignupResponse], error) {
	email := req.Msg.Email
	password := req.Msg.Password
	name := req.Msg.Name

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if user already exists
	if _, exists := s.users[email]; exists {
		return nil, connect.NewError(
			connect.CodeAlreadyExists,
			fmt.Errorf("user already exists: %s", email),
		)
	}

	// Generate user ID
	userID, err := generateToken()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Create user (placeholder: not hashing password)
	_ = password // In production, hash the password

	user := &loginv1.User{
		Id:    userID,
		Email: email,
		Name:  name,
	}

	s.users[email] = user

	// Generate session token
	token, err := generateToken()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&loginv1.SignupResponse{
		Token: token,
		User:  user,
	}), nil
}

// generateToken creates a random session token.
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
