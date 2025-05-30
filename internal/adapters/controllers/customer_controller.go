package controllers

import (
	"strconv"

	"github.com/cloudparallax/parallax/internal/adapters/http/dto"
	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/cloudparallax/parallax/internal/usecases"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// CustomerController handles customer HTTP requests
type CustomerController struct {
	customerUseCase *usecases.CustomerUseCase
}

// NewCustomerController creates a new customer controller
func NewCustomerController(customerUseCase *usecases.CustomerUseCase) *CustomerController {
	return &CustomerController{
		customerUseCase: customerUseCase,
	}
}

// CreateCustomer creates a new customer
func (cc *CustomerController) CreateCustomer(c fiber.Ctx) error {
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

	var req dto.CreateCustomerRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid request body",
			},
		})
	}

	customer, err := cc.customerUseCase.CreateCustomer(c.Context(), tenantID, req.FirstName, req.LastName, req.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	response := cc.toCustomerResponse(customer)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetCustomer retrieves a customer by ID
func (cc *CustomerController) GetCustomer(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid customer ID",
			},
		})
	}

	customer, err := cc.customerUseCase.GetCustomer(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusNotFound,
				"message": "Customer not found",
			},
		})
	}

	response := cc.toCustomerResponse(customer)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetCustomersByTenant retrieves customers by tenant ID with pagination
func (cc *CustomerController) GetCustomersByTenant(c fiber.Ctx) error {
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

	customers, err := cc.customerUseCase.GetCustomersByTenant(c.Context(), tenantID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Failed to retrieve customers",
			},
		})
	}

	var responses []dto.CustomerResponse
	for _, customer := range customers {
		responses = append(responses, cc.toCustomerResponse(customer))
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

// SearchCustomers searches customers by name
func (cc *CustomerController) SearchCustomers(c fiber.Ctx) error {
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

	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Search query is required",
			},
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	customers, err := cc.customerUseCase.SearchCustomersByName(c.Context(), tenantID, query, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Failed to search customers",
			},
		})
	}

	var responses []dto.CustomerResponse
	for _, customer := range customers {
		responses = append(responses, cc.toCustomerResponse(customer))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    responses,
		"meta": fiber.Map{
			"limit":  limit,
			"offset": offset,
			"count":  len(responses),
			"query":  query,
		},
	})
}

// UpdateCustomer updates customer information
func (cc *CustomerController) UpdateCustomer(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid customer ID",
			},
		})
	}

	var req dto.UpdateCustomerRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid request body",
			},
		})
	}

	customer, err := cc.customerUseCase.UpdateCustomer(c.Context(), id, req.FirstName, req.LastName, req.Email, req.Phone, req.Address, req.City, req.State, req.Country, req.PostalCode, req.CompanyName, req.JobTitle, req.Notes, req.Tags)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	response := cc.toCustomerResponse(customer)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// ActivateCustomer activates a customer
func (cc *CustomerController) ActivateCustomer(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid customer ID",
			},
		})
	}

	customer, err := cc.customerUseCase.ActivateCustomer(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	response := cc.toCustomerResponse(customer)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeactivateCustomer deactivates a customer
func (cc *CustomerController) DeactivateCustomer(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid customer ID",
			},
		})
	}

	customer, err := cc.customerUseCase.DeactivateCustomer(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	response := cc.toCustomerResponse(customer)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// AddTag adds a tag to a customer
func (cc *CustomerController) AddTag(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid customer ID",
			},
		})
	}

	var req dto.AddTagRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid request body",
			},
		})
	}

	customer, err := cc.customerUseCase.AddCustomerTag(c.Context(), id, req.Tag)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	response := cc.toCustomerResponse(customer)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// RemoveTag removes a tag from a customer
func (cc *CustomerController) RemoveTag(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid customer ID",
			},
		})
	}

	var req dto.RemoveTagRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid request body",
			},
		})
	}

	customer, err := cc.customerUseCase.RemoveCustomerTag(c.Context(), id, req.Tag)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	response := cc.toCustomerResponse(customer)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeleteCustomer deletes a customer
func (cc *CustomerController) DeleteCustomer(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid customer ID",
			},
		})
	}

	err = cc.customerUseCase.DeleteCustomer(c.Context(), id)
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
		"message": "Customer deleted successfully",
	})
}

// toCustomerResponse converts entity to response DTO
func (cc *CustomerController) toCustomerResponse(customer *entities.Customer) dto.CustomerResponse {
	return dto.CustomerResponse{
		ID:          customer.ID,
		TenantID:    customer.TenantID,
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		Email:       customer.Email,
		Phone:       customer.Phone,
		Address:     customer.Address,
		City:        customer.City,
		State:       customer.State,
		Country:     customer.Country,
		PostalCode:  customer.PostalCode,
		CompanyName: customer.CompanyName,
		JobTitle:    customer.JobTitle,
		Notes:       customer.Notes,
		IsActive:    customer.IsActive,
		Tags:        customer.Tags,
		CreatedAt:   customer.CreatedAt,
		UpdatedAt:   customer.UpdatedAt,
	}
}