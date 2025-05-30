package entities

import (
	"time"

	"github.com/google/uuid"
)

// Tenant represents a tenant in the workplace management system
type Tenant struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Domain      string    `json:"domain" db:"domain"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	Plan        string    `json:"plan" db:"plan"`
	MaxUsers    int       `json:"max_users" db:"max_users"`
	MaxLocations int      `json:"max_locations" db:"max_locations"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// NewTenant creates a new tenant instance
func NewTenant(name, domain, plan string, maxUsers, maxLocations int) *Tenant {
	return &Tenant{
		ID:           uuid.New(),
		Name:         name,
		Domain:       domain,
		IsActive:     true,
		Plan:         plan,
		MaxUsers:     maxUsers,
		MaxLocations: maxLocations,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// Update updates tenant information
func (t *Tenant) Update(name, plan string, maxUsers, maxLocations int) {
	t.Name = name
	t.Plan = plan
	t.MaxUsers = maxUsers
	t.MaxLocations = maxLocations
	t.UpdatedAt = time.Now()
}

// Activate activates the tenant
func (t *Tenant) Activate() {
	t.IsActive = true
	t.UpdatedAt = time.Now()
}

// Deactivate deactivates the tenant
func (t *Tenant) Deactivate() {
	t.IsActive = false
	t.UpdatedAt = time.Now()
}