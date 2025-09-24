package services

import (
	"context"

	"ndugu-backend/internal/common"
	"ndugu-backend/internal/models"
	"ndugu-backend/internal/repository"
)

// AuthService interface pour le service d'authentification
type AuthService interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error)
	GetUser(ctx context.Context, userID string) (*models.UserResponse, error)
	ValidateSession(ctx context.Context, req *models.ValidateSessionRequest) (*models.ValidateSessionResponse, error)
	CreateOAuth2Client(ctx context.Context, req *models.CreateOAuth2ClientRequest) (*models.OAuth2Client, error)
	CreatePermission(ctx context.Context, req *models.CreatePermissionRequest) (*models.PermissionResponse, error)
	CheckPermission(ctx context.Context, req *models.CheckPermissionRequest) (*models.PermissionResponse, error)
}

// authService implémentation du service d'authentification
type authService struct {
	userRepo  repository.UserRepository
	oryClient repository.OryClient
	logger    common.Logger
}

// NewAuthService crée une nouvelle instance du service d'authentification
func NewAuthService(
	userRepo repository.UserRepository,
	oryClient repository.OryClient,
	logger common.Logger,
) AuthService {
	return &authService{
		userRepo:  userRepo,
		oryClient: oryClient,
		logger:    logger,
	}
}

// CreateUser crée un nouvel utilisateur
func (s *authService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	s.logger.Info("Début de création d'utilisateur", "email", req.Email, "firstName", req.FirstName, "lastName", req.LastName)

	// Validation
	if err := s.validateCreateUserRequest(req); err != nil {
		s.logger.Error("Validation échouée pour la création d'utilisateur", "email", req.Email, "error", err)
		return nil, err
	}

	// Vérifier si l'utilisateur existe déjà
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		s.logger.Warn("Tentative de création d'un utilisateur existant", "email", req.Email)
		return nil, common.ErrUserExists
	}

	s.logger.Debug("Utilisateur n'existe pas, création via Kratos", "email", req.Email)

	// Créer l'utilisateur via Kratos
	user, err := s.oryClient.CreateUser(ctx, req.Email, req.FirstName, req.LastName)
	if err != nil {
		s.logger.Error("Erreur lors de la création de l'utilisateur via Kratos", "email", req.Email, "error", err)
		return nil, common.NewAppError(common.ErrCodeKratosError, "Erreur lors de la création de l'utilisateur", err.Error())
	}

	s.logger.Info("Utilisateur créé avec succès via Kratos", "userId", user.ID, "email", user.Email)

	// Sauvegarder en base de données locale
	dbUser := &models.User{
		ID:        user.ID,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Traits:    user.Traits,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if err := s.userRepo.Create(ctx, dbUser); err != nil {
		s.logger.Error("Erreur lors de la sauvegarde de l'utilisateur en base locale", "userId", user.ID, "error", err)
		// Ne pas retourner d'erreur car l'utilisateur a été créé dans Kratos
	} else {
		s.logger.Debug("Utilisateur sauvegardé avec succès en base locale", "userId", user.ID)
	}

	return dbUser.ToResponse(), nil
}

// GetUser récupère un utilisateur par son ID
func (s *authService) GetUser(ctx context.Context, userID string) (*models.UserResponse, error) {
	s.logger.Info("Début de récupération d'utilisateur", "userId", userID)

	if userID == "" {
		s.logger.Error("ID utilisateur vide fourni")
		return nil, common.ErrInvalidInput
	}

	// Récupérer depuis la base de données locale
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Debug("Utilisateur non trouvé en base locale, tentative via Kratos", "userId", userID)
		// Essayer de récupérer depuis Kratos
		kratosUser, kratosErr := s.oryClient.GetUser(ctx, userID)
		if kratosErr != nil {
			s.logger.Error("Erreur lors de la récupération de l'utilisateur", "userId", userID, "localError", err, "kratosError", kratosErr)
			return nil, common.ErrUserNotFound
		}

		s.logger.Info("Utilisateur trouvé via Kratos", "userId", userID)

		// Convertir l'utilisateur Kratos
		user = &models.User{
			ID:        kratosUser.ID,
			Email:     kratosUser.Email,
			FirstName: kratosUser.FirstName,
			LastName:  kratosUser.LastName,
			Traits:    kratosUser.Traits,
			CreatedAt: kratosUser.CreatedAt,
			UpdatedAt: kratosUser.UpdatedAt,
		}
	}

	return user.ToResponse(), nil
}

