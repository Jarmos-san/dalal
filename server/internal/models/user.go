// Package models define the core domain entities used across the application.
//
// Models represent the internal data structures of the system and should remain
// independent of transport-layer concerns such as HTTP request/response formats.
//
// These structs may contain sensitive or internal-only fields that must not be exposed
// directly to external clients.
package models

// User represents an application user within the system.
//
// It contains core identity attributes such as username and email, along with a
// password hash used for authentication.
//
// The `PasswordHash` field is intentionally excluded from JSON serialization to prevent
// accidental exposure in API responses.
type User struct {
	// ID is the unique identifier for the user.
	ID string `json:"id"`

	// Username is the display or login name of the user.
	Username string `json:"username"`

	// Email is the user's email address.
	Email string `json:"email"`

	// PasswordHash stores the hashed password. It is excluded from JSON output for
	// security reasons.
	PasswordHash string `json:"-"`
}
