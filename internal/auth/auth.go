package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	hydra "github.com/ory/hydra-client-go/v2"
	// "github.com/ory/keto-client-go" // Temporairement commenté
	kratos "github.com/ory/kratos-client-go"
)

// OryClient encapsule les clients pour les services Ory
type OryClient struct {
	Kratos *kratos.APIClient
	Hydra  *hydra.APIClient
	// Keto   *keto.APIClient // Temporairement commenté
}

// NewOryClient crée une nouvelle instance du client Ory
func NewOryClient() *OryClient {
	// Configuration Kratos
	kratosConfig := kratos.NewConfiguration()
	kratosConfig.Servers = []kratos.ServerConfiguration{
		{
			URL: "http://kratos:4434", // Admin API
		},
	}
	kratosClient := kratos.NewAPIClient(kratosConfig)

	// Configuration Hydra
	hydraConfig := hydra.NewConfiguration()
	hydraConfig.Servers = []hydra.ServerConfiguration{
		{
			URL: "http://hydra:4445", // Admin API
		},
	}
	hydraClient := hydra.NewAPIClient(hydraConfig)

	// Configuration Keto - Temporairement commenté
	// ketoConfig := keto.NewConfiguration()
	// ketoConfig.Servers = []keto.ServerConfiguration{
	// 	{
	// 		URL: "http://keto:4467", // Write API
	// 	},
	// }
	// ketoClient := keto.NewAPIClient(ketoConfig)

	return &OryClient{
		Kratos: kratosClient,
		Hydra:  hydraClient,
		// Keto:   ketoClient, // Temporairement commenté
	}
}

// User représente un utilisateur dans le système
type User struct {
	ID        string                 `json:"id"`
	Email     string                 `json:"email"`
	Name      map[string]interface{} `json:"name"`
	Traits    map[string]interface{} `json:"traits"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// CreateUser crée un nouvel utilisateur via Kratos
func (c *OryClient) CreateUser(ctx context.Context, email, firstName, lastName string) (*User, error) {
	// Créer l'identité
	createIdentityBody := kratos.CreateIdentityBody{
		SchemaId: "default",
		Traits: map[string]interface{}{
			"email": email,
			"name": map[string]interface{}{
				"first": firstName,
				"last":  lastName,
			},
		},
	}

	identity, _, err := c.Kratos.IdentityApi.CreateIdentity(ctx).CreateIdentityBody(createIdentityBody).Execute()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de l'identité: %w", err)
	}

	// Convertir en User
	traits := identity.Traits.(map[string]interface{})
	user := &User{
		ID:     identity.Id,
		Email:  traits["email"].(string),
		Name:   traits["name"].(map[string]interface{}),
		Traits: traits,
	}

	if identity.CreatedAt != nil {
		user.CreatedAt = *identity.CreatedAt
	}
	if identity.UpdatedAt != nil {
		user.UpdatedAt = *identity.UpdatedAt
	}

	return user, nil
}

// GetUser récupère un utilisateur par son ID
func (c *OryClient) GetUser(ctx context.Context, userID string) (*User, error) {
	identity, _, err := c.Kratos.IdentityApi.GetIdentity(ctx, userID).Execute()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de l'utilisateur: %w", err)
	}

	traits := identity.Traits.(map[string]interface{})
	user := &User{
		ID:     identity.Id,
		Email:  traits["email"].(string),
		Name:   traits["name"].(map[string]interface{}),
		Traits: traits,
	}

	if identity.CreatedAt != nil {
		user.CreatedAt = *identity.CreatedAt
	}
	if identity.UpdatedAt != nil {
		user.UpdatedAt = *identity.UpdatedAt
	}

	return user, nil
}

// CreateOAuth2Client crée un client OAuth2 via Hydra - Temporairement commenté
// func (c *OryClient) CreateOAuth2Client(ctx context.Context, clientID, clientName, redirectURI string) (*hydra.OAuth2Client, error) {
// 	client := hydra.OAuth2Client{
// 		ClientId:                &clientID,
// 		ClientName:              &clientName,
// 		RedirectUris:            []string{redirectURI},
// 		GrantTypes:              []string{"authorization_code", "refresh_token"},
// 		ResponseTypes:           []string{"code"},
// 		Scope:                   hydra.PtrString("openid profile email"),
// 		TokenEndpointAuthMethod: hydra.PtrString("client_secret_basic"),
// 	}

// 	createdClient, _, err := c.Hydra.OAuth2Api.CreateOAuth2Client(ctx).OAuth2Client(client).Execute()
// 	if err != nil {
// 		return nil, fmt.Errorf("erreur lors de la création du client OAuth2: %w", err)
// 	}

// 	return createdClient, nil
// }

// CreatePermission crée une permission via Keto - Temporairement commenté
// func (c *OryClient) CreatePermission(ctx context.Context, namespace, object, relation, subject string) error {
// 	relationTuple := keto.RelationQuery{
// 		Namespace: &namespace,
// 		Object:    &object,
// 		Relation:  &relation,
// 		Subject:   &subject,
// 	}

// 	_, err := c.Keto.RelationTuplesApi.CreateRelationTuples(ctx).RelationQuery(relationTuple).Execute()
// 	if err != nil {
// 		return fmt.Errorf("erreur lors de la création de la permission: %w", err)
// 	}

// 	return nil
// }

// CheckPermission vérifie une permission via Keto - Temporairement commenté
// func (c *OryClient) CheckPermission(ctx context.Context, namespace, object, relation, subject string) (bool, error) {
// 	relationQuery := keto.RelationQuery{
// 		Namespace: &namespace,
// 		Object:    &object,
// 		Relation:  &relation,
// 		Subject:   &subject,
// 	}

// 	response, _, err := c.Keto.RelationTuplesApi.GetRelationTuples(ctx).RelationQuery(relationQuery).Execute()
// 	if err != nil {
// 		return false, fmt.Errorf("erreur lors de la vérification de la permission: %w", err)
// 	}

// 	return len(response.RelationTuples) > 0, nil
// }

// ValidateSession valide une session Kratos
func (c *OryClient) ValidateSession(ctx context.Context, sessionToken string) (*kratos.Session, error) {
	// Créer une requête avec le cookie de session
	req, err := http.NewRequestWithContext(ctx, "GET", "http://kratos:4433/sessions/whoami", nil)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la requête: %w", err)
	}

	req.Header.Set("Cookie", fmt.Sprintf("ory_kratos_session=%s", sessionToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la requête: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("session invalide: status %d", resp.StatusCode)
	}

	var session kratos.Session
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("erreur lors du décodage de la session: %w", err)
	}

	return &session, nil
}

// GetLoginFlow récupère un flux de connexion Kratos
func (c *OryClient) GetLoginFlow(ctx context.Context, flowID string) (*kratos.LoginFlow, error) {
	flow, _, err := c.Kratos.FrontendApi.GetLoginFlow(ctx).Id(flowID).Execute()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du flux de connexion: %w", err)
	}

	return flow, nil
}

// GetRegistrationFlow récupère un flux d'inscription Kratos
func (c *OryClient) GetRegistrationFlow(ctx context.Context, flowID string) (*kratos.RegistrationFlow, error) {
	flow, _, err := c.Kratos.FrontendApi.GetRegistrationFlow(ctx).Id(flowID).Execute()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du flux d'inscription: %w", err)
	}

	return flow, nil
}
