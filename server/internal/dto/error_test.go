package dto_test

import (
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/dto"
)

func TestJSONAPI_Error_Validate(t *testing.T) {
	t.Parallel()

	t.Run("valid with title", func(t *testing.T) {
		t.Parallel()
		e := dto.ErrorObject{
			Title: "Invalid request",
		}

		err := e.Validate()
		if err != nil {
			t.Fatalf("expected valid error object, got %v", err)
		}
	})

	t.Run("valid with status", func(t *testing.T) {
		t.Parallel()

		e := dto.ErrorObject{
			Status: "400",
		}

		err := e.Validate()
		if err != nil {
			t.Fatalf("expected valid error object, got %v", err)
		}
	})

	t.Run("invalid empty error object", func(t *testing.T) {
		t.Parallel()

		e := dto.ErrorObject{}

		err := e.Validate()
		if err == nil {
			t.Fatalf("expected validation error, got nil")
		}
	})

	t.Run("valid with meta", func(t *testing.T) {
		t.Parallel()

		e := dto.ErrorObject{
			Meta: map[string]any{
				"debug": "info",
			},
		}

		err := e.Validate()
		if err != nil {
			t.Fatalf("expected valid error object, got %v", err)
		}
	})
}
