package main

import (
	"context"

	"ndugu-backend/internal/common"
	v1 "ndugu-backend/internal/grpc/api/v1"
	"ndugu-backend/internal/models"
	"ndugu-backend/internal/services"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// gRPCServer encapsule le serveur gRPC
type gRPCServer struct {
	v1.UnimplementedAuthServiceServer
	authService services.AuthService
	logger      common.Logger
}

// NewGRPCServer crée une nouvelle instance du serveur gRPC
func NewGRPCServer(authService services.AuthService, logger common.Logger) *grpc.Server {
	server := grpc.NewServer()

	// Créer l'implémentation du service
	grpcService := &gRPCServer{
		authService: authService,
		logger:      logger,
	}

	// Enregistrer les services gRPC
	v1.RegisterAuthServiceServer(server, grpcService)

	// Activer la réflexion gRPC pour le débogage
	reflection.Register(server)

	return server
}

// ============================================================================
// AuthService Implementation
// ============================================================================

// CreateUser crée un nouvel utilisateur via Kratos
func (s *gRPCServer) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	s.logger.Info("gRPC CreateUser appelé", "email", req.Email, "firstName", req.FirstName, "lastName", req.LastName)

	// Validation des données d'entrée
	if req.Email == "" {
		s.logger.Error("Email requis manquant dans la requête gRPC CreateUser")
		return nil, status.Error(codes.InvalidArgument, "Email requis")
	}
	if req.FirstName == "" {
		s.logger.Error("Prénom requis manquant dans la requête gRPC CreateUser", "email", req.Email)
		return nil, status.Error(codes.InvalidArgument, "Prénom requis")
	}
	if req.LastName == "" {
		s.logger.Error("Nom requis manquant dans la requête gRPC CreateUser", "email", req.Email)
		return nil, status.Error(codes.InvalidArgument, "Nom requis")
	}

	// Créer la requête pour le service
	createReq := &models.CreateUserRequest{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// Appeler le service
	user, err := s.authService.CreateUser(ctx, createReq)
	if err != nil {
		s.logger.Error("Erreur lors de la création de l'utilisateur via gRPC", "email", req.Email, "error", err)
		return nil, status.Error(codes.Internal, "Erreur lors de la création de l'utilisateur")
	}

	s.logger.Info("Utilisateur créé avec succès via gRPC", "userId", user.ID, "email", user.Email)

	// Convertir en réponse gRPC
	return &v1.CreateUserResponse{
		UserId:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}

// GetUser récupère un utilisateur par son ID
func (s *gRPCServer) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	s.logger.Info("gRPC GetUser appelé pour ID: %s", req.UserId)

	// Validation des données d'entrée
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "ID utilisateur requis")
	}

	// Appeler le service
	user, err := s.authService.GetUser(ctx, req.UserId)
	if err != nil {
		s.logger.Error("Erreur lors de la récupération de l'utilisateur: %v", err)
		if isNotFoundError(err) {
			return nil, status.Error(codes.NotFound, "Utilisateur non trouvé")
		}
		return nil, status.Error(codes.Internal, "Erreur lors de la récupération de l'utilisateur")
	}

	// Convertir en réponse gRPC
	return &v1.GetUserResponse{
		UserId:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// ValidateSession valide une session
func (s *gRPCServer) ValidateSession(ctx context.Context, req *v1.ValidateSessionRequest) (*v1.ValidateSessionResponse, error) {
	s.logger.Info("gRPC ValidateSession appelé")

	// Validation des données d'entrée
	if req.SessionToken == "" {
		return nil, status.Error(codes.InvalidArgument, "Token de session requis")
	}

	// Créer la requête pour le service
	validateReq := &models.ValidateSessionRequest{
		SessionToken: req.SessionToken,
	}

	// Appeler le service
	session, err := s.authService.ValidateSession(ctx, validateReq)
	if err != nil {
		s.logger.Error("Erreur lors de la validation de la session: %v", err)
		return nil, status.Error(codes.Internal, "Erreur lors de la validation de la session")
	}

	// Convertir en réponse gRPC
	response := &v1.ValidateSessionResponse{
		Valid: session.Valid,
	}

	if session.Valid {
		response.UserId = session.UserID
		response.Email = session.Email
		response.ExpiresAt = timestamppb.New(session.ExpiresAt)
	}

	return response, nil
}

// CreateOAuth2Client crée un client OAuth2
func (s *gRPCServer) CreateOAuth2Client(ctx context.Context, req *v1.CreateOAuth2ClientRequest) (*v1.CreateOAuth2ClientResponse, error) {
	s.logger.Info("gRPC CreateOAuth2Client appelé pour client: %s", req.ClientId)

	// Validation des données d'entrée
	if req.ClientId == "" {
		return nil, status.Error(codes.InvalidArgument, "ID client requis")
	}
	if req.ClientName == "" {
		return nil, status.Error(codes.InvalidArgument, "Nom client requis")
	}
	if req.RedirectUri == "" {
		return nil, status.Error(codes.InvalidArgument, "URI de redirection requise")
	}

	// Créer la requête pour le service
	createReq := &models.CreateOAuth2ClientRequest{
		ID:          req.ClientId,
		Name:        req.ClientName,
		RedirectURI: req.RedirectUri,
	}

	// Appeler le service
	client, err := s.authService.CreateOAuth2Client(ctx, createReq)
	if err != nil {
		s.logger.Error("Erreur lors de la création du client OAuth2: %v", err)
		return nil, status.Error(codes.Internal, "Erreur lors de la création du client OAuth2")
	}

	// Convertir en réponse gRPC
	return &v1.CreateOAuth2ClientResponse{
		ClientId:     client.ID,
		ClientName:   client.Name,
		ClientSecret: client.Secret,
		RedirectUris: client.RedirectURIs,
	}, nil
}

// CreatePermission crée une permission
func (s *gRPCServer) CreatePermission(ctx context.Context, req *v1.CreatePermissionRequest) (*v1.CreatePermissionResponse, error) {
	s.logger.Info("gRPC CreatePermission appelé pour %s:%s#%s@%s", req.Namespace, req.Object, req.Relation, req.Subject)

	// Validation des données d'entrée
	if req.Namespace == "" {
		return nil, status.Error(codes.InvalidArgument, "Namespace requis")
	}
	if req.Object == "" {
		return nil, status.Error(codes.InvalidArgument, "Objet requis")
	}
	if req.Relation == "" {
		return nil, status.Error(codes.InvalidArgument, "Relation requise")
	}
	if req.Subject == "" {
		return nil, status.Error(codes.InvalidArgument, "Sujet requis")
	}

	// Créer la requête pour le service
	createReq := &models.CreatePermissionRequest{
		Namespace: req.Namespace,
		Object:    req.Object,
		Relation:  req.Relation,
		Subject:   req.Subject,
	}

	// Appeler le service
	permission, err := s.authService.CreatePermission(ctx, createReq)
	if err != nil {
		s.logger.Error("Erreur lors de la création de la permission: %v", err)
		return nil, status.Error(codes.Internal, "Erreur lors de la création de la permission")
	}

	// Convertir en réponse gRPC
	return &v1.CreatePermissionResponse{
		Success: permission.HasPermission,
		Message: permission.Message,
	}, nil
}

// CheckPermission vérifie une permission
func (s *gRPCServer) CheckPermission(ctx context.Context, req *v1.CheckPermissionRequest) (*v1.CheckPermissionResponse, error) {
	s.logger.Info("gRPC CheckPermission appelé pour %s:%s#%s@%s", req.Namespace, req.Object, req.Relation, req.Subject)

	// Validation des données d'entrée
	if req.Namespace == "" {
		return nil, status.Error(codes.InvalidArgument, "Namespace requis")
	}
	if req.Object == "" {
		return nil, status.Error(codes.InvalidArgument, "Objet requis")
	}
	if req.Relation == "" {
		return nil, status.Error(codes.InvalidArgument, "Relation requise")
	}
	if req.Subject == "" {
		return nil, status.Error(codes.InvalidArgument, "Sujet requis")
	}

	// Créer la requête pour le service
	checkReq := &models.CheckPermissionRequest{
		Namespace: req.Namespace,
		Object:    req.Object,
		Relation:  req.Relation,
		Subject:   req.Subject,
	}

	// Appeler le service
	permission, err := s.authService.CheckPermission(ctx, checkReq)
	if err != nil {
		s.logger.Error("Erreur lors de la vérification de la permission: %v", err)
		return nil, status.Error(codes.Internal, "Erreur lors de la vérification de la permission")
	}

	// Convertir en réponse gRPC
	return &v1.CheckPermissionResponse{
		HasPermission: permission.HasPermission,
		Message:       permission.Message,
	}, nil
}

// ============================================================================
// Helper Functions
// ============================================================================

// isNotFoundError vérifie si une erreur est de type "non trouvé"
func isNotFoundError(err error) bool {
	if appErr, ok := err.(*common.AppError); ok {
		return appErr.Code == common.ErrCodeNotFound ||
			appErr.Code == common.ErrCodeUserNotFound ||
			appErr.Code == common.ErrCodeCustomerNotFound
	}
	return false
}
