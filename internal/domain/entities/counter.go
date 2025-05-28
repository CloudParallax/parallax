package entities

import (
	"errors"
	"time"
)

// Counter represents a counter entity in the domain
type Counter struct {
	ID        string    `json:"id"`
	Value     int       `json:"value"`
	MinValue  int       `json:"min_value"`
	MaxValue  int       `json:"max_value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewCounter creates a new counter entity
func NewCounter(id string, initialValue, minValue, maxValue int) (*Counter, error) {
	if initialValue < minValue || initialValue > maxValue {
		return nil, errors.New("initial value must be within min and max bounds")
	}
	
	now := time.Now()
	return &Counter{
		ID:        id,
		Value:     initialValue,
		MinValue:  minValue,
		MaxValue:  maxValue,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Increment increases the counter value by 1
func (c *Counter) Increment() error {
	if c.Value >= c.MaxValue {
		return errors.New("counter value cannot exceed maximum")
	}
	c.Value++
	c.UpdatedAt = time.Now()
	return nil
}

// Decrement decreases the counter value by 1
func (c *Counter) Decrement() error {
	if c.Value <= c.MinValue {
		return errors.New("counter value cannot go below minimum")
	}
	c.Value--
	c.UpdatedAt = time.Now()
	return nil
}

// SetValue sets the counter to a specific value
func (c *Counter) SetValue(value int) error {
	if value < c.MinValue || value > c.MaxValue {
		return errors.New("value must be within min and max bounds")
	}
	c.Value = value
	c.UpdatedAt = time.Now()
	return nil
}

// IsAtMinimum checks if the counter is at its minimum value
func (c *Counter) IsAtMinimum() bool {
	return c.Value == c.MinValue
}

// IsAtMaximum checks if the counter is at its maximum value
func (c *Counter) IsAtMaximum() bool {
	return c.Value == c.MaxValue
}