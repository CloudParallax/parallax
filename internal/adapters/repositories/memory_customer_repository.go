package repositories

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/cloudparallax/parallax/internal/domain/repositories"
	"github.com/google/uuid"
)

// MemoryCustomerRepository implements CustomerRepository using in-memory storage
type MemoryCustomerRepository struct {
	customers map[uuid.UUID]*entities.Customer
	mutex     sync.RWMutex
}

// NewMemoryCustomerRepository creates a new memory-based customer repository
func NewMemoryCustomerRepository() repositories.CustomerRepository {
	return &MemoryCustomerRepository{
		customers: make(map[uuid.UUID]*entities.Customer),
	}
}

// Create stores a new customer
func (r *MemoryCustomerRepository) Create(ctx context.Context, customer *entities.Customer) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if email already exists for this tenant
	for _, c := range r.customers {
		if c.TenantID == customer.TenantID && c.Email == customer.Email {
			return errors.New("customer with this email already exists in tenant")
		}
	}

	r.customers[customer.ID] = customer
	return nil
}

// GetByID retrieves a customer by ID
func (r *MemoryCustomerRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Customer, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	customer, exists := r.customers[id]
	if !exists {
		return nil, errors.New("customer not found")
	}

	return customer, nil
}

// GetByTenantID retrieves customers by tenant ID with pagination
func (r *MemoryCustomerRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entities.Customer, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var customers []*entities.Customer
	count := 0

	for _, customer := range r.customers {
		if customer.TenantID == tenantID {
			if count < offset {
				count++
				continue
			}

			if len(customers) >= limit {
				break
			}

			customers = append(customers, customer)
			count++
		}
	}

	return customers, nil
}

// GetByEmail retrieves a customer by email within a tenant
func (r *MemoryCustomerRepository) GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*entities.Customer, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, customer := range r.customers {
		if customer.TenantID == tenantID && customer.Email == email {
			return customer, nil
		}
	}

	return nil, errors.New("customer not found")
}

// SearchByName searches customers by name within a tenant
func (r *MemoryCustomerRepository) SearchByName(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*entities.Customer, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var customers []*entities.Customer
	count := 0
	query = strings.ToLower(query)

	for _, customer := range r.customers {
		if customer.TenantID == tenantID {
			fullName := strings.ToLower(customer.GetFullName())
			if strings.Contains(fullName, query) || strings.Contains(strings.ToLower(customer.FirstName), query) || strings.Contains(strings.ToLower(customer.LastName), query) {
				if count < offset {
					count++
					continue
				}

				if len(customers) >= limit {
					break
				}

				customers = append(customers, customer)
				count++
			}
		}
	}

	return customers, nil
}

// GetByTags retrieves customers by tags within a tenant
func (r *MemoryCustomerRepository) GetByTags(ctx context.Context, tenantID uuid.UUID, tags []string, limit, offset int) ([]*entities.Customer, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var customers []*entities.Customer
	count := 0

	for _, customer := range r.customers {
		if customer.TenantID == tenantID && r.hasAnyTag(customer.Tags, tags) {
			if count < offset {
				count++
				continue
			}

			if len(customers) >= limit {
				break
			}

			customers = append(customers, customer)
			count++
		}
	}

	return customers, nil
}

// hasAnyTag checks if customer has any of the specified tags
func (r *MemoryCustomerRepository) hasAnyTag(customerTags, searchTags []string) bool {
	for _, searchTag := range searchTags {
		for _, customerTag := range customerTags {
			if customerTag == searchTag {
				return true
			}
		}
	}
	return false
}

// Update updates an existing customer
func (r *MemoryCustomerRepository) Update(ctx context.Context, customer *entities.Customer) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.customers[customer.ID]; !exists {
		return errors.New("customer not found")
	}

	r.customers[customer.ID] = customer
	return nil
}

// Delete removes a customer
func (r *MemoryCustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.customers[id]; !exists {
		return errors.New("customer not found")
	}

	delete(r.customers, id)
	return nil
}

// CountByTenantID returns the count of customers for a tenant
func (r *MemoryCustomerRepository) CountByTenantID(ctx context.Context, tenantID uuid.UUID) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := 0
	for _, customer := range r.customers {
		if customer.TenantID == tenantID {
			count++
		}
	}

	return count, nil
}
