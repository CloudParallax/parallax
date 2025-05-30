package repositories

import (
	"context"
	"errors"
	"sync"

	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/cloudparallax/parallax/internal/domain/repositories"
	"github.com/google/uuid"
)

// MemoryTenantRepository implements TenantRepository using in-memory storage
type MemoryTenantRepository struct {
	tenants map[uuid.UUID]*entities.Tenant
	mutex   sync.RWMutex
}

// NewMemoryTenantRepository creates a new memory-based tenant repository
func NewMemoryTenantRepository() repositories.TenantRepository {
	return &MemoryTenantRepository{
		tenants: make(map[uuid.UUID]*entities.Tenant),
	}
}

// Create stores a new tenant
func (r *MemoryTenantRepository) Create(ctx context.Context, tenant *entities.Tenant) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if domain already exists
	for _, t := range r.tenants {
		if t.Domain == tenant.Domain {
			return errors.New("tenant with this domain already exists")
		}
	}

	r.tenants[tenant.ID] = tenant
	return nil
}

// GetByID retrieves a tenant by ID
func (r *MemoryTenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Tenant, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tenant, exists := r.tenants[id]
	if !exists {
		return nil, errors.New("tenant not found")
	}

	return tenant, nil
}

// GetByDomain retrieves a tenant by domain
func (r *MemoryTenantRepository) GetByDomain(ctx context.Context, domain string) (*entities.Tenant, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, tenant := range r.tenants {
		if tenant.Domain == domain {
			return tenant, nil
		}
	}

	return nil, errors.New("tenant not found")
}

// GetAll retrieves all tenants with pagination
func (r *MemoryTenantRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.Tenant, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var tenants []*entities.Tenant
	count := 0

	for _, tenant := range r.tenants {
		if count < offset {
			count++
			continue
		}

		if len(tenants) >= limit {
			break
		}

		tenants = append(tenants, tenant)
		count++
	}

	return tenants, nil
}

// Update updates an existing tenant
func (r *MemoryTenantRepository) Update(ctx context.Context, tenant *entities.Tenant) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tenants[tenant.ID]; !exists {
		return errors.New("tenant not found")
	}

	r.tenants[tenant.ID] = tenant
	return nil
}

// Delete removes a tenant
func (r *MemoryTenantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tenants[id]; !exists {
		return errors.New("tenant not found")
	}

	delete(r.tenants, id)
	return nil
}

// GetActiveCount returns the count of active tenants
func (r *MemoryTenantRepository) GetActiveCount(ctx context.Context) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := 0
	for _, tenant := range r.tenants {
		if tenant.IsActive {
			count++
		}
	}

	return count, nil
}
