// Package dto contains data transfer objects used for communication between the API
// layer and external clients.
//
// DTOs define the shape of request and response payloads. They are designed to be
// serialized into formats such as JSON and must not contain business logic.
//
// The types in this package partially implement the JSON:API specification:
// https://jsonapi.org/
package dto

import (
	"errors"
)

// ResourceObject represents a resource object in a JSON:API document.
//
// A resource object MUST contain a "type" member and MAY contain an "id". Attributes
// and relationships are intentionally flexible to allow arbitrary payloads.
//
// Note:
// This implementation keeps attributes and relationships loosely typed for simplicity.
// Stronger typing can be introduced later using generics if needed.
type ResourceObject struct {
	Type          string            `json:"type"`
	ID            string            `json:"id,omitempty"`
	Attributes    map[string]any    `json:"attributes,omitempty"`
	Relationships map[string]any    `json:"relationships,omitempty"`
	Links         map[string]string `json:"links,omitempty"`
}

// JSONAPI describes the server's implementation details of the JSON:API specification.
//
// It MAY include a version string and optional meta information.
type JSONAPI struct {
	Version string         `json:"version,omitempty"`
	Meta    map[string]any `json:"meta,omitempty"`
}

// ErrorLinks contains URLs that provide more information about the error.
type ErrorLinks struct {
	// About links to details about this specific occurrence of the error.
	About string `json:"about,omitempty"`

	// Type links to a general description of this type of error.
	Type string `json:"type,omitempty"`
}

// ErrorSource identifies the source of an error.
//
// Only one of the fields should typically be set.
type ErrorSource struct {
	// Pointer is a JSON Pointer to the offending value in the request document.
	Pointer string `json:"pointer,omitempty"`

	// Parameter indicates which query parameter caused the error.
	Parameter string `json:"parameter,omitempty"`

	// Header indicates which HTTP header caused the error.
	Header string `json:"header,omitempty"`
}

// ErrorObject represents a single error in a JSON:API-compliant response.
//
// This is a simplified representation of the JSON:API error object. Only the most
// commonly used fields are included.
//
// Spec reference:
// https://jsonapi.org/format/#error-objects
type ErrorObject struct {
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `json:"id,omitempty"`

	// Links contain references to additional information about the error.
	Links *ErrorLinks `json:"links,omitempty"`

	// Status is the HTTP status code associated with the error, expressed as a string
	// (e.g., "404", "422", etc).
	Status string `json:"status,omitempty"`

	// Code is an application-specific error code.
	Code string `json:"code,omitempty"`

	// Title is a short, human-readable summary of the problem. It should remain stable
	// across occurrences of the same error.
	Title string `json:"title,omitempty"`

	// Detail provides a more detailed explanation of the error.
	Detail string `json:"detail,omitempty"`

	// Source identifies th source of the error (e.g., request body, query param).
	Source *ErrorSource `json:"source,omitempty"`

	// Meta contains non-standard  meta-information about the error.
	Meta map[string]any `json:"meta,omitempty"`
}

// NewErrorDocument constructs a JSON:API document containing one or more error objects.
//
// The returned document will have the "errors" member populated and will not include
// a "data" member, ensuring compliance with the rule that "data" and "errors"
// MUST NOT coexist.
//
// Each ErrorObject in the provided slice is expected to satisfy its own validation
// constraints. These will be verified when calling Document.Validate().
//
// Example JSON output:
//
//	{
//	  "errors": [
//	    { "status": "404", "title": "Not Found" }
//	  ]
//	}
func NewErrorDocument(errs []ErrorObject) Document[any] {
	return Document[any]{
		JSONAPI: &JSONAPI{
			Version: "1.0",
			Meta:    nil,
		},
		Errors:   errs,
		Links:    nil,
		Data:     nil,
		Meta:     nil,
		Included: nil,
	}
}

// Validate checks whether the ErrorObject satisfies the minimal requirements defined by
// the JSON:API specification.
//
// An error object MUST contain at least one of its defined members. This method ensures
// that at least one of the following fields is present:
//
//   - id
//   - links
//   - status
//   - code
//   - title
//   - detail
//   - source
//   - meta
//
// It does not perform deep validation of nested fields (e.g., Links or Source).
func (e ErrorObject) Validate() error {
	if e.ID == "" &&
		e.Links == nil &&
		e.Status == "" &&
		e.Code == "" &&
		e.Title == "" &&
		e.Detail == "" &&
		e.Source == nil &&
		len(e.Meta) == 0 {
		return errors.New( //nolint:err113
			"jsonapi: error object must contain at least one field",
		)
	}

	return nil
}

