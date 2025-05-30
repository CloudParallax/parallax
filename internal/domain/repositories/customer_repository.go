package repositories

import (
	"context"

	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/google/uuid"
)

// CustomerRepository defines the interface for customer data operations
type CustomerRepository interface {
	Create(ctx context.Context, customer *entities.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Customer, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entities.Customer, error)
	GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*entities.Customer, error)
	SearchByName(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*entities.Customer, error)
	GetByTags(ctx context.Context, tenantID uuid.UUID, tags []string, limit, offset int) ([]*entities.Customer, error)
	Update(ctx context.Context, customer *entities.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByTenantID(ctx context.Context, tenantID uuid.UUID) (int, error)
}
