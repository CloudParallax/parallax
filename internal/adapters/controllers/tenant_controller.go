package controllers

import (
	"strconv"

	"github.com/cloudparallax/parallax/internal/adapters/http/dto"
	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/cloudparallax/parallax/internal/usecases"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// TenantController handles tenant HTTP requests
type TenantController struct {
	tenantUseCase *usecases.TenantUseCase
}

// NewTenantController creates a new tenant controller
func NewTenantController(tenantUseCase *usecases.TenantUseCase) *TenantController {
	return &TenantController{
		tenantUseCase: tenantUseCase,
	}
}

// CreateTenant creates a new tenant
func (tc *TenantController) CreateTenant(c fiber.Ctx) error {
	var req dto.CreateTenantRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid request body",
			},
		})
	}

	tenant, err := tc.tenantUseCase.CreateTenant(c.RequestCtx(), req.Name, req.Domain, req.Plan, req.MaxUsers, req.MaxLocations)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	response := tc.toTenantResponse(tenant)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetTenant retrieves a tenant by ID
func (tc *TenantController) GetTenant(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid tenant ID",
			},
		})
	}

	tenant, err := tc.tenantUseCase.GetTenant(c.RequestCtx(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusNotFound,
				"message": "Tenant not found",
			},
		})
	}

	response := tc.toTenantResponse(tenant)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetTenants retrieves all tenants with pagination
func (tc *TenantController) GetTenants(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	tenants, err := tc.tenantUseCase.GetAllTenants(c.RequestCtx(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Failed to retrieve tenants",
			},
		})
	}

	var responses []dto.TenantResponse
	for _, tenant := range tenants {
		responses = append(responses, tc.toTenantResponse(tenant))
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

// UpdateTenant updates tenant information
func (tc *TenantController) UpdateTenant(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid tenant ID",
			},
		})
	}

	var req dto.UpdateTenantRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid request body",
			},
		})
	}

	tenant, err := tc.tenantUseCase.UpdateTenant(c.RequestCtx(), id, req.Name, req.Plan, req.MaxUsers, req.MaxLocations)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	response := tc.toTenantResponse(tenant)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// ActivateTenant activates a tenant
func (tc *TenantController) ActivateTenant(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid tenant ID",
			},
		})
	}

	tenant, err := tc.tenantUseCase.ActivateTenant(c.RequestCtx(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	response := tc.toTenantResponse(tenant)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeactivateTenant deactivates a tenant
func (tc *TenantController) DeactivateTenant(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid tenant ID",
			},
		})
	}

	tenant, err := tc.tenantUseCase.DeactivateTenant(c.RequestCtx(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	response := tc.toTenantResponse(tenant)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeleteTenant deletes a tenant
func (tc *TenantController) DeleteTenant(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid tenant ID",
			},
		})
	}

	err = tc.tenantUseCase.DeleteTenant(c.RequestCtx(), id)
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
		"message": "Tenant deleted successfully",
	})
}

// toTenantResponse converts entity to response DTO
func (tc *TenantController) toTenantResponse(tenant *entities.Tenant) dto.TenantResponse {
	return dto.TenantResponse{
		ID:           tenant.ID,
		Name:         tenant.Name,
		Domain:       tenant.Domain,
		IsActive:     tenant.IsActive,
		Plan:         tenant.Plan,
		MaxUsers:     tenant.MaxUsers,
		MaxLocations: tenant.MaxLocations,
		CreatedAt:    tenant.CreatedAt,
		UpdatedAt:    tenant.UpdatedAt,
	}
}
