package repository

import (
	"context"
	"time"
)

// hydraClient implémentation du client Hydra
type hydraClient struct {
	// Temporairement vide - sera implémenté quand Hydra sera corrigé
}

// NewHydraClient crée une nouvelle instance du client Hydra
func NewHydraClient() HydraClient {
	return &hydraClient{}
}

// CreateOAuth2Client crée un client OAuth2 via Hydra
func (c *hydraClient) CreateOAuth2Client(ctx context.Context, clientID, clientName, redirectURI string) (*HydraOAuth2Client, error) {
	// Temporairement désactivé - à implémenter quand Hydra sera corrigé
	return &HydraOAuth2Client{
		ID:           clientID,
		Name:         clientName,
		Secret:       "temporairement-désactivé",
		RedirectURIs: []string{redirectURI},
		GrantTypes:   []string{"authorization_code", "refresh_token"},
		Scopes:       []string{"openid", "profile", "email"},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
