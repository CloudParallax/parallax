package dto

import (
	"time"

	"github.com/cloudparallax/parallax/internal/domain/entities"
)

// CreateCounterRequest represents the request to create a counter
type CreateCounterRequest struct {
	ID           string `json:"id" validate:"required,min=1,max=50"`
	InitialValue int    `json:"initial_value"`
	MinValue     int    `json:"min_value"`
	MaxValue     int    `json:"max_value" validate:"gtfield=MinValue"`
}

// UpdateCounterValueRequest represents the request to update a counter value
type UpdateCounterValueRequest struct {
	Value int `json:"value" validate:"required"`
}

// CounterResponse represents the response for a counter
type CounterResponse struct {
	ID           string    `json:"id"`
	Value        int       `json:"value"`
	MinValue     int       `json:"min_value"`
	MaxValue     int       `json:"max_value"`
	IsAtMinimum  bool      `json:"is_at_minimum"`
	IsAtMaximum  bool      `json:"is_at_maximum"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CounterListResponse represents the response for a list of counters
type CounterListResponse struct {
	Counters []CounterResponse `json:"counters"`
}

// CounterOperationResponse represents the response for counter operations
type CounterOperationResponse struct {
	Counter   CounterResponse `json:"counter"`
	Operation string          `json:"operation"`
	Success   bool            `json:"success"`
	Message   string          `json:"message,omitempty"`
}

// ToCounterResponse converts a counter entity to response DTO
func ToCounterResponse(counter *entities.Counter) *CounterResponse {
	return &CounterResponse{
		ID:          counter.ID,
		Value:       counter.Value,
		MinValue:    counter.MinValue,
		MaxValue:    counter.MaxValue,
		IsAtMinimum: counter.IsAtMinimum(),
		IsAtMaximum: counter.IsAtMaximum(),
		CreatedAt:   counter.CreatedAt,
		UpdatedAt:   counter.UpdatedAt,
	}
}

// ToCounterListResponse converts a list of counter entities to list response DTO
func ToCounterListResponse(counters []*entities.Counter) *CounterListResponse {
	responses := make([]CounterResponse, len(counters))
	for i, counter := range counters {
		responses[i] = *ToCounterResponse(counter)
	}

	return &CounterListResponse{
		Counters: responses,
	}
}

// ToCounterOperationResponse converts a counter entity to operation response DTO
func ToCounterOperationResponse(counter *entities.Counter, operation string, success bool, message string) *CounterOperationResponse {
	return &CounterOperationResponse{
		Counter:   *ToCounterResponse(counter),
		Operation: operation,
		Success:   success,
		Message:   message,
	}
}