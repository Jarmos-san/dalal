// Package handlers provides HTTP transport layer implementations.
//
// It is responsible for handling incoming HTTP requests, delegating business logic to
// the appropriate services, and formatting HTTP responses.
//
// Handlers should remain thin and only deal with HTTP-specific concerns such as:
//   - request parsing
//   - response encoding
//   - status code handling
//
// Business logic must be delegated to the service layer.
package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Jarmos-san/arthika/server/internal/services"
)

// UserHandler handles HTTP requests related to user resources.
//
// UserHandler acts as the transport layer for user-related endpoints. It delegates
// business logic to the injected UserService and is responsible for constructing HTTP
// responses based on the results.
//
// The handler is safe for concurrent use provided its dependencies are also
// concurrency-safe.
type UserHandler struct {
	service services.UserService
	logger  *slog.Logger
}

// NewUserHandler constructs a new UserHandler with its required dependencies.
//
// Parameters:
//   - service: provides user-related business logic
//   - logger:  used for structured logging within the handler
//
// The returned handler is ready to be registered with an HTTP router.
func NewUserHandler(service services.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// UserResponse represents the JSON response payload for a user resource.
//
// This type is a transport-layer DTO and defines the external API contract. It should
// remain decoupled from internal domain models to avoid leaking implementation details.
type UserResponse struct {
	Name string `json:"name"`
}

// GetUser handles HTTP GET requests for retrieving a user.
//
// It invokes the underlying UserService to fetch user data and returns a JSON response
// to the client.
//
// Success Response:
//   - Status: 200 OK
//   - Content-Type: application/vnd.api+json
//   - Body: JSON-encoded UserResponse
//
// Error Response:
//   - Status: 500 Internal Server Error
//   - Body: plain text error message
//
// Notes:
//   - The method uses a value receiver; this is acceptable since the handler
//     struct contains only references. Pointer receivers are still preferred
//     for consistency across handler methods.
//   - JSON encoding errors are logged but cannot alter the response once headers
//     have been written.
//   - A new logger instance is created during encoding failure, which is
//     inefficient and should be avoided in favor of the injected logger.
func (u UserHandler) GetUser(writer http.ResponseWriter, _ *http.Request) {
	user, serviceErr := u.service.GetUser()
	if serviceErr != nil {
		u.logger.Error(
			"failed to fetch user",
			slog.String("error", serviceErr.Error()),
		)
		http.Error(writer, "internal server error", http.StatusInternalServerError)

		return
	}

	resp := UserResponse{
		Name: user.Name,
	}

	writer.Header().Set("Content-Type", "application/vnd.api+json")
	writer.WriteHeader(http.StatusOK)

	encodingErr := json.NewEncoder(writer).Encode(resp)
	if encodingErr != nil {
		u.logger.Error(
			"JSON encoding failed",
			slog.String("error", encodingErr.Error()),
		)
	}
}

// CreateUser ...
func (u UserHandler) CreateUser( //nolint:funlen
	writer http.ResponseWriter,
	request *http.Request,
) {
	request.Body = http.MaxBytesReader(writer, request.Body, http.DefaultMaxHeaderBytes)

	formErr := request.ParseForm()
	if formErr != nil {
		u.logger.Error("failed to parse form", slog.String("error", formErr.Error()))

		http.Error(writer, "form parsing failed", http.StatusInternalServerError)

		return
	}

	username := request.FormValue("username")
	email := request.FormValue("email")
	password := request.FormValue("password")

	if username == "" {
		u.logger.Error("missing required field", slog.String("name", "username"))

		http.Error(
			writer,
			"missing required field: 'username'",
			http.StatusUnprocessableEntity,
		)

		return
	}

	if email == "" {
		u.logger.Error("missing required field", slog.String("name", "email"))

		http.Error(
			writer,
			"missing required field: 'email'",
			http.StatusUnprocessableEntity,
		)

		return
	}

	if password == "" {
		u.logger.Error("missing required field", slog.String("name", "password"))

		http.Error(
			writer,
			"missing required field: 'password'",
			http.StatusUnprocessableEntity,
		)

		return
	}

	resp, serviceErr := u.service.CreateUser(username, []byte(password))
	if serviceErr != nil {
		u.logger.Error("failed to create user")

		http.Error(writer, "internal server error", http.StatusInternalServerError)

		return
	}

	writer.Header().Set("Content-Type", "application/vnd.api+json")
	writer.WriteHeader(http.StatusCreated)

	encodingErr := json.NewEncoder(writer).Encode(resp)
	if encodingErr != nil {
		u.logger.Error(
			"JSON encoding failed",
			slog.String("error", encodingErr.Error()),
		)

		http.Error(writer, "internal server error", http.StatusInternalServerError)

		return
	}

	u.logger.Info(
		"successfully created new user",
		slog.String("id", resp.ID),
		slog.String("username", resp.Name),
	)
}
