package repositories

import (
	"context"

	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/google/uuid"
)

// TenantRepository defines the interface for tenant data operations
type TenantRepository interface {
	Create(ctx context.Context, tenant *entities.Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Tenant, error)
	GetByDomain(ctx context.Context, domain string) (*entities.Tenant, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Tenant, error)
	Update(ctx context.Context, tenant *entities.Tenant) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetActiveCount(ctx context.Context) (int, error)
}
