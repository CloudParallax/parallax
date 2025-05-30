package dto

import (
	"time"

	"github.com/google/uuid"
)

// CustomerResponse represents a customer in API responses
type CustomerResponse struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	Country     string    `json:"country"`
	PostalCode  string    `json:"postal_code"`
	CompanyName string    `json:"company_name"`
	JobTitle    string    `json:"job_title"`
	Notes       string    `json:"notes"`
	IsActive    bool      `json:"is_active"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateCustomerRequest represents a request to create a customer
type CreateCustomerRequest struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string `json:"last_name" validate:"required,min=1,max=50"`
	Email     string `json:"email" validate:"required,email,max=100"`
}

// UpdateCustomerRequest represents a request to update a customer
type UpdateCustomerRequest struct {
	FirstName   string   `json:"first_name" validate:"required,min=1,max=50"`
	LastName    string   `json:"last_name" validate:"required,min=1,max=50"`
	Email       string   `json:"email" validate:"required,email,max=100"`
	Phone       string   `json:"phone" validate:"omitempty,max=20"`
	Address     string   `json:"address" validate:"omitempty,max=200"`
	City        string   `json:"city" validate:"omitempty,max=50"`
	State       string   `json:"state" validate:"omitempty,max=50"`
	Country     string   `json:"country" validate:"omitempty,max=50"`
	PostalCode  string   `json:"postal_code" validate:"omitempty,max=20"`
	CompanyName string   `json:"company_name" validate:"omitempty,max=100"`
	JobTitle    string   `json:"job_title" validate:"omitempty,max=100"`
	Notes       string   `json:"notes" validate:"omitempty,max=1000"`
	Tags        []string `json:"tags" validate:"omitempty"`
}

// AddTagRequest represents a request to add a tag to a customer
type AddTagRequest struct {
	Tag string `json:"tag" validate:"required,min=1,max=50"`
}

// RemoveTagRequest represents a request to remove a tag from a customer
type RemoveTagRequest struct {
	Tag string `json:"tag" validate:"required,min=1,max=50"`
}
