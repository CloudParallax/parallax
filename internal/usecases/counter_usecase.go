package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/cloudparallax/parallax/internal/domain/repositories"
)

// CounterUseCase defines the interface for counter business logic
type CounterUseCase interface {
	CreateCounter(ctx context.Context, input CreateCounterInput) (*entities.Counter, error)
	GetCounter(ctx context.Context, id string) (*entities.Counter, error)
	IncrementCounter(ctx context.Context, id string) (*entities.Counter, error)
	DecrementCounter(ctx context.Context, id string) (*entities.Counter, error)
	SetCounterValue(ctx context.Context, id string, value int) (*entities.Counter, error)
	DeleteCounter(ctx context.Context, id string) error
	GetAllCounters(ctx context.Context) ([]*entities.Counter, error)
	ResetCounter(ctx context.Context, id string, value int) (*entities.Counter, error)
}

// counterUseCase implements CounterUseCase interface
type counterUseCase struct {
	counterRepo repositories.CounterRepository
}

// NewCounterUseCase creates a new counter use case
func NewCounterUseCase(counterRepo repositories.CounterRepository) CounterUseCase {
	return &counterUseCase{
		counterRepo: counterRepo,
	}
}

// CreateCounterInput represents input for creating a counter
type CreateCounterInput struct {
	ID           string `json:"id" validate:"required,min=1,max=50"`
	InitialValue int    `json:"initial_value"`
	MinValue     int    `json:"min_value"`
	MaxValue     int    `json:"max_value" validate:"gtfield=MinValue"`
}

// CreateCounter creates a new counter
func (uc *counterUseCase) CreateCounter(ctx context.Context, input CreateCounterInput) (*entities.Counter, error) {
	// Validate input
	if err := uc.validateCreateInput(input); err != nil {
		return nil, err
	}

	// Check if counter already exists
	exists, err := uc.counterRepo.Exists(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check counter existence: %w", err)
	}
	if exists {
		return nil, errors.New("counter with this ID already exists")
	}

	// Create new counter entity
	counter, err := entities.NewCounter(input.ID, input.InitialValue, input.MinValue, input.MaxValue)
	if err != nil {
		return nil, fmt.Errorf("failed to create counter entity: %w", err)
	}

	// Save to repository
	if err := uc.counterRepo.Create(ctx, counter); err != nil {
		return nil, fmt.Errorf("failed to create counter: %w", err)
	}

	return counter, nil
}

// GetCounter retrieves a counter by its ID
func (uc *counterUseCase) GetCounter(ctx context.Context, id string) (*entities.Counter, error) {
	if id == "" {
		return nil, errors.New("counter ID is required")
	}

	counter, err := uc.counterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get counter: %w", err)
	}

	return counter, nil
}

// IncrementCounter increments a counter value
func (uc *counterUseCase) IncrementCounter(ctx context.Context, id string) (*entities.Counter, error) {
	if id == "" {
		return nil, errors.New("counter ID is required")
	}

	// Get counter
	counter, err := uc.counterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get counter: %w", err)
	}

	// Increment counter
	if err := counter.Increment(); err != nil {
		return nil, fmt.Errorf("failed to increment counter: %w", err)
	}

	// Save changes
	if err := uc.counterRepo.Update(ctx, counter); err != nil {
		return nil, fmt.Errorf("failed to update counter: %w", err)
	}

	return counter, nil
}

// DecrementCounter decrements a counter value
func (uc *counterUseCase) DecrementCounter(ctx context.Context, id string) (*entities.Counter, error) {
	if id == "" {
		return nil, errors.New("counter ID is required")
	}

	// Get counter
	counter, err := uc.counterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get counter: %w", err)
	}

	// Decrement counter
	if err := counter.Decrement(); err != nil {
		return nil, fmt.Errorf("failed to decrement counter: %w", err)
	}

	// Save changes
	if err := uc.counterRepo.Update(ctx, counter); err != nil {
		return nil, fmt.Errorf("failed to update counter: %w", err)
	}

	return counter, nil
}

// SetCounterValue sets a counter to a specific value
func (uc *counterUseCase) SetCounterValue(ctx context.Context, id string, value int) (*entities.Counter, error) {
	if id == "" {
		return nil, errors.New("counter ID is required")
	}

	// Get counter
	counter, err := uc.counterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get counter: %w", err)
	}

	// Set value
	if err := counter.SetValue(value); err != nil {
		return nil, fmt.Errorf("failed to set counter value: %w", err)
	}

	// Save changes
	if err := uc.counterRepo.Update(ctx, counter); err != nil {
		return nil, fmt.Errorf("failed to update counter: %w", err)
	}

	return counter, nil
}

// DeleteCounter deletes a counter
func (uc *counterUseCase) DeleteCounter(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("counter ID is required")
	}

	// Check if counter exists
	_, err := uc.counterRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get counter: %w", err)
	}

	// Delete counter
	if err := uc.counterRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete counter: %w", err)
	}

	return nil
}

// GetAllCounters retrieves all counters
func (uc *counterUseCase) GetAllCounters(ctx context.Context) ([]*entities.Counter, error) {
	counters, err := uc.counterRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all counters: %w", err)
	}

	return counters, nil
}

// ResetCounter resets a counter to a specific value
func (uc *counterUseCase) ResetCounter(ctx context.Context, id string, value int) (*entities.Counter, error) {
	return uc.SetCounterValue(ctx, id, value)
}

// Helper methods

func (uc *counterUseCase) validateCreateInput(input CreateCounterInput) error {
	if input.ID == "" {
		return errors.New("counter ID is required")
	}
	if input.MaxValue <= input.MinValue {
		return errors.New("max value must be greater than min value")
	}
	if input.InitialValue < input.MinValue || input.InitialValue > input.MaxValue {
		return errors.New("initial value must be within min and max bounds")
	}
	return nil
}