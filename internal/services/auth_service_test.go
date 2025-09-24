package services

import (
	"context"
	"testing"
	"time"

	"ndugu-backend/internal/common"
	"ndugu-backend/internal/models"
)

// MockUserRepository pour les tests
type MockUserRepository struct {
	users map[string]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, common.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, common.ErrUserNotFound
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	users := make([]*models.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

// MockOryClient pour les tests
type MockOryClient struct {
	users map[string]*models.User
}

func NewMockOryClient() *MockOryClient {
	return &MockOryClient{
		users: make(map[string]*models.User),
	}
}

func (m *MockOryClient) CreateUser(ctx context.Context, email, firstName, lastName string) (*models.User, error) {
	user := &models.User{
		ID:        "test-id-" + email,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Traits:    map[string]interface{}{"role": "user"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.users[user.ID] = user
	return user, nil
}

func (m *MockOryClient) GetUser(ctx context.Context, userID string) (*models.User, error) {
	user, exists := m.users[userID]
	if !exists {
		return nil, common.ErrUserNotFound
	}
	return user, nil
}

func (m *MockOryClient) ValidateSession(ctx context.Context, sessionToken string) (*models.Session, error) {
	return &models.Session{
		ID:        "session-id",
		UserID:    "test-user-id",
		Token:     sessionToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}, nil
}

func (m *MockOryClient) CreateOAuth2Client(ctx context.Context, clientID, clientName, redirectURI string) (*models.OAuth2Client, error) {
	return &models.OAuth2Client{
		ID:           clientID,
		Name:         clientName,
		Secret:       "test-secret",
		RedirectURIs: []string{redirectURI},
		GrantTypes:   []string{"authorization_code"},
		Scopes:       []string{"openid", "profile"},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (m *MockOryClient) CreatePermission(ctx context.Context, namespace, object, relation, subject string) error {
	return nil
}

func (m *MockOryClient) CheckPermission(ctx context.Context, namespace, object, relation, subject string) (bool, error) {
	return true, nil
}

func TestAuthService_CreateUser(t *testing.T) {
	// Arrange
	mockUserRepo := NewMockUserRepository()
	mockOryClient := NewMockOryClient()
	logger := common.NewSimpleLogger()
	authService := NewAuthService(mockUserRepo, mockOryClient, logger)

	ctx := context.Background()
	req := &models.CreateUserRequest{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Act
	response, err := authService.CreateUser(ctx, req)

	// Assert
	if err != nil {
		t.Errorf("CreateUser() error = %v, wantErr false", err)
	}
	if response == nil {
		t.Error("CreateUser() response is nil")
		return
	}
	if response.Email != req.Email {
		t.Errorf("CreateUser() email = %v, want %v", response.Email, req.Email)
	}
	if response.FirstName != req.FirstName {
		t.Errorf("CreateUser() firstName = %v, want %v", response.FirstName, req.FirstName)
	}
	if response.LastName != req.LastName {
		t.Errorf("CreateUser() lastName = %v, want %v", response.LastName, req.LastName)
	}
}

func TestAuthService_CreateUser_InvalidEmail(t *testing.T) {
	// Arrange
	mockUserRepo := NewMockUserRepository()
	mockOryClient := NewMockOryClient()
	logger := common.NewSimpleLogger()
	authService := NewAuthService(mockUserRepo, mockOryClient, logger)

	ctx := context.Background()
	req := &models.CreateUserRequest{
		Email:     "invalid-email",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Act
	_, err := authService.CreateUser(ctx, req)

	// Assert
	if err == nil {
		t.Error("CreateUser() expected error for invalid email")
	}
}

func TestAuthService_GetUser(t *testing.T) {
	// Arrange
	mockUserRepo := NewMockUserRepository()
	mockOryClient := NewMockOryClient()
	logger := common.NewSimpleLogger()
	authService := NewAuthService(mockUserRepo, mockOryClient, logger)

	ctx := context.Background()
	userID := "test-user-id"

	// Créer un utilisateur dans le mock
	user := &models.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockUserRepo.Create(ctx, user)

	// Act
	response, err := authService.GetUser(ctx, userID)

	// Assert
	if err != nil {
		t.Errorf("GetUser() error = %v, wantErr false", err)
	}
	if response == nil {
		t.Error("GetUser() response is nil")
		return
	}
	if response.ID != userID {
		t.Errorf("GetUser() id = %v, want %v", response.ID, userID)
	}
}

func TestAuthService_ValidateSession(t *testing.T) {
	// Arrange
	mockUserRepo := NewMockUserRepository()
	mockOryClient := NewMockOryClient()
	logger := common.NewSimpleLogger()
	authService := NewAuthService(mockUserRepo, mockOryClient, logger)

	ctx := context.Background()
	req := &models.ValidateSessionRequest{
		SessionToken: "valid-token",
	}

	// Act
	response, err := authService.ValidateSession(ctx, req)

	// Assert
	if err != nil {
		t.Errorf("ValidateSession() error = %v, wantErr false", err)
	}
	if response == nil {
		t.Error("ValidateSession() response is nil")
		return
	}
	if !response.Valid {
		t.Error("ValidateSession() valid = false, want true")
	}
}
