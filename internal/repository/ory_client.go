package repository

import (
	"context"
	"fmt"
	"time"

	"ndugu-backend/internal/common"
	"ndugu-backend/internal/models"
)

// oryClient implémentation du client Ory
type oryClient struct {
	kratosClient KratosClient
	hydraClient  HydraClient
	ketoClient   KetoClient
	logger       common.Logger
}

// NewOryClient crée une nouvelle instance du client Ory
func NewOryClient(logger common.Logger) OryClient {
	return &oryClient{
		kratosClient: NewKratosClient(),
		hydraClient:  NewHydraClient(),
		ketoClient:   NewKetoClient(),
		logger:       logger,
	}
}

// CreateUser crée un utilisateur via Kratos
func (c *oryClient) CreateUser(ctx context.Context, email, firstName, lastName string) (*models.User, error) {
	user, err := c.kratosClient.CreateUser(ctx, email, firstName, lastName)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de l'utilisateur: %w", err)
	}

	// Convertir en modèle User
	return &models.User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: firstName,
		LastName:  lastName,
		Traits:    user.Traits,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUser récupère un utilisateur via Kratos
func (c *oryClient) GetUser(ctx context.Context, userID string) (*models.User, error) {
	user, err := c.kratosClient.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de l'utilisateur: %w", err)
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

	return &models.User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: firstName,
		LastName:  lastName,
		Traits:    user.Traits,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// ValidateSession valide une session via Kratos
func (c *oryClient) ValidateSession(ctx context.Context, sessionToken string) (*models.Session, error) {
	session, err := c.kratosClient.ValidateSession(ctx, sessionToken)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la validation de la session: %w", err)
	}

	// Extraire l'email depuis l'identité
	if emailVal, ok := session.Identity.Traits["email"].(string); ok {
		_ = emailVal // Utilisé pour l'extraction future
	}

	return &models.Session{
		ID:        session.Id,
		UserID:    session.Identity.ID,
		Token:     sessionToken,
		Traits:    session.Identity.Traits,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Exemple
		CreatedAt: time.Now(),
	}, nil
}

// CreateOAuth2Client crée un client OAuth2 via Hydra
func (c *oryClient) CreateOAuth2Client(ctx context.Context, clientID, clientName, redirectURI string) (*models.OAuth2Client, error) {
	client, err := c.hydraClient.CreateOAuth2Client(ctx, clientID, clientName, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du client OAuth2: %w", err)
	}

	return &models.OAuth2Client{
		ID:           client.ID,
		Name:         client.Name,
		Secret:       client.Secret,
		RedirectURIs: client.RedirectURIs,
		GrantTypes:   client.GrantTypes,
		Scopes:       client.Scopes,
		CreatedAt:    client.CreatedAt,
		UpdatedAt:    client.UpdatedAt,
	}, nil
}

// CreatePermission crée une permission via Keto
func (c *oryClient) CreatePermission(ctx context.Context, namespace, object, relation, subject string) error {
	return c.ketoClient.CreatePermission(ctx, namespace, object, relation, subject)
}

// CheckPermission vérifie une permission via Keto
func (c *oryClient) CheckPermission(ctx context.Context, namespace, object, relation, subject string) (bool, error) {
	return c.ketoClient.CheckPermission(ctx, namespace, object, relation, subject)
}

// Types temporaires pour les clients Ory
type KratosUser struct {
	ID        string                 `json:"id"`
	Email     string                 `json:"email"`
	Name      map[string]interface{} `json:"name"`
	Traits    map[string]interface{} `json:"traits"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type KratosSession struct {
	Id        string     `json:"id"`
	Identity  KratosUser `json:"identity"`
	ExpiresAt time.Time  `json:"expires_at"`
}

type HydraOAuth2Client struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Secret       string    `json:"secret"`
	RedirectURIs []string  `json:"redirect_uris"`
	GrantTypes   []string  `json:"grant_types"`
	Scopes       []string  `json:"scopes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
