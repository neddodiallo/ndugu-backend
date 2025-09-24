package models

import (
	"time"
)

// Session représente une session utilisateur
type Session struct {
	ID        string                 `json:"id" db:"id"`
	UserID    string                 `json:"userId" db:"user_id"`
	Token     string                 `json:"token" db:"token"`
	Traits    map[string]interface{} `json:"traits" db:"traits"`
	ExpiresAt time.Time              `json:"expiresAt" db:"expires_at"`
	CreatedAt time.Time              `json:"createdAt" db:"created_at"`
}

// ValidateSessionRequest représente la requête de validation de session
type ValidateSessionRequest struct {
	SessionToken string `json:"sessionToken" validate:"required"`
}

// ValidateSessionResponse représente la réponse de validation de session
type ValidateSessionResponse struct {
	Valid     bool      `json:"valid"`
	UserID    string    `json:"userId,omitempty"`
	Email     string    `json:"email,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

// OAuth2Client représente un client OAuth2
type OAuth2Client struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Secret       string    `json:"secret" db:"secret"`
	RedirectURIs []string  `json:"redirectUris" db:"redirect_uris"`
	GrantTypes   []string  `json:"grantTypes" db:"grant_types"`
	Scopes       []string  `json:"scopes" db:"scopes"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

// CreateOAuth2ClientRequest représente la requête de création de client OAuth2
type CreateOAuth2ClientRequest struct {
	ID          string `json:"id" validate:"required,min=3,max=50"`
	Name        string `json:"name" validate:"required,min=3,max=100"`
	RedirectURI string `json:"redirectUri" validate:"required,url"`
}

// Permission représente une permission dans le système
type Permission struct {
	ID        string    `json:"id" db:"id"`
	Namespace string    `json:"namespace" db:"namespace"`
	Object    string    `json:"object" db:"object"`
	Relation  string    `json:"relation" db:"relation"`
	Subject   string    `json:"subject" db:"subject"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// CreatePermissionRequest représente la requête de création de permission
type CreatePermissionRequest struct {
	Namespace string `json:"namespace" validate:"required,min=1,max=50"`
	Object    string `json:"object" validate:"required,min=1,max=100"`
	Relation  string `json:"relation" validate:"required,min=1,max=50"`
	Subject   string `json:"subject" validate:"required,min=1,max=100"`
}

// CheckPermissionRequest représente la requête de vérification de permission
type CheckPermissionRequest struct {
	Namespace string `json:"namespace" validate:"required,min=1,max=50"`
	Object    string `json:"object" validate:"required,min=1,max=100"`
	Relation  string `json:"relation" validate:"required,min=1,max=50"`
	Subject   string `json:"subject" validate:"required,min=1,max=100"`
}

// PermissionResponse représente la réponse de permission
type PermissionResponse struct {
	HasPermission bool   `json:"hasPermission"`
	Message       string `json:"message"`
}
