package repository

import (
	"context"

	"ndugu-backend/internal/models"
)

// UserRepository interface pour la gestion des utilisateurs
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
}

// CustomerRepository interface pour la gestion des clients
type CustomerRepository interface {
	Create(ctx context.Context, customer *models.Customer) error
	GetByID(ctx context.Context, id string) (*models.Customer, error)
	GetByPhone(ctx context.Context, phoneCode, phoneNumber string) (*models.Customer, error)
	Update(ctx context.Context, customer *models.Customer) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*models.Customer, error)
}

// OryClient interface pour les services Ory
type OryClient interface {
	CreateUser(ctx context.Context, email, firstName, lastName string) (*models.User, error)
	GetUser(ctx context.Context, userID string) (*models.User, error)
	ValidateSession(ctx context.Context, sessionToken string) (*models.Session, error)
	CreateOAuth2Client(ctx context.Context, clientID, clientName, redirectURI string) (*models.OAuth2Client, error)
	CreatePermission(ctx context.Context, namespace, object, relation, subject string) error
	CheckPermission(ctx context.Context, namespace, object, relation, subject string) (bool, error)
}

// Interfaces pour les clients Ory individuels
type KratosClient interface {
	CreateUser(ctx context.Context, email, firstName, lastName string) (*KratosUser, error)
	GetUser(ctx context.Context, userID string) (*KratosUser, error)
	ValidateSession(ctx context.Context, sessionToken string) (*KratosSession, error)
}

type HydraClient interface {
	CreateOAuth2Client(ctx context.Context, clientID, clientName, redirectURI string) (*HydraOAuth2Client, error)
}

type KetoClient interface {
	CreatePermission(ctx context.Context, namespace, object, relation, subject string) error
	CheckPermission(ctx context.Context, namespace, object, relation, subject string) (bool, error)
}
