package models

import (
	"time"
)

// User représente un utilisateur dans le système
type User struct {
	ID        string                 `json:"id" db:"id"`
	Email     string                 `json:"email" db:"email"`
	FirstName string                 `json:"firstName" db:"first_name"`
	LastName  string                 `json:"lastName" db:"last_name"`
	Traits    map[string]interface{} `json:"traits" db:"traits"`
	CreatedAt time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time              `json:"updatedAt" db:"updated_at"`
}

// CreateUserRequest représente la requête de création d'utilisateur
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required,min=2,max=50"`
	LastName  string `json:"lastName" validate:"required,min=2,max=50"`
}

// UpdateUserRequest représente la requête de mise à jour d'utilisateur
type UpdateUserRequest struct {
	ID        string `json:"id" validate:"required"`
	Email     string `json:"email" validate:"omitempty,email"`
	FirstName string `json:"firstName" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"lastName" validate:"omitempty,min=2,max=50"`
}

// UserResponse représente la réponse utilisateur
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToResponse convertit un User en UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
