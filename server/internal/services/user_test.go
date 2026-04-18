package services_test

import (
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/services"
)

// TestNewUserService verifies that the constructor returns a non-nil service.
//
// This ensures that the service is properly initialized and ready for use.
func TestNewUserService(t *testing.T) {
	t.Parallel()

	svc := services.NewUserService()

	if svc == nil {
		t.Fatal("expected non-nil UserService, got nil")
	}
}

// TestUserService_GetUser validates the behavior of GetUser.
//
// Since the implementation is deterministic (returns a static user),
// we assert exact values.
func TestUserService_GetUser(t *testing.T) {
	t.Parallel()

	svc := services.NewUserService()

	user, err := svc.GetUser()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Validate returned user fields
	if user.Name != "John Doe" {
		t.Errorf("expected Name to be 'John Doe', got '%s'", user.Name)
	}
}

// TestUserService_InterfaceCompliance ensures that userService
// satisfies the UserService interface.
//
// This is a compile-time check and will fail if the implementation
// diverges from the interface contract.
func TestUserService_InterfaceCompliance(t *testing.T) {
	t.Parallel()

	var _ services.UserService
}
