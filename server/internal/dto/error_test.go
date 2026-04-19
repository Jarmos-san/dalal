package dto_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/dto"
)

func TestErrorResponse_JSONSerialization(t *testing.T) {
	t.Parallel()

	resp := dto.ErrorResponse{
		Errors: []dto.ErrorObject{
			{
				Status: "400",
				Title:  "Bad Request",
				Detail: "invalid input",
			},
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal error response: %v", err)
	}

	jsonStr := string(data)

	// Ensure top-level structure
	if !strings.Contains(jsonStr, `"errors"`) {
		t.Errorf("expected 'errors' field in JSON, got %s", jsonStr)
	}

	// Ensure nested fields
	if !strings.Contains(jsonStr, `"status":"400"`) {
		t.Errorf("expected status field, got %s", jsonStr)
	}

	if !strings.Contains(jsonStr, `"title":"Bad Request"`) {
		t.Errorf("expected title field, got %s", jsonStr)
	}

	if !strings.Contains(jsonStr, `"detail":"invalid input"`) {
		t.Errorf("expected detail field, got %s", jsonStr)
	}
}

func TestErrorResponse_JSONUnmarshal(t *testing.T) {
	t.Parallel()

	input := `{
		"errors": [
			{
				"status": "400",
				"title": "Bad Request",
				"detail": "invalid input"
			}
		]
	}`

	var resp dto.ErrorResponse

	err := json.Unmarshal([]byte(input), &resp)
	if err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if len(resp.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(resp.Errors))
	}

	errObj := resp.Errors[0]

	if errObj.Status != "400" {
		t.Errorf("expected status=400, got %s", errObj.Status)
	}

	if errObj.Title != "Bad Request" {
		t.Errorf("expected title='Bad Request', got %s", errObj.Title)
	}

	if errObj.Detail != "invalid input" {
		t.Errorf("expected detail='invalid input', got %s", errObj.Detail)
	}
}
