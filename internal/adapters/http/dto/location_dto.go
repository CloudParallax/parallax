package dto

import (
	"time"

	"github.com/google/uuid"
)

// LocationResponse represents a location in API responses
type LocationResponse struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	Country     string    `json:"country"`
	PostalCode  string    `json:"postal_code"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	IsActive    bool      `json:"is_active"`
	Capacity    int       `json:"capacity"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateLocationRequest represents a request to create a location
type CreateLocationRequest struct {
	Name       string `json:"name" validate:"required,min=1,max=100"`
	Address    string `json:"address" validate:"required,min=1,max=200"`
	City       string `json:"city" validate:"required,min=1,max=50"`
	State      string `json:"state" validate:"required,min=1,max=50"`
	Country    string `json:"country" validate:"required,min=1,max=50"`
	PostalCode string `json:"postal_code" validate:"required,min=1,max=20"`
}

// UpdateLocationRequest represents a request to update a location
type UpdateLocationRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Address     string `json:"address" validate:"required,min=1,max=200"`
	City        string `json:"city" validate:"required,min=1,max=50"`
	State       string `json:"state" validate:"required,min=1,max=50"`
	Country     string `json:"country" validate:"required,min=1,max=50"`
	PostalCode  string `json:"postal_code" validate:"required,min=1,max=20"`
	Phone       string `json:"phone" validate:"omitempty,max=20"`
	Email       string `json:"email" validate:"omitempty,email,max=100"`
	Description string `json:"description" validate:"omitempty,max=500"`
	Capacity    int    `json:"capacity" validate:"omitempty,min=0"`
}
