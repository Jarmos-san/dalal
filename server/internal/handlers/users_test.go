// Package handlers_test contains black-box tests for HTTP handlers.
//
// These tests validate the behaviour of the transport layer in isolation by:
//   - mocking service dependencies
//   - issuing HTTP requests using httptest utilities
//   - asserting on HTTP responses (status, headers, body)
//
// The goal is NOT to test business logic, but to ensure correct interaction
// between the handler and its dependencies, along with proper HTTP semantics.
package handlers_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/handlers"
	"github.com/Jarmos-san/arthika/server/internal/services"
)

// mockUserService is a test double that implements the UserService interface.
//
// It allows us to:
//   - control the output of the service layer
//   - simulate both success and failure scenarios deterministically
type mockUserService struct {
	user services.User
	err  error
}

// GetUser returns preconfigured values for testing.
//
// NOTE: This signature must match the service interface exactly.
// If the real service uses context, this should too.
func (m mockUserService) GetUser() (services.User, error) {
	return m.user, m.err
}

// TestGetUser_Success verifies the happy-path behaviour of the handler.
//
// It ensures that:
//   - the handler returns HTTP 200
//   - the correct Content-Type header is set
//   - the response body is valid JSON
//   - the payload matches expected data from the service
func TestGetUser_Success(t *testing.T) {
	t.Parallel()

	// Arrange: mock service is returning a valid user.
	mockSvc := mockUserService{ //nolint:exhaustruct
		user: services.User{
			Name: "Test User",
		},
	}

	logger := slog.Default()

	// Inject mock dependencies into the handler.
	handler := handlers.NewUserHandler(mockSvc, logger)

	// Create a test HTTP request
	req := httptest.NewRequestWithContext(
		context.TODO(), // context can be enriched later (timeouts, tracing, etc.)
		http.MethodGet,
		"/users/",
		nil,
	)

	// Recorder captures the HTTP response
	recorder := httptest.NewRecorder()

	// Act: invoke handler directly
	handler.GetUser(recorder, req)

	// Assert: HTTP status code
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	// Assert: Content-Type header (important for API correctness)
	if ct := recorder.Header().Get("Content-Type"); ct != "application/vnd.api+json" {
		t.Fatalf("unexpected content-type: %s", ct)
	}

	// Assert: Response body structure and content
	var resp handlers.UserResponse

	err := json.NewDecoder(recorder.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Name != "Test User" {
		t.Errorf("expected name %q, got %q", "Test User", resp.Name)
	}
}

// assertError is a simple sentinel error used to simulate failure scenarios.
type assertError struct{}

// Error implements the error interface.
func (assertError) Error() string { return "test error" }

// TestGetUser_Error verifies handler behaviour when the service returns an error.
//
// It ensures that:
//   - the handler responds with HTTP 500
//   - no successful payload is returned
//
// NOTE: The response body is intentionally not asserted here. That decision depends on
// the API error format (e.g., JSON:API error objects). Such assertions will be added
// once the error contract is formalised.
func TestGetUser_Error(t *testing.T) {
	t.Parallel()

	// Arrange: Mock service returning an error
	mockSvc := &mockUserService{ //nolint:exhaustruct
		err: assertError{},
	}

	logger := slog.Default()

	// Inject mock dependencies into the handler
	handler := handlers.NewUserHandler(mockSvc, logger)

	// Create a test HTTP request
	req := httptest.NewRequestWithContext(
		context.TODO(), // The context can be enriched in the future
		http.MethodGet,
		"/users/",
		nil,
	)

	// Recorder captures the HTTP response
	recorder := httptest.NewRecorder()

	// Act: Invoke handler directly
	handler.GetUser(recorder, req)

	// Assert: HTTP status code
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			recorder.Code,
		)
	}
}
