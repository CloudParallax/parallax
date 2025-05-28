package repositories

import (
	"context"
	"github.com/cloudparallax/parallax/internal/domain/entities"
)

// CounterRepository defines the interface for counter data operations
type CounterRepository interface {
	// Create creates a new counter
	Create(ctx context.Context, counter *entities.Counter) error
	
	// GetByID retrieves a counter by its ID
	GetByID(ctx context.Context, id string) (*entities.Counter, error)
	
	// Update updates an existing counter
	Update(ctx context.Context, counter *entities.Counter) error
	
	// Delete deletes a counter by ID
	Delete(ctx context.Context, id string) error
	
	// GetAll retrieves all counters
	GetAll(ctx context.Context) ([]*entities.Counter, error)
	
	// Exists checks if a counter with the given ID exists
	Exists(ctx context.Context, id string) (bool, error)
}