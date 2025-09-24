package repository

import (
	"context"
	"fmt"
	"time"

	"ndugu-backend/internal/auth"
)

// kratosClient implémentation du client Kratos
type kratosClient struct {
	client *auth.OryClient
}

// NewKratosClient crée une nouvelle instance du client Kratos
func NewKratosClient() KratosClient {
	return &kratosClient{
		client: auth.NewOryClient(),
	}
}

// CreateUser crée un utilisateur via Kratos
func (c *kratosClient) CreateUser(ctx context.Context, email, firstName, lastName string) (*KratosUser, error) {
	user, err := c.client.CreateUser(ctx, email, firstName, lastName)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de l'utilisateur: %w", err)
	}

	// Convertir en KratosUser
	return &KratosUser{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Traits:    user.Traits,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUser récupère un utilisateur via Kratos
func (c *kratosClient) GetUser(ctx context.Context, userID string) (*KratosUser, error) {
	user, err := c.client.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de l'utilisateur: %w", err)
	}

	// Convertir en KratosUser
	return &KratosUser{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Traits:    user.Traits,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// ValidateSession valide une session via Kratos
func (c *kratosClient) ValidateSession(ctx context.Context, sessionToken string) (*KratosSession, error) {
	session, err := c.client.ValidateSession(ctx, sessionToken)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la validation de la session: %w", err)
	}

	// Convertir en KratosSession
	return &KratosSession{
		Id: session.Id,
		Identity: KratosUser{
			ID:        session.Identity.Id,
			Email:     "", // Sera extrait depuis les traits
			Name:      make(map[string]interface{}),
			Traits:    session.Identity.Traits.(map[string]interface{}),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ExpiresAt: time.Now().Add(24 * time.Hour), // Exemple
	}, nil
}
