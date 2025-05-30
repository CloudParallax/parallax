package controllers

import (
	"strconv"

	"github.com/cloudparallax/parallax/internal/adapters/http/dto"
	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/cloudparallax/parallax/internal/usecases"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// LocationController handles location HTTP requests
type LocationController struct {
	locationUseCase *usecases.LocationUseCase
}

// NewLocationController creates a new location controller
func NewLocationController(locationUseCase *usecases.LocationUseCase) *LocationController {
	return &LocationController{
		locationUseCase: locationUseCase,
	}
}

// CreateLocation creates a new location
func (lc *LocationController) CreateLocation(c fiber.Ctx) error {
	tenantIDParam := c.Params("tenantId")
	tenantID, err := uuid.Parse(tenantIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid tenant ID",
			},
		})
	}

	var req dto.CreateLocationRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid request body",
			},
		})
	}

	location, err := lc.locationUseCase.CreateLocation(c.Context(), tenantID, req.Name, req.Address, req.City, req.State, req.Country, req.PostalCode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	response := lc.toLocationResponse(location)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetLocation retrieves a location by ID
func (lc *LocationController) GetLocation(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid location ID",
			},
		})
	}

	location, err := lc.locationUseCase.GetLocation(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusNotFound,
				"message": "Location not found",
			},
		})
	}

	response := lc.toLocationResponse(location)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetLocationsByTenant retrieves locations by tenant ID with pagination
func (lc *LocationController) GetLocationsByTenant(c fiber.Ctx) error {
	tenantIDParam := c.Params("tenantId")
	tenantID, err := uuid.Parse(tenantIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid tenant ID",
			},
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	locations, err := lc.locationUseCase.GetLocationsByTenant(c.Context(), tenantID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Failed to retrieve locations",
			},
		})
	}

	var responses []dto.LocationResponse
	for _, location := range locations {
		responses = append(responses, lc.toLocationResponse(location))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    responses,
		"meta": fiber.Map{
			"limit":  limit,
			"offset": offset,
			"count":  len(responses),
		},
	})
}

// GetActiveLocationsByTenant retrieves active locations by tenant ID
func (lc *LocationController) GetActiveLocationsByTenant(c fiber.Ctx) error {
	tenantIDParam := c.Params("tenantId")
	tenantID, err := uuid.Parse(tenantIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid tenant ID",
			},
		})
	}

	locations, err := lc.locationUseCase.GetActiveLocationsByTenant(c.Context(), tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Failed to retrieve locations",
			},
		})
	}

	var responses []dto.LocationResponse
	for _, location := range locations {
		responses = append(responses, lc.toLocationResponse(location))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    responses,
	})
}

// UpdateLocation updates location information
func (lc *LocationController) UpdateLocation(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid location ID",
			},
		})
	}

	var req dto.UpdateLocationRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid request body",
			},
		})
	}

	location, err := lc.locationUseCase.UpdateLocation(c.Context(), id, req.Name, req.Address, req.City, req.State, req.Country, req.PostalCode, req.Phone, req.Email, req.Description, req.Capacity)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	response := lc.toLocationResponse(location)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// ActivateLocation activates a location
func (lc *LocationController) ActivateLocation(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid location ID",
			},
		})
	}

	location, err := lc.locationUseCase.ActivateLocation(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	response := lc.toLocationResponse(location)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeactivateLocation deactivates a location
func (lc *LocationController) DeactivateLocation(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid location ID",
			},
		})
	}

	location, err := lc.locationUseCase.DeactivateLocation(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	response := lc.toLocationResponse(location)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeleteLocation deletes a location
func (lc *LocationController) DeleteLocation(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid location ID",
			},
		})
	}

	err = lc.locationUseCase.DeleteLocation(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Location deleted successfully",
	})
}

// toLocationResponse converts entity to response DTO
func (lc *LocationController) toLocationResponse(location *entities.Location) dto.LocationResponse {
	return dto.LocationResponse{
		ID:          location.ID,
		TenantID:    location.TenantID,
		Name:        location.Name,
		Address:     location.Address,
		City:        location.City,
		State:       location.State,
		Country:     location.Country,
		PostalCode:  location.PostalCode,
		Phone:       location.Phone,
		Email:       location.Email,
		IsActive:    location.IsActive,
		Capacity:    location.Capacity,
		Description: location.Description,
		CreatedAt:   location.CreatedAt,
		UpdatedAt:   location.UpdatedAt,
	}
}