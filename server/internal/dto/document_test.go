package dto_test

import (
	"encoding/json"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/dto"
)

func Test_JSONAPI_Document_Validate(t *testing.T) {
	t.Parallel()

	t.Run("valid single resource", func(t *testing.T) {
		t.Parallel()

		ro := dto.ResourceObject{Type: "users", ID: "1"}
		doc := dto.NewSingleDocument(ro)

		err := doc.Validate()
		if err != nil {
			t.Fatalf("expected valid document, got error: %v", err)
		}
	})

	t.Run("valid collection", func(t *testing.T) {
		t.Parallel()

		ros := []dto.ResourceObject{
			{Type: "users", ID: "1"},
			{Type: "users", ID: "2"},
		}

		doc := dto.NewCollection(ros)

		err := doc.Validate()
		if err != nil {
			t.Fatalf("expected valid document, got error: %v", err)
		}
	})

	t.Run("valid null data", func(t *testing.T) {
		t.Parallel()

		doc := dto.NewNullDocument()

		err := doc.Validate()
		if err != nil {
			t.Fatalf("expected valid null document, got error: %v", err)
		}
	})

	t.Run("valid error document", func(t *testing.T) {
		t.Parallel()

		doc := dto.NewErrorDocument([]dto.ErrorObject{
			{Title: "Not Found"},
		})

		err := doc.Validate()
		if err != nil {
			t.Fatalf("expected valid error document, got error: %v", err)
		}
	})

	t.Run("invalid: empty document", func(t *testing.T) {
		t.Parallel()

		doc := dto.Document[any]{}

		err := doc.Validate()
		if err == nil {
			t.Fatalf("expected validation error, got nil")
		}
	})

	t.Run("invalid: data and errors coexist", func(t *testing.T) {
		t.Parallel()

		ro := dto.ResourceObject{Type: "users"}
		doc := dto.Document[dto.ResourceObject]{
			Data: &ro,
			Errors: []dto.ErrorObject{
				{Title: "Bad Request"},
			},
		}

		err := doc.Validate()
		if err == nil {
			t.Fatalf("expected validation error, got nil")
		}
	})
}

func TestJSONAPI_Document_JSONSerialization(t *testing.T) {
	t.Parallel()

	t.Run("single resource serialization", func(t *testing.T) {
		t.Parallel()

		ro := dto.ResourceObject{Type: "users", ID: "1"}
		doc := dto.NewSingleDocument(ro)

		b, err := json.Marshal(doc)
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}

		if string(b) == "{}" {
			t.Fatalf("expected non-empty JSON")
		}
	})

	t.Run("collection serialization", func(t *testing.T) {
		t.Parallel()

		ros := []dto.ResourceObject{
			{Type: "users", ID: "1"},
		}
		doc := dto.NewCollection(ros)

		b, err := json.Marshal(doc)
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}

		if string(b) == "{}" {
			t.Fatalf("expected non-empty JSON")
		}
	})

	t.Run("null data serialization", func(t *testing.T) {
		t.Parallel()

		doc := dto.NewNullDocument()

		body, err := json.Marshal(doc)
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}

		// Expect "data":null to exist
		if !contains(string(body), `"data":null`) {
			t.Fatalf("expected data:null, got %s", string(body))
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (stringIndex(s, substr) >= 0)
}

func stringIndex(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}
