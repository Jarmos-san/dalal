package models_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/models"
)

func TestUser_JSONSerialization(t *testing.T) {
	t.Parallel()

	user := models.User{
		ID:           "123",
		Username:     "admin",
		Email:        "admin@example.com",
		PasswordHash: "secret-hash",
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal user: %v", err)
	}

	jsonStr := string(data)

	// Ensure expected fields exist
	if !strings.Contains(jsonStr, `"id":"123"`) {
		t.Errorf("expected id field in JSON, got %s", jsonStr)
	}

	if !strings.Contains(jsonStr, `"username":"admin"`) {
		t.Errorf("expected username field in JSON, got %s", jsonStr)
	}

	if !strings.Contains(jsonStr, `"email":"admin@example.com"`) {
		t.Errorf("expected email field in JSON, got %s", jsonStr)
	}

	// Ensure password hash is NOT present
	if strings.Contains(jsonStr, "secret-hash") {
		t.Errorf("password hash should not be serialized, got %s", jsonStr)
	}
}

func TestUser_JSONUnmarshal(t *testing.T) {
	t.Parallel()

	input := `{
		"id": "123",
		"username": "admin",
		"email": "admin@example.com"
	}`

	var user models.User

	err := json.Unmarshal([]byte(input), &user)
	if err != nil {
		t.Fatalf("failed to unmarshal user: %v", err)
	}

	if user.ID != "123" {
		t.Errorf("expected ID=123, got %s", user.ID)
	}

	if user.Username != "admin" {
		t.Errorf("expected Username=admin, got %s", user.Username)
	}

	if user.Email != "admin@example.com" {
		t.Errorf("expected Email=admin@example.com, got %s", user.Email)
	}

	// PasswordHash should remain empty
	if user.PasswordHash != "" {
		t.Errorf("expected PasswordHash to be empty, got %s", user.PasswordHash)
	}
}
