package grpc

// Types pour les messages gRPC AuthService
// Ces types correspondent aux messages définis dans ndugu.proto

// Messages pour la gestion des utilisateurs
type CreateUserRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type CreateUserResponse struct {
	UserId    string      `json:"userId"`
	Email     string      `json:"email"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	CreatedAt interface{} `json:"createdAt"` // Timestamp
}

type GetUserRequest struct {
	UserId string `json:"userId"`
}

type GetUserResponse struct {
	UserId    string      `json:"userId"`
	Email     string      `json:"email"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	CreatedAt interface{} `json:"createdAt"` // Timestamp
	UpdatedAt interface{} `json:"updatedAt"` // Timestamp
}

type ValidateSessionRequest struct {
	SessionToken string `json:"sessionToken"`
}

type ValidateSessionResponse struct {
	Valid     bool        `json:"valid"`
	UserId    string      `json:"userId"`
	Email     string      `json:"email"`
	ExpiresAt interface{} `json:"expiresAt"` // Timestamp
}

// Messages pour OAuth2
type CreateOAuth2ClientRequest struct {
	ClientId    string `json:"clientId"`
	ClientName  string `json:"clientName"`
	RedirectUri string `json:"redirectUri"`
}

type CreateOAuth2ClientResponse struct {
	ClientId     string   `json:"clientId"`
	ClientName   string   `json:"clientName"`
	ClientSecret string   `json:"clientSecret"`
	RedirectUris []string `json:"redirectUris"`
}

// Messages pour les permissions
type CreatePermissionRequest struct {
	Namespace string `json:"namespace"`
	Object    string `json:"object"`
	Relation  string `json:"relation"`
	Subject   string `json:"subject"`
}

type CreatePermissionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type CheckPermissionRequest struct {
	Namespace string `json:"namespace"`
	Object    string `json:"object"`
	Relation  string `json:"relation"`
	Subject   string `json:"subject"`
}

type CheckPermissionResponse struct {
	HasPermission bool   `json:"hasPermission"`
	Message       string `json:"message"`
}
