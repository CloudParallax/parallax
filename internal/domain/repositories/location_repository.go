package repositories

import (
	"context"

	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/google/uuid"
)

// LocationRepository defines the interface for location data operations
type LocationRepository interface {
	Create(ctx context.Context, location *entities.Location) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Location, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entities.Location, error)
	GetActivByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entities.Location, error)
	Update(ctx context.Context, location *entities.Location) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByTenantID(ctx context.Context, tenantID uuid.UUID) (int, error)
}
