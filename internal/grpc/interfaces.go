package grpc

import (
	"context"

	"google.golang.org/grpc"
)

// AuthServiceServerInterface définit l'interface pour le service d'authentification
type AuthServiceServerInterface interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
	GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error)
	ValidateSession(ctx context.Context, req *ValidateSessionRequest) (*ValidateSessionResponse, error)
	CreateOAuth2Client(ctx context.Context, req *CreateOAuth2ClientRequest) (*CreateOAuth2ClientResponse, error)
	CreatePermission(ctx context.Context, req *CreatePermissionRequest) (*CreatePermissionResponse, error)
	CheckPermission(ctx context.Context, req *CheckPermissionRequest) (*CheckPermissionResponse, error)
}

// RegisterAuthServiceServer enregistre le service d'authentification avec le serveur gRPC
func RegisterAuthServiceServer(s *grpc.Server, srv *AuthServiceServer) {
	// Cette fonction sera implémentée quand nous aurons le code généré par protoc
	// Pour l'instant, c'est un placeholder
}
