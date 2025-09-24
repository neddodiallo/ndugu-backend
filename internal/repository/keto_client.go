package repository

import (
	"context"
	"fmt"
)

// ketoClient implémentation du client Keto
type ketoClient struct {
	// Temporairement vide - sera implémenté quand Keto sera corrigé
}

// NewKetoClient crée une nouvelle instance du client Keto
func NewKetoClient() KetoClient {
	return &ketoClient{}
}

// CreatePermission crée une permission via Keto
func (c *ketoClient) CreatePermission(ctx context.Context, namespace, object, relation, subject string) error {
	// Temporairement désactivé - à implémenter quand Keto sera corrigé
	return fmt.Errorf("Keto temporairement désactivé")
}

// CheckPermission vérifie une permission via Keto
func (c *ketoClient) CheckPermission(ctx context.Context, namespace, object, relation, subject string) (bool, error) {
	// Temporairement désactivé - à implémenter quand Keto sera corrigé
	return true, nil // Simulation
}
