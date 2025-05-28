package repositories

import (
	"context"
	"errors"
	"sync"

	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/cloudparallax/parallax/internal/domain/repositories"
)

// memoryCounterRepository implements CounterRepository interface using in-memory storage
type memoryCounterRepository struct {
	counters map[string]*entities.Counter
	mutex    sync.RWMutex
}

// NewMemoryCounterRepository creates a new in-memory counter repository
func NewMemoryCounterRepository() repositories.CounterRepository {
	repo := &memoryCounterRepository{
		counters: make(map[string]*entities.Counter),
		mutex:    sync.RWMutex{},
	}
	
	// Initialize with a default counter
	repo.initializeDefaultCounter()
	
	return repo
}

// Create creates a new counter
func (r *memoryCounterRepository) Create(ctx context.Context, counter *entities.Counter) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.counters[counter.ID]; exists {
		return errors.New("counter already exists")
	}
	
	// Create a copy to avoid external modifications
	counterCopy := *counter
	r.counters[counter.ID] = &counterCopy
	
	return nil
}

// GetByID retrieves a counter by its ID
func (r *memoryCounterRepository) GetByID(ctx context.Context, id string) (*entities.Counter, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	counter, exists := r.counters[id]
	if !exists {
		return nil, errors.New("counter not found")
	}
	
	// Return a copy to avoid external modifications
	counterCopy := *counter
	return &counterCopy, nil
}

// Update updates an existing counter
func (r *memoryCounterRepository) Update(ctx context.Context, counter *entities.Counter) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.counters[counter.ID]; !exists {
		return errors.New("counter not found")
	}
	
	// Create a copy to avoid external modifications
	counterCopy := *counter
	r.counters[counter.ID] = &counterCopy
	
	return nil
}

// Delete deletes a counter by ID
func (r *memoryCounterRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.counters[id]; !exists {
		return errors.New("counter not found")
	}
	
	delete(r.counters, id)
	return nil
}

// GetAll retrieves all counters
func (r *memoryCounterRepository) GetAll(ctx context.Context) ([]*entities.Counter, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var counters []*entities.Counter
	
	for _, counter := range r.counters {
		// Create a copy to avoid external modifications
		counterCopy := *counter
		counters = append(counters, &counterCopy)
	}
	
	return counters, nil
}

// Exists checks if a counter with the given ID exists
func (r *memoryCounterRepository) Exists(ctx context.Context, id string) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	_, exists := r.counters[id]
	return exists, nil
}

// Helper methods

func (r *memoryCounterRepository) initializeDefaultCounter() {
	defaultCounter, _ := entities.NewCounter("default", 1, 1, 100)
	r.counters["default"] = defaultCounter
}