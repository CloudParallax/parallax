package entities

import (
	"time"

	"github.com/google/uuid"
)

// Location represents a physical location within a tenant
type Location struct {
	ID          uuid.UUID `json:"id" db:"id"`
	TenantID    uuid.UUID `json:"tenant_id" db:"tenant_id"`
	Name        string    `json:"name" db:"name"`
	Address     string    `json:"address" db:"address"`
	City        string    `json:"city" db:"city"`
	State       string    `json:"state" db:"state"`
	Country     string    `json:"country" db:"country"`
	PostalCode  string    `json:"postal_code" db:"postal_code"`
	Phone       string    `json:"phone" db:"phone"`
	Email       string    `json:"email" db:"email"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	Capacity    int       `json:"capacity" db:"capacity"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// NewLocation creates a new location instance
func NewLocation(tenantID uuid.UUID, name, address, city, state, country, postalCode string) *Location {
	return &Location{
		ID:         uuid.New(),
		TenantID:   tenantID,
		Name:       name,
		Address:    address,
		City:       city,
		State:      state,
		Country:    country,
		PostalCode: postalCode,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// Update updates location information
func (l *Location) Update(name, address, city, state, country, postalCode, phone, email, description string, capacity int) {
	l.Name = name
	l.Address = address
	l.City = city
	l.State = state
	l.Country = country
	l.PostalCode = postalCode
	l.Phone = phone
	l.Email = email
	l.Description = description
	l.Capacity = capacity
	l.UpdatedAt = time.Now()
}

// Activate activates the location
func (l *Location) Activate() {
	l.IsActive = true
	l.UpdatedAt = time.Now()
}

// Deactivate deactivates the location
func (l *Location) Deactivate() {
	l.IsActive = false
	l.UpdatedAt = time.Now()
}
