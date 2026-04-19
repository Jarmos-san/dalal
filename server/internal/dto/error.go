// Package dto contains data transfer objects used for communication between the APPI
// layer and external clients.
//
// DTOs define the shape of the request and response payloads and should not contain
// business logic. They are designed to be serialised into formats such as JSON.
package dto

// ErrorObject represents a single error in a JSON:API compliant response.
//
// It includes a status code, a short title, and a more detailed description of the
// error.
//
// This structure is intended to be used in API responses and follows a simplified
// version of thee JSON:API error object specification.
type ErrorObject struct {
	// Status is the HTTP status code associated with the error.
	Status string `json:"status"`

	// Title is a short, human-readable summary of the problem.
	Title string `json:"title"`

	// Detail provides a more detailed explanation of the error.
	Detail string `json:"detail"`
}

// ErrorResponse represents a JSON:API-style error response.
//
// It wraps one or more ErrorObject instances under the "errors" field, allowing
// multiple errors to be returned in a single response.
type ErrorResponse struct {
	Errors []ErrorObject `json:"errors"`
}
