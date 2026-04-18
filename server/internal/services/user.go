// Package services provides business logic abstractions and implementations.
//
// It defines service interfaces and their concrete implementations. Services should
// encapsulate domain logic and remain independent of transport layers.
package services

// User represents a basic user entity exposed by the service layer.
//
// This struct is intentionally minimal but can be extended with additional
// domain-specific fields (e.g., ID, Email, etc.) as needed.
type User struct {
	Name string `json:"name"`
}

// UserService defines the contract for user-related operations.
//
// This abstraction allows for multiple implementations (e.g., mock, DB-backed,
// API-based) and enables easier unit testing via dependency injection.
type UserService interface {
	// GetUser retrieves a user entity.
	//
	// Returns:
	//   - User: the retrieved user object
	//   - error: non-nil if retrieval fails
	GetUser() (User, error)
}

// userService is a concrete implementation of UserService.
//
// NOTE: This is an unexported struct to enforce usage through the constructor. This
// helps maintain control over instantiation and future extensibility.
type userService struct{}

// NewUserService constructs a new UserService implementation.
//
// Returns:
//   - *userService: a concrete implementation of UserService
//
// NOTE:
// Although it returns a concrete type, it is generally recommended to return the
// interface (UserService) unless you have a specific reason not to.
func NewUserService() *userService { //nolint:revive
	return &userService{}
}

// GetUser returns a static user.
//
// This is a stub implementation and should be replaced with real logic (e.g., database
// query, external API call, etc.) in production systems.
func (s userService) GetUser() (User, error) {
	return User{
		Name: "John Doe",
	}, nil
}