// ValidateSession valide une session
func (s *authService) ValidateSession(ctx context.Context, req *models.ValidateSessionRequest) (*models.ValidateSessionResponse, error) {
	s.logger.Info("Début de validation de session", "sessionToken", req.SessionToken[:10]+"...")

	if req.SessionToken == "" {
		s.logger.Error("Token de session vide fourni")
		return nil, common.ErrInvalidInput
	}

	session, err := s.oryClient.ValidateSession(ctx, req.SessionToken)
	if err != nil {
		s.logger.Error("Erreur lors de la validation de session", "error", err)
		return &models.ValidateSessionResponse{
			Valid: false,
		}, nil
	}

	s.logger.Info("Session validée avec succès", "userId", session.UserID)

	// Extraire l'email depuis l'identité
	email := ""
	if emailVal, ok := session.Traits["email"].(string); ok {
		email = emailVal
	}

	return &models.ValidateSessionResponse{
		Valid:     true,
		UserID:    session.UserID,
		Email:     email,
		ExpiresAt: session.ExpiresAt,
	}, nil
}

// CreateOAuth2Client crée un client OAuth2
func (s *authService) CreateOAuth2Client(ctx context.Context, req *models.CreateOAuth2ClientRequest) (*models.OAuth2Client, error) {
	// Validation
	if err := s.validateCreateOAuth2ClientRequest(req); err != nil {
		return nil, err
	}

	// Créer le client via Hydra (temporairement désactivé)
	// client, err := s.oryClient.CreateOAuth2Client(ctx, req.ID, req.Name, req.RedirectURI)
	// if err != nil {
	// 	return nil, common.NewAppError(common.ErrCodeHydraError, "Erreur lors de la création du client OAuth2", err.Error())
	// }

	// Simulation temporaire
	return &models.OAuth2Client{
		ID:           req.ID,
		Name:         req.Name,
		Secret:       "temporairement-désactivé",
		RedirectURIs: []string{req.RedirectURI},
		GrantTypes:   []string{"authorization_code", "refresh_token"},
		Scopes:       []string{"openid", "profile", "email"},
	}, nil
}

// CreatePermission crée une permission
func (s *authService) CreatePermission(ctx context.Context, req *models.CreatePermissionRequest) (*models.PermissionResponse, error) {
	// Validation
	if err := s.validateCreatePermissionRequest(req); err != nil {
		return nil, err
	}

	// Créer la permission via Keto (temporairement désactivé)
	// err := s.oryClient.CreatePermission(ctx, req.Namespace, req.Object, req.Relation, req.Subject)
	// if err != nil {
	// 	return nil, common.NewAppError(common.ErrCodeKetoError, "Erreur lors de la création de la permission", err.Error())
	// }

	// Simulation temporaire
	return &models.PermissionResponse{
		HasPermission: true,
		Message:       "Permission créée (simulation - Keto temporairement désactivé)",
	}, nil
}

// CheckPermission vérifie une permission
func (s *authService) CheckPermission(ctx context.Context, req *models.CheckPermissionRequest) (*models.PermissionResponse, error) {
	// Validation
	if err := s.validateCheckPermissionRequest(req); err != nil {
		return nil, err
	}

	// Vérifier la permission via Keto (temporairement désactivé)
	// hasPermission, err := s.oryClient.CheckPermission(ctx, req.Namespace, req.Object, req.Relation, req.Subject)
	// if err != nil {
	// 	return nil, common.NewAppError(common.ErrCodeKetoError, "Erreur lors de la vérification de la permission", err.Error())
	// }

	// Simulation temporaire
	return &models.PermissionResponse{
		HasPermission: true,
		Message:       "Permission accordée (simulation - Keto temporairement désactivé)",
	}, nil
}

// Méthodes de validation privées
func (s *authService) validateCreateUserRequest(req *models.CreateUserRequest) error {
	if err := common.ValidateEmail(req.Email); err != nil {
		return err
	}
	if err := common.ValidateRequired(req.FirstName, "Prénom"); err != nil {
		return err
	}
	if err := common.ValidateRequired(req.LastName, "Nom"); err != nil {
		return err
	}
	return nil
}

func (s *authService) validateCreateOAuth2ClientRequest(req *models.CreateOAuth2ClientRequest) error {
	if err := common.ValidateRequired(req.ID, "ID du client"); err != nil {
		return err
	}
	if err := common.ValidateRequired(req.Name, "Nom du client"); err != nil {
		return err
	}
	if err := common.ValidateRequired(req.RedirectURI, "URI de redirection"); err != nil {
		return err
	}
	return nil
}

func (s *authService) validateCreatePermissionRequest(req *models.CreatePermissionRequest) error {
	if err := common.ValidateRequired(req.Namespace, "Namespace"); err != nil {
		return err
	}
	if err := common.ValidateRequired(req.Object, "Objet"); err != nil {
		return err
	}
	if err := common.ValidateRequired(req.Relation, "Relation"); err != nil {
		return err
	}
	if err := common.ValidateRequired(req.Subject, "Sujet"); err != nil {
		return err
	}
	return nil
}

func (s *authService) validateCheckPermissionRequest(req *models.CheckPermissionRequest) error {
	return s.validateCreatePermissionRequest(&models.CreatePermissionRequest{
		Namespace: req.Namespace,
		Object:    req.Object,
		Relation:  req.Relation,
		Subject:   req.Subject,
	})
}
