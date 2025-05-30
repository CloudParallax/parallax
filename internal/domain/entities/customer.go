package entities

import (
	"time"

	"github.com/google/uuid"
)

// Customer represents a customer within a tenant
type Customer struct {
	ID          uuid.UUID `json:"id" db:"id"`
	TenantID    uuid.UUID `json:"tenant_id" db:"tenant_id"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	Address     string    `json:"address" db:"address"`
	City        string    `json:"city" db:"city"`
	State       string    `json:"state" db:"state"`
	Country     string    `json:"country" db:"country"`
	PostalCode  string    `json:"postal_code" db:"postal_code"`
	CompanyName string    `json:"company_name" db:"company_name"`
	JobTitle    string    `json:"job_title" db:"job_title"`
	Notes       string    `json:"notes" db:"notes"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	Tags        []string  `json:"tags" db:"tags"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// NewCustomer creates a new customer instance
func NewCustomer(tenantID uuid.UUID, firstName, lastName, email string) *Customer {
	return &Customer{
		ID:        uuid.New(),
		TenantID:  tenantID,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		IsActive:  true,
		Tags:      []string{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Update updates customer information
func (c *Customer) Update(firstName, lastName, email, phone, address, city, state, country, postalCode, companyName, jobTitle, notes string, tags []string) {
	c.FirstName = firstName
	c.LastName = lastName
	c.Email = email
	c.Phone = phone
	c.Address = address
	c.City = city
	c.State = state
	c.Country = country
	c.PostalCode = postalCode
	c.CompanyName = companyName
	c.JobTitle = jobTitle
	c.Notes = notes
	c.Tags = tags
	c.UpdatedAt = time.Now()
}

// GetFullName returns the customers full name
func (c *Customer) GetFullName() string {
	return c.FirstName + " " + c.LastName
}

// Activate activates the customer
func (c *Customer) Activate() {
	c.IsActive = true
	c.UpdatedAt = time.Now()
}

// Deactivate deactivates the customer
func (c *Customer) Deactivate() {
	c.IsActive = false
	c.UpdatedAt = time.Now()
}

// AddTag adds a tag to the customer
func (c *Customer) AddTag(tag string) {
	for _, existingTag := range c.Tags {
		if existingTag == tag {
			return // Tag already exists
		}
	}
	c.Tags = append(c.Tags, tag)
	c.UpdatedAt = time.Now()
}

// RemoveTag removes a tag from the customer
func (c *Customer) RemoveTag(tag string) {
	for i, existingTag := range c.Tags {
		if existingTag == tag {
			c.Tags = append(c.Tags[:i], c.Tags[i+1:]...)
			c.UpdatedAt = time.Now()
			return
		}
	}
}
