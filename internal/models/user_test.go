package models

import (
	"testing"
	"time"
)

func TestUser_ToResponse(t *testing.T) {
	// Arrange
	user := &User{
		ID:        "test-id",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Traits:    map[string]interface{}{"role": "user"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	response := user.ToResponse()

	// Assert
	if response.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, response.ID)
	}
	if response.Email != user.Email {
		t.Errorf("Expected Email %s, got %s", user.Email, response.Email)
	}
	if response.FirstName != user.FirstName {
		t.Errorf("Expected FirstName %s, got %s", user.FirstName, response.FirstName)
	}
	if response.LastName != user.LastName {
		t.Errorf("Expected LastName %s, got %s", user.LastName, response.LastName)
	}
	if !response.CreatedAt.Equal(user.CreatedAt) {
		t.Errorf("Expected CreatedAt %v, got %v", user.CreatedAt, response.CreatedAt)
	}
	if !response.UpdatedAt.Equal(user.UpdatedAt) {
		t.Errorf("Expected UpdatedAt %v, got %v", user.UpdatedAt, response.UpdatedAt)
	}
}

func TestCreateUserRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateUserRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreateUserRequest{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: false,
		},
		{
			name: "empty email",
			request: CreateUserRequest{
				Email:     "",
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			request: CreateUserRequest{
				Email:     "invalid-email",
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: true,
		},
		{
			name: "empty first name",
			request: CreateUserRequest{
				Email:     "test@example.com",
				FirstName: "",
				LastName:  "Doe",
			},
			wantErr: true,
		},
		{
			name: "empty last name",
			request: CreateUserRequest{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := validateCreateUserRequest(&tt.request)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCreateUserRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Fonction de validation pour les tests
func validateCreateUserRequest(req *CreateUserRequest) error {
	if req.Email == "" {
		return &ValidationError{Field: "email", Message: "email is required"}
	}
	// Validation basique de l'email
	if !contains(req.Email, "@") {
		return &ValidationError{Field: "email", Message: "invalid email format"}
	}
	if req.FirstName == "" {
		return &ValidationError{Field: "firstName", Message: "firstName is required"}
	}
	if req.LastName == "" {
		return &ValidationError{Field: "lastName", Message: "lastName is required"}
	}
	return nil
}

// Fonction utilitaire pour vérifier si une chaîne contient un caractère
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ValidationError pour les tests
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
