package common

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with subdomain",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
		{
			name:    "invalid email format",
			email:   "invalid-email",
			wantErr: true,
		},
		{
			name:    "email without domain",
			email:   "test@",
			wantErr: true,
		},
		{
			name:    "email without @",
			email:   "testexample.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{
			name:    "valid phone with country code",
			phone:   "+1234567890",
			wantErr: false,
		},
		{
			name:    "valid phone without country code",
			phone:   "1234567890",
			wantErr: false,
		},
		{
			name:    "empty phone",
			phone:   "",
			wantErr: true,
		},
		{
			name:    "invalid phone format",
			phone:   "123",
			wantErr: true,
		},
		{
			name:    "phone with letters",
			phone:   "123abc456",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhone(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "valid password with special chars",
			password: "P@ssw0rd!",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "short password",
			password: "1234567",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		wantErr   bool
	}{
		{
			name:      "valid value",
			value:     "test",
			fieldName: "field",
			wantErr:   false,
		},
		{
			name:      "empty value",
			value:     "",
			fieldName: "field",
			wantErr:   true,
		},
		{
			name:      "whitespace only",
			value:     "   ",
			fieldName: "field",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequired(tt.value, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLength(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		min       int
		max       int
		wantErr   bool
	}{
		{
			name:      "valid length",
			value:     "test",
			fieldName: "field",
			min:       2,
			max:       10,
			wantErr:   false,
		},
		{
			name:      "too short",
			value:     "a",
			fieldName: "field",
			min:       2,
			max:       10,
			wantErr:   true,
		},
		{
			name:      "too long",
			value:     "verylongstring",
			fieldName: "field",
			min:       2,
			max:       10,
			wantErr:   true,
		},
		{
			name:      "exact min length",
			value:     "ab",
			fieldName: "field",
			min:       2,
			max:       10,
			wantErr:   false,
		},
		{
			name:      "exact max length",
			value:     "abcdefghij",
			fieldName: "field",
			min:       2,
			max:       10,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLength(tt.value, tt.fieldName, tt.min, tt.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
