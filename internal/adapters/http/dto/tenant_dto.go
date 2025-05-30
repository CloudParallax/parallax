package dto

import (
	"time"

	"github.com/google/uuid"
)

// TenantResponse represents a tenant in API responses
type TenantResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Domain       string    `json:"domain"`
	IsActive     bool      `json:"is_active"`
	Plan         string    `json:"plan"`
	MaxUsers     int       `json:"max_users"`
	MaxLocations int       `json:"max_locations"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateTenantRequest represents a request to create a tenant
type CreateTenantRequest struct {
	Name         string `json:"name" validate:"required,min=1,max=100"`
	Domain       string `json:"domain" validate:"required,min=1,max=50"`
	Plan         string `json:"plan" validate:"required,oneof=basic premium enterprise"`
	MaxUsers     int    `json:"max_users" validate:"required,min=1"`
	MaxLocations int    `json:"max_locations" validate:"required,min=1"`
}

// UpdateTenantRequest represents a request to update a tenant
type UpdateTenantRequest struct {
	Name         string `json:"name" validate:"required,min=1,max=100"`
	Plan         string `json:"plan" validate:"required,oneof=basic premium enterprise"`
	MaxUsers     int    `json:"max_users" validate:"required,min=1"`
	MaxLocations int    `json:"max_locations" validate:"required,min=1"`
}
