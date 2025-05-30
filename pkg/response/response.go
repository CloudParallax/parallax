package response

import (
	"log"
	"net/http"

	"github.com/cloudparallax/parallax/pkg/errors"
	"github.com/gofiber/fiber/v3"
)

// Response represents a standard API response
type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
	Meta    *Meta      `json:"meta,omitempty"`
}

// ErrorInfo represents error information in the response
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// Meta represents metadata for the response
type Meta struct {
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// Success sends a successful response
func Success(c fiber.Ctx, data any) error {
	return c.Status(http.StatusOK).JSON(Response{
		Success: true,
		Data:    data,
	})
}

// SuccessWithStatus sends a successful response with custom status
func SuccessWithStatus(c fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(Response{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMeta sends a successful response with metadata
func SuccessWithMeta(c fiber.Ctx, data any, meta *Meta) error {
	return c.Status(http.StatusOK).JSON(Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Created sends a 201 Created response
func Created(c fiber.Ctx, data any) error {
	return c.Status(http.StatusCreated).JSON(Response{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a 204 No Content response
func NoContent(c fiber.Ctx) error {
	return c.SendStatus(http.StatusNoContent)
}

// Error sends an error response
func Error(c fiber.Ctx, err error) error {
	// Log the error
	log.Printf("API Error: %v", err)

	// Check if it's an AppError
	if appErr, ok := errors.GetAppError(err); ok {
		return c.Status(appErr.Status).JSON(Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Details,
			},
		})
	}

	// Check if it's a ValidationErrors
	if valErr, ok := err.(*errors.ValidationErrors); ok {
		return c.Status(http.StatusBadRequest).JSON(Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    "VALIDATION_FAILED",
				Message: "Validation failed",
				Details: valErr.Errors,
			},
		})
	}

	// Generic error
	return c.Status(http.StatusInternalServerError).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "An unexpected error occurred",
		},
	})
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c fiber.Ctx, message string) error {
	return c.Status(http.StatusBadRequest).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "BAD_REQUEST",
			Message: message,
		},
	})
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c fiber.Ctx, message string) error {
	return c.Status(http.StatusUnauthorized).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "UNAUTHORIZED",
			Message: message,
		},
	})
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c fiber.Ctx, message string) error {
	return c.Status(http.StatusForbidden).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "FORBIDDEN",
			Message: message,
		},
	})
}

// NotFound sends a 404 Not Found response
func NotFound(c fiber.Ctx, message string) error {
	return c.Status(http.StatusNotFound).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "NOT_FOUND",
			Message: message,
		},
	})
}

// Conflict sends a 409 Conflict response
func Conflict(c fiber.Ctx, message string) error {
	return c.Status(http.StatusConflict).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "CONFLICT",
			Message: message,
		},
	})
}

// UnprocessableEntity sends a 422 Unprocessable Entity response
func UnprocessableEntity(c fiber.Ctx, message string) error {
	return c.Status(http.StatusUnprocessableEntity).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "UNPROCESSABLE_ENTITY",
			Message: message,
		},
	})
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c fiber.Ctx, message string) error {
	return c.Status(http.StatusInternalServerError).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: message,
		},
	})
}

// ParseJSON parses JSON request body
func ParseJSON(c fiber.Ctx, v any) error {
	if err := c.Bind().JSON(v); err != nil {
		return errors.NewAppErrorWithDetails(
			"INVALID_JSON",
			"Invalid JSON format",
			err.Error(),
			http.StatusBadRequest,
		)
	}
	return nil
}

// NewMeta creates pagination metadata
func NewMeta(page, limit, total int) *Meta {
	totalPages := (total + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	return &Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
