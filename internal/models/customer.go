package models

import (
	"time"
)

// Customer représente un client dans le système
type Customer struct {
	ID          string    `json:"id" db:"id"`
	PhoneCode   string    `json:"phoneCode" db:"phone_code"`
	PhoneNumber string    `json:"phoneNumber" db:"phone_number"`
	Password    string    `json:"-" db:"password"` // Masqué dans les réponses JSON
	IsActive    bool      `json:"isActive" db:"is_active"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// CreateCustomerRequest représente la requête de création de client
type CreateCustomerRequest struct {
	PhoneCode   string `json:"phoneCode" validate:"required,min=1,max=5"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=7,max=15"`
	Password    string `json:"password" validate:"required,min=8"`
}

// CreateCustomerResponse représente la réponse de création de client
type CreateCustomerResponse struct {
	CustomerID string `json:"customerId"`
}

// CustomerResponse représente la réponse client
type CustomerResponse struct {
	ID          string    `json:"id"`
	PhoneCode   string    `json:"phoneCode"`
	PhoneNumber string    `json:"phoneNumber"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ToResponse convertit un Customer en CustomerResponse
func (c *Customer) ToResponse() *CustomerResponse {
	return &CustomerResponse{
		ID:          c.ID,
		PhoneCode:   c.PhoneCode,
		PhoneNumber: c.PhoneNumber,
		IsActive:    c.IsActive,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
