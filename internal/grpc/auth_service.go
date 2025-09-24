package grpc

import (
	"context"
	"time"

	"ndugu-backend/internal/auth"
)

// AuthServiceServer implémente le service gRPC AuthService
type AuthServiceServer struct {
	oryClient *auth.OryClient
}

// NewAuthServiceServer crée une nouvelle instance du serveur AuthService
func NewAuthServiceServer(oryClient *auth.OryClient) *AuthServiceServer {
	return &AuthServiceServer{
		oryClient: oryClient,
	}
}

// CreateUser crée un nouvel utilisateur via Kratos
func (s *AuthServiceServer) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	user, err := s.oryClient.CreateUser(ctx, req.Email, req.FirstName, req.LastName)
	if err != nil {
		return nil, err
	}

	// Extraire le nom depuis les traits
	firstName := ""
	lastName := ""
	if name, ok := user.Name["first"].(string); ok {
		firstName = name
	}
	if name, ok := user.Name["last"].(string); ok {
		lastName = name
	}

	return &CreateUserResponse{
		UserId:    user.ID,
		Email:     user.Email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: user.CreatedAt,
	}, nil
}

// GetUser récupère un utilisateur par son ID
func (s *AuthServiceServer) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	user, err := s.oryClient.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	// Extraire le nom depuis les traits
	firstName := ""
	lastName := ""
	if name, ok := user.Name["first"].(string); ok {
		firstName = name
	}
	if name, ok := user.Name["last"].(string); ok {
		lastName = name
	}

	response := &GetUserResponse{
		UserId:    user.ID,
		Email:     user.Email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: user.CreatedAt,
	}

	if !user.UpdatedAt.IsZero() {
		response.UpdatedAt = user.UpdatedAt
	}

	return response, nil
}

// ValidateSession valide une session Kratos
func (s *AuthServiceServer) ValidateSession(ctx context.Context, req *ValidateSessionRequest) (*ValidateSessionResponse, error) {
	session, err := s.oryClient.ValidateSession(ctx, req.SessionToken)
	if err != nil {
		return &ValidateSessionResponse{
			Valid: false,
		}, nil
	}

	// Extraire l'email depuis l'identité
	email := ""
	if traits, ok := session.Identity.Traits.(map[string]interface{}); ok {
		if emailVal, ok := traits["email"].(string); ok {
			email = emailVal
		}
	}

	return &ValidateSessionResponse{
		Valid:     true,
		UserId:    session.Identity.Id,
		Email:     email,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Exemple
	}, nil
}

// CreateOAuth2Client crée un client OAuth2 via Hydra (temporairement désactivé)
func (s *AuthServiceServer) CreateOAuth2Client(ctx context.Context, req *CreateOAuth2ClientRequest) (*CreateOAuth2ClientResponse, error) {
	// Temporairement désactivé - à implémenter quand Hydra sera corrigé
	return &CreateOAuth2ClientResponse{
		ClientId:     req.ClientId,
		ClientName:   req.ClientName,
		ClientSecret: "temporairement-désactivé",
		RedirectUris: []string{req.RedirectUri},
	}, nil
}

// CreatePermission crée une permission via Keto (temporairement désactivé)
func (s *AuthServiceServer) CreatePermission(ctx context.Context, req *CreatePermissionRequest) (*CreatePermissionResponse, error) {
	// Temporairement désactivé - à implémenter quand Keto sera corrigé
	return &CreatePermissionResponse{
		Success: true,
		Message: "Permission créée (simulation - Keto temporairement désactivé)",
	}, nil
}

// CheckPermission vérifie une permission via Keto (temporairement désactivé)
func (s *AuthServiceServer) CheckPermission(ctx context.Context, req *CheckPermissionRequest) (*CheckPermissionResponse, error) {
	// Temporairement désactivé - à implémenter quand Keto sera corrigé
	return &CheckPermissionResponse{
		HasPermission: true,
		Message:       "Permission accordée (simulation - Keto temporairement désactivé)",
	}, nil
}
