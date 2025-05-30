package usecases

import (
	"context"

	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/cloudparallax/parallax/internal/domain/repositories"
	"github.com/google/uuid"
)

// TenantUseCase handles tenant business logic
type TenantUseCase struct {
	tenantRepo repositories.TenantRepository
}

// NewTenantUseCase creates a new tenant use case
func NewTenantUseCase(tenantRepo repositories.TenantRepository) *TenantUseCase {
	return &TenantUseCase{
		tenantRepo: tenantRepo,
	}
}

// CreateTenant creates a new tenant
func (uc *TenantUseCase) CreateTenant(ctx context.Context, name, domain, plan string, maxUsers, maxLocations int) (*entities.Tenant, error) {
	tenant := entities.NewTenant(name, domain, plan, maxUsers, maxLocations)
	
	err := uc.tenantRepo.Create(ctx, tenant)
	if err != nil {
		return nil, err
	}
	
	return tenant, nil
}

// GetTenant retrieves a tenant by ID
func (uc *TenantUseCase) GetTenant(ctx context.Context, id uuid.UUID) (*entities.Tenant, error) {
	return uc.tenantRepo.GetByID(ctx, id)
}

// GetTenantByDomain retrieves a tenant by domain
func (uc *TenantUseCase) GetTenantByDomain(ctx context.Context, domain string) (*entities.Tenant, error) {
	return uc.tenantRepo.GetByDomain(ctx, domain)
}

// GetAllTenants retrieves all tenants with pagination
func (uc *TenantUseCase) GetAllTenants(ctx context.Context, limit, offset int) ([]*entities.Tenant, error) {
	return uc.tenantRepo.GetAll(ctx, limit, offset)
}

// UpdateTenant updates tenant information
func (uc *TenantUseCase) UpdateTenant(ctx context.Context, id uuid.UUID, name, plan string, maxUsers, maxLocations int) (*entities.Tenant, error) {
	tenant, err := uc.tenantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	tenant.Update(name, plan, maxUsers, maxLocations)
	
	err = uc.tenantRepo.Update(ctx, tenant)
	if err != nil {
		return nil, err
	}
	
	return tenant, nil
}

// ActivateTenant activates a tenant
func (uc *TenantUseCase) ActivateTenant(ctx context.Context, id uuid.UUID) (*entities.Tenant, error) {
	tenant, err := uc.tenantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	tenant.Activate()
	
	err = uc.tenantRepo.Update(ctx, tenant)
	if err != nil {
		return nil, err
	}
	
	return tenant, nil
}

// DeactivateTenant deactivates a tenant
func (uc *TenantUseCase) DeactivateTenant(ctx context.Context, id uuid.UUID) (*entities.Tenant, error) {
	tenant, err := uc.tenantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	tenant.Deactivate()
	
	err = uc.tenantRepo.Update(ctx, tenant)
	if err != nil {
		return nil, err
	}
	
	return tenant, nil
}

// DeleteTenant deletes a tenant
func (uc *TenantUseCase) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	return uc.tenantRepo.Delete(ctx, id)
}

// GetActiveTenantCount returns the count of active tenants
func (uc *TenantUseCase) GetActiveTenantCount(ctx context.Context) (int, error) {
	return uc.tenantRepo.GetActiveCount(ctx)
}
