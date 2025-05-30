package repositories

import (
	"context"
	"errors"
	"sync"

	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/cloudparallax/parallax/internal/domain/repositories"
	"github.com/google/uuid"
)

// MemoryLocationRepository implements LocationRepository using in-memory storage
type MemoryLocationRepository struct {
	locations map[uuid.UUID]*entities.Location
	mutex     sync.RWMutex
}

// NewMemoryLocationRepository creates a new memory-based location repository
func NewMemoryLocationRepository() repositories.LocationRepository {
	return &MemoryLocationRepository{
		locations: make(map[uuid.UUID]*entities.Location),
	}
}

// Create stores a new location
func (r *MemoryLocationRepository) Create(ctx context.Context, location *entities.Location) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.locations[location.ID] = location
	return nil
}

// GetByID retrieves a location by ID
func (r *MemoryLocationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Location, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	location, exists := r.locations[id]
	if !exists {
		return nil, errors.New("location not found")
	}

	return location, nil
}

// GetByTenantID retrieves locations by tenant ID with pagination
func (r *MemoryLocationRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entities.Location, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var locations []*entities.Location
	count := 0

	for _, location := range r.locations {
		if location.TenantID == tenantID {
			if count < offset {
				count++
				continue
			}

			if len(locations) >= limit {
				break
			}

			locations = append(locations, location)
			count++
		}
	}

	return locations, nil
}

// GetActivByTenantID retrieves active locations by tenant ID
func (r *MemoryLocationRepository) GetActivByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entities.Location, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var locations []*entities.Location

	for _, location := range r.locations {
		if location.TenantID == tenantID && location.IsActive {
			locations = append(locations, location)
		}
	}

	return locations, nil
}

// Update updates an existing location
func (r *MemoryLocationRepository) Update(ctx context.Context, location *entities.Location) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.locations[location.ID]; !exists {
		return errors.New("location not found")
	}

	r.locations[location.ID] = location
	return nil
}

// Delete removes a location
func (r *MemoryLocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.locations[id]; !exists {
		return errors.New("location not found")
	}

	delete(r.locations, id)
	return nil
}

// CountByTenantID returns the count of locations for a tenant
func (r *MemoryLocationRepository) CountByTenantID(ctx context.Context, tenantID uuid.UUID) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := 0
	for _, location := range r.locations {
		if location.TenantID == tenantID {
			count++
		}
	}

	return count, nil
}
