package repository

import (
	"context"
	"sync"
	"time"

	"ndugu-backend/internal/common"
	"ndugu-backend/internal/models"
)

// MockUserRepository implémentation mock du repository utilisateur
type MockUserRepository struct {
	users map[string]*models.User
	mutex sync.RWMutex
}

// NewMockUserRepository crée une nouvelle instance du mock repository
func NewMockUserRepository() UserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

// Create crée un utilisateur
func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Vérifier si l'utilisateur existe déjà
	for _, existingUser := range m.users {
		if existingUser.Email == user.Email {
			return common.ErrUserExists
		}
	}

	// Générer un ID si non fourni
	if user.ID == "" {
		user.ID = "user_" + time.Now().Format("20060102150405")
	}

	// Définir les timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	m.users[user.ID] = user
	return nil
}

// GetByID récupère un utilisateur par son ID
func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	user, exists := m.users[id]
	if !exists {
		return nil, common.ErrUserNotFound
	}

	// Retourner une copie pour éviter les modifications accidentelles
	userCopy := *user
	return &userCopy, nil
}

// GetByEmail récupère un utilisateur par son email
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, user := range m.users {
		if user.Email == email {
			// Retourner une copie pour éviter les modifications accidentelles
			userCopy := *user
			return &userCopy, nil
		}
	}

	return nil, common.ErrUserNotFound
}

// Update met à jour un utilisateur
func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Vérifier si l'utilisateur existe
	if _, exists := m.users[user.ID]; !exists {
		return common.ErrUserNotFound
	}

	// Mettre à jour le timestamp
	user.UpdatedAt = time.Now()

	m.users[user.ID] = user
	return nil
}

// Delete supprime un utilisateur
func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.users[id]; !exists {
		return common.ErrUserNotFound
	}

	delete(m.users, id)
	return nil
}

// List liste les utilisateurs avec pagination
func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	users := make([]*models.User, 0, len(m.users))
	count := 0
	skipped := 0

	for _, user := range m.users {
		if skipped < offset {
			skipped++
			continue
		}
		if count >= limit {
			break
		}

		// Retourner une copie pour éviter les modifications accidentelles
		userCopy := *user
		users = append(users, &userCopy)
		count++
	}

	return users, nil
}

// GetStats retourne des statistiques sur les utilisateurs
func (m *MockUserRepository) GetStats(ctx context.Context) (map[string]int, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := map[string]int{
		"total_users": len(m.users),
	}

	return stats, nil
}

// Clear supprime tous les utilisateurs (utile pour les tests)
func (m *MockUserRepository) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.users = make(map[string]*models.User)
}
