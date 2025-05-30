package usecases

import (
	"context"
	"errors"

	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/cloudparallax/parallax/internal/domain/repositories"
	"github.com/google/uuid"
)

// CustomerUseCase handles customer business logic
type CustomerUseCase struct {
	customerRepo repositories.CustomerRepository
	tenantRepo   repositories.TenantRepository
}

// NewCustomerUseCase creates a new customer use case
func NewCustomerUseCase(customerRepo repositories.CustomerRepository, tenantRepo repositories.TenantRepository) *CustomerUseCase {
	return &CustomerUseCase{
		customerRepo: customerRepo,
		tenantRepo:   tenantRepo,
	}
}

// CreateCustomer creates a new customer
func (uc *CustomerUseCase) CreateCustomer(ctx context.Context, tenantID uuid.UUID, firstName, lastName, email string) (*entities.Customer, error) {
	// Verify tenant exists and is active
	tenant, err := uc.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, errors.New("tenant not found")
	}
	
	if !tenant.IsActive {
		return nil, errors.New("tenant is not active")
	}
	
	customer := entities.NewCustomer(tenantID, firstName, lastName, email)
	
	err = uc.customerRepo.Create(ctx, customer)
	if err != nil {
		return nil, err
	}
	
	return customer, nil
}

// GetCustomer retrieves a customer by ID
func (uc *CustomerUseCase) GetCustomer(ctx context.Context, id uuid.UUID) (*entities.Customer, error) {
	return uc.customerRepo.GetByID(ctx, id)
}

// GetCustomersByTenant retrieves customers by tenant ID with pagination
func (uc *CustomerUseCase) GetCustomersByTenant(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entities.Customer, error) {
	return uc.customerRepo.GetByTenantID(ctx, tenantID, limit, offset)
}

// GetCustomerByEmail retrieves a customer by email within a tenant
func (uc *CustomerUseCase) GetCustomerByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*entities.Customer, error) {
	return uc.customerRepo.GetByEmail(ctx, tenantID, email)
}

// SearchCustomersByName searches customers by name within a tenant
func (uc *CustomerUseCase) SearchCustomersByName(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*entities.Customer, error) {
	return uc.customerRepo.SearchByName(ctx, tenantID, query, limit, offset)
}

// GetCustomersByTags retrieves customers by tags within a tenant
func (uc *CustomerUseCase) GetCustomersByTags(ctx context.Context, tenantID uuid.UUID, tags []string, limit, offset int) ([]*entities.Customer, error) {
	return uc.customerRepo.GetByTags(ctx, tenantID, tags, limit, offset)
}

// UpdateCustomer updates customer information
func (uc *CustomerUseCase) UpdateCustomer(ctx context.Context, id uuid.UUID, firstName, lastName, email, phone, address, city, state, country, postalCode, companyName, jobTitle, notes string, tags []string) (*entities.Customer, error) {
	customer, err := uc.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	customer.Update(firstName, lastName, email, phone, address, city, state, country, postalCode, companyName, jobTitle, notes, tags)
	
	err = uc.customerRepo.Update(ctx, customer)
	if err != nil {
		return nil, err
	}
	
	return customer, nil
}

// ActivateCustomer activates a customer
func (uc *CustomerUseCase) ActivateCustomer(ctx context.Context, id uuid.UUID) (*entities.Customer, error) {
	customer, err := uc.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	customer.Activate()
	
	err = uc.customerRepo.Update(ctx, customer)
	if err != nil {
		return nil, err
	}
	
	return customer, nil
}

// DeactivateCustomer deactivates a customer
func (uc *CustomerUseCase) DeactivateCustomer(ctx context.Context, id uuid.UUID) (*entities.Customer, error) {
	customer, err := uc.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	customer.Deactivate()
	
	err = uc.customerRepo.Update(ctx, customer)
	if err != nil {
		return nil, err
	}
	
	return customer, nil
}

// AddCustomerTag adds a tag to a customer
func (uc *CustomerUseCase) AddCustomerTag(ctx context.Context, id uuid.UUID, tag string) (*entities.Customer, error) {
	customer, err := uc.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	customer.AddTag(tag)
	
	err = uc.customerRepo.Update(ctx, customer)
	if err != nil {
		return nil, err
	}
	
	return customer, nil
}

// RemoveCustomerTag removes a tag from a customer
func (uc *CustomerUseCase) RemoveCustomerTag(ctx context.Context, id uuid.UUID, tag string) (*entities.Customer, error) {
	customer, err := uc.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	customer.RemoveTag(tag)
	
	err = uc.customerRepo.Update(ctx, customer)
	if err != nil {
		return nil, err
	}
	
	return customer, nil
}

// DeleteCustomer deletes a customer
func (uc *CustomerUseCase) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
	return uc.customerRepo.Delete(ctx, id)
}

// GetCustomerCount returns the count of customers for a tenant
func (uc *CustomerUseCase) GetCustomerCount(ctx context.Context, tenantID uuid.UUID) (int, error) {
	return uc.customerRepo.CountByTenantID(ctx, tenantID)
}
