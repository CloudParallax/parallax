package controllers

import (
	"strconv"
	"strings"

	"github.com/cloudparallax/parallax/internal/adapters/http/dto"
	"github.com/cloudparallax/parallax/internal/usecases"
	"github.com/cloudparallax/parallax/pkg/response"
	"github.com/gofiber/fiber/v3"
)

// CounterController handles HTTP requests for counter operations
type CounterController struct {
	counterUseCase usecases.CounterUseCase
}

// NewCounterController creates a new counter controller
func NewCounterController(counterUseCase usecases.CounterUseCase) *CounterController {
	return &CounterController{
		counterUseCase: counterUseCase,
	}
}

// CreateCounter handles creating a new counter
func (c *CounterController) CreateCounter(ctx fiber.Ctx) error {
	var req dto.CreateCounterRequest
	
	if err := response.ParseJSON(ctx, &req); err != nil {
		return response.Error(ctx, err)
	}

	input := usecases.CreateCounterInput{
		ID:           req.ID,
		InitialValue: req.InitialValue,
		MinValue:     req.MinValue,
		MaxValue:     req.MaxValue,
	}

	counter, err := c.counterUseCase.CreateCounter(ctx.Context(), input)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return response.Conflict(ctx, "Counter with this ID already exists")
		}
		return response.Error(ctx, err)
	}

	return response.Created(ctx, dto.ToCounterResponse(counter))
}

// GetCounter handles getting a counter by ID
func (c *CounterController) GetCounter(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return response.BadRequest(ctx, "Counter ID is required")
	}

	counter, err := c.counterUseCase.GetCounter(ctx.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(ctx, "Counter not found")
		}
		return response.Error(ctx, err)
	}

	return response.Success(ctx, dto.ToCounterResponse(counter))
}

// GetAllCounters handles getting all counters
func (c *CounterController) GetAllCounters(ctx fiber.Ctx) error {
	counters, err := c.counterUseCase.GetAllCounters(ctx.Context())
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, dto.ToCounterListResponse(counters))
}

// IncrementCounter handles incrementing a counter
func (c *CounterController) IncrementCounter(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return response.BadRequest(ctx, "Counter ID is required")
	}

	counter, err := c.counterUseCase.IncrementCounter(ctx.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(ctx, "Counter not found")
		}
		if strings.Contains(err.Error(), "cannot exceed maximum") {
			return response.UnprocessableEntity(ctx, "Counter cannot exceed maximum value")
		}
		return response.Error(ctx, err)
	}

	return response.Success(ctx, dto.ToCounterOperationResponse(counter, "increment", true, "Counter incremented successfully"))
}

// DecrementCounter handles decrementing a counter
func (c *CounterController) DecrementCounter(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return response.BadRequest(ctx, "Counter ID is required")
	}

	counter, err := c.counterUseCase.DecrementCounter(ctx.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(ctx, "Counter not found")
		}
		if strings.Contains(err.Error(), "cannot go below minimum") {
			return response.UnprocessableEntity(ctx, "Counter cannot go below minimum value")
		}
		return response.Error(ctx, err)
	}

	return response.Success(ctx, dto.ToCounterOperationResponse(counter, "decrement", true, "Counter decremented successfully"))
}

// SetCounterValue handles setting a counter to a specific value
func (c *CounterController) SetCounterValue(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return response.BadRequest(ctx, "Counter ID is required")
	}

	var req dto.UpdateCounterValueRequest
	
	if err := response.ParseJSON(ctx, &req); err != nil {
		return response.Error(ctx, err)
	}

	counter, err := c.counterUseCase.SetCounterValue(ctx.Context(), id, req.Value)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(ctx, "Counter not found")
		}
		if strings.Contains(err.Error(), "within min and max bounds") {
			return response.UnprocessableEntity(ctx, "Value must be within counter bounds")
		}
		return response.Error(ctx, err)
	}

	return response.Success(ctx, dto.ToCounterOperationResponse(counter, "set_value", true, "Counter value set successfully"))
}

// ResetCounter handles resetting a counter to a specific value
func (c *CounterController) ResetCounter(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return response.BadRequest(ctx, "Counter ID is required")
	}

	// Get reset value from query parameter or default to min value
	resetValue := 0
	if valueStr := ctx.Query("value"); valueStr != "" {
		if v, err := strconv.Atoi(valueStr); err == nil {
			resetValue = v
		}
	} else {
		// If no value provided, get the counter's min value
		counter, err := c.counterUseCase.GetCounter(ctx.Context(), id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return response.NotFound(ctx, "Counter not found")
			}
			return response.Error(ctx, err)
		}
		resetValue = counter.MinValue
	}

	counter, err := c.counterUseCase.ResetCounter(ctx.Context(), id, resetValue)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(ctx, "Counter not found")
		}
		if strings.Contains(err.Error(), "within min and max bounds") {
			return response.UnprocessableEntity(ctx, "Reset value must be within counter bounds")
		}
		return response.Error(ctx, err)
	}

	return response.Success(ctx, dto.ToCounterOperationResponse(counter, "reset", true, "Counter reset successfully"))
}

// DeleteCounter handles deleting a counter
func (c *CounterController) DeleteCounter(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return response.BadRequest(ctx, "Counter ID is required")
	}

	err := c.counterUseCase.DeleteCounter(ctx.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(ctx, "Counter not found")
		}
		return response.Error(ctx, err)
	}

	return response.NoContent(ctx)
}