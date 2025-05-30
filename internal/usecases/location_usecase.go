package usecases

import (
	"context"
	"errors"

	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/cloudparallax/parallax/internal/domain/repositories"
	"github.com/google/uuid"
)

// LocationUseCase handles location business logic
type LocationUseCase struct {
	locationRepo repositories.LocationRepository
	tenantRepo   repositories.TenantRepository
}

// NewLocationUseCase creates a new location use case
func NewLocationUseCase(locationRepo repositories.LocationRepository, tenantRepo repositories.TenantRepository) *LocationUseCase {
	return &LocationUseCase{
		locationRepo: locationRepo,
		tenantRepo:   tenantRepo,
	}
}

// CreateLocation creates a new location
func (uc *LocationUseCase) CreateLocation(ctx context.Context, tenantID uuid.UUID, name, address, city, state, country, postalCode string) (*entities.Location, error) {
	// Verify tenant exists and is active
	tenant, err := uc.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, errors.New("tenant not found")
	}
	
	if !tenant.IsActive {
		return nil, errors.New("tenant is not active")
	}
	
	// Check location limit
	locationCount, err := uc.locationRepo.CountByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	
	if locationCount >= tenant.MaxLocations {
		return nil, errors.New("maximum number of locations reached for this tenant")
	}
	
	location := entities.NewLocation(tenantID, name, address, city, state, country, postalCode)
	
	err = uc.locationRepo.Create(ctx, location)
	if err != nil {
		return nil, err
	}
	
	return location, nil
}

// GetLocation retrieves a location by ID
func (uc *LocationUseCase) GetLocation(ctx context.Context, id uuid.UUID) (*entities.Location, error) {
	return uc.locationRepo.GetByID(ctx, id)
}

// GetLocationsByTenant retrieves locations by tenant ID with pagination
func (uc *LocationUseCase) GetLocationsByTenant(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entities.Location, error) {
	return uc.locationRepo.GetByTenantID(ctx, tenantID, limit, offset)
}

// GetActiveLocationsByTenant retrieves active locations by tenant ID
func (uc *LocationUseCase) GetActiveLocationsByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entities.Location, error) {
	return uc.locationRepo.GetActivByTenantID(ctx, tenantID)
}

// UpdateLocation updates location information
func (uc *LocationUseCase) UpdateLocation(ctx context.Context, id uuid.UUID, name, address, city, state, country, postalCode, phone, email, description string, capacity int) (*entities.Location, error) {
	location, err := uc.locationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	location.Update(name, address, city, state, country, postalCode, phone, email, description, capacity)
	
	err = uc.locationRepo.Update(ctx, location)
	if err != nil {
		return nil, err
	}
	
	return location, nil
}

// ActivateLocation activates a location
func (uc *LocationUseCase) ActivateLocation(ctx context.Context, id uuid.UUID) (*entities.Location, error) {
	location, err := uc.locationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	location.Activate()
	
	err = uc.locationRepo.Update(ctx, location)
	if err != nil {
		return nil, err
	}
	
	return location, nil
}

// DeactivateLocation deactivates a location
func (uc *LocationUseCase) DeactivateLocation(ctx context.Context, id uuid.UUID) (*entities.Location, error) {
	location, err := uc.locationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	location.Deactivate()
	
	err = uc.locationRepo.Update(ctx, location)
	if err != nil {
		return nil, err
	}
	
	return location, nil
}

// DeleteLocation deletes a location
func (uc *LocationUseCase) DeleteLocation(ctx context.Context, id uuid.UUID) error {
	return uc.locationRepo.Delete(ctx, id)
}

// GetLocationCount returns the count of locations for a tenant
func (uc *LocationUseCase) GetLocationCount(ctx context.Context, tenantID uuid.UUID) (int, error) {
	return uc.locationRepo.CountByTenantID(ctx, tenantID)
}