// Document represents a top-level JSON:API document.
//
// A document MUST contain at least one of the following top-level members:
//   - data: the document's primary data
//   - errors: an array of error objects
//   - meta: a meta object containing non-standard information
//
// Additionally:
//   - The members "data" and "errors" MUST NOT coexist in the same document.
//
// The Data field is generic and represents the JSON:API "data" member, which can be:
//   - a single resource object
//   - an array of resource objects
//   - null
//
// Data is defined as a pointer (*T) to distinguish between:
//   - nil pointer     -> "data" is omitted
//   - non-nil pointer -> "data" is present (including explicit null)
//
// This struct does not enforce all constraints at compile time. Hence, call Validate()
// to ensure the document complies with these rules.
type Document[T any] struct {
	JSONAPI  *JSONAPI          `json:"jsonapi,omitempty"`
	Links    map[string]string `json:"links,omitempty"`
	Data     *T                `json:"data,omitempty"`
	Errors   []ErrorObject     `json:"errors,omitempty"`
	Meta     map[string]any    `json:"meta,omitempty"`
	Included []ResourceObject  `json:"included,omitempty"`
}

// NewSingleDocument constructs a JSON:API document containing a single resource object.
//
// The returned document will have the "data" member set to the provided resource
// object, making it suitable for endpoints that return a single resource.
//
// The "errors" member will be empty, ensuring compliance with the rule that
// "data" and "errors" MUST NOT coexist.
//
// Example JSON output:
//
//	{
//	  "data": { "type": "...", "id": "..." }
//	}
func NewSingleDocument(resource ResourceObject) Document[ResourceObject] {
	return Document[ResourceObject]{
		JSONAPI: &JSONAPI{
			Version: "1.0",
			Meta:    nil,
		},
		Data:     &resource,
		Links:    nil,
		Errors:   nil,
		Meta:     nil,
		Included: nil,
	}
}

// NewCollection constructs a JSON:API document containing a collection of resource
// objects.
//
// The returned document will have the "data" member set to an array of resource
// objects,
// making it suitable for endpoints that return multiple resources.
//
// The "errors" member will be empty, ensuring compliance with the rule that
// "data" and "errors" MUST NOT coexist.
//
// Example JSON output:
//
//	{
//	  "data": [
//	    { "type": "...", "id": "..." },
//	    { "type": "...", "id": "..." }
//	  ]
//	}
func NewCollection(resource []ResourceObject) Document[[]ResourceObject] {
	return Document[[]ResourceObject]{
		JSONAPI: &JSONAPI{
			Version: "1.0",
			Meta:    nil,
		},
		Data:     &resource,
		Links:    nil,
		Errors:   nil,
		Meta:     nil,
		Included: nil,
	}
}

// NewNullDocument constructs a JSON:API document with an explicit null "data" member.
//
// This is typically used for endpoints where the primary data is intentionally absent,
// such as when a to-one relationship is empty.
//
// The "data" member will be present in the serialized output with a null value:
//
//	{
//	  "data": null
//	}
//
// This differs from omitting the "data" field entirely, which would be represented
// by a nil Data pointer.
//
// The "errors" member will be empty, ensuring compliance with the rule that
// "data" and "errors" MUST NOT coexist.
func NewNullDocument() Document[any] {
	var value any

	return Document[any]{
		JSONAPI: &JSONAPI{
			Version: "1.0",
			Meta:    nil,
		},
		Data:     &value,
		Links:    nil,
		Errors:   nil,
		Meta:     nil,
		Included: nil,
	}
}

// Validate checks whether the Document satisfies core JSON:API constraints.
//
// It enforces the following rules derived from the JSON:API specification:
//
//   - A document MUST contain at least one of the following top-level members:
//     "data", "errors", or "meta".
//   - The members "data" and "errors" MUST NOT coexist in the same document.
//
// In addition, this method performs shallow validation of error objects by
// invoking Validate() on each element in the Errors slice.
//
// The presence of "data" is determined by whether the Data pointer is non-nil.
// Note that a non-nil Data pointer may still represent a JSON null value.
//
// This method does not perform deep validation of:
//   - resource objects contained in "data"
//   - relationship structures
//   - the "included" member
//
// It is the caller's responsibility to ensure those structures are valid.
func (d Document[T]) Validate() error {
	if d.Data == nil && len(d.Errors) == 0 && len(d.Meta) == 0 {
		return errors.New( //nolint:err113
			"jsonapi: document must contain at least one of data, errors, or meta",
		)
	}

	if d.Data != nil && len(d.Errors) > 0 {
		return errors.New( //nolint:err113
			"jsonapi: document must not contain both data and errors",
		)
	}

	for _, e := range d.Errors {
		err := e.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
