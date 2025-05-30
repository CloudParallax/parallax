package http

import (
	"github.com/cloudparallax/parallax/internal/adapters/controllers"
	"github.com/cloudparallax/parallax/internal/adapters/http/middleware"
	"github.com/cloudparallax/parallax/internal/adapters/repositories"
	"github.com/cloudparallax/parallax/internal/usecases"
	"github.com/gofiber/fiber/v3"
)

// Router handles API route setup
type Router struct {
	app        *fiber.App
	middleware *middleware.MiddlewareManager
}

// NewRouter creates a new router instance
func NewRouter(app *fiber.App) *Router {
	return &Router{
		app:        app,
		middleware: middleware.NewMiddlewareManager(),
	}
}

// SetupRoutes configures all API routes
func (r *Router) SetupRoutes() {
	// Setup global middleware
	r.middleware.SetupGlobalMiddleware(r.app)
	
	// Start cleanup routine for middlewares
	r.middleware.StartCleanupRoutine()

	// Initialize repositories
	tenantRepo := repositories.NewMemoryTenantRepository()
	locationRepo := repositories.NewMemoryLocationRepository()
	customerRepo := repositories.NewMemoryCustomerRepository()

	// Initialize use cases
	tenantUseCase := usecases.NewTenantUseCase(tenantRepo)
	locationUseCase := usecases.NewLocationUseCase(locationRepo, tenantRepo)
	customerUseCase := usecases.NewCustomerUseCase(customerRepo, tenantRepo)

	// Initialize controllers
	tenantController := controllers.NewTenantController(tenantUseCase)
	locationController := controllers.NewLocationController(locationUseCase)
	customerController := controllers.NewCustomerController(customerUseCase)

	// Setup API routes
	api := r.app.Group("/api/v1")

	// Health check (no auth required)
	api.Get("/health", r.healthCheck)

	// Hello World endpoint (no auth required)
	api.Get("/hello", r.helloWorld)

	// Auth routes (no auth required)
	r.setupAuthRoutes(api)

	// Public routes (optional auth)
	r.setupPublicRoutes(api, tenantController, locationController, customerController)

	// Protected routes (auth required)
	r.setupProtectedRoutes(api, tenantController, locationController, customerController)
}

// setupAuthRoutes configures authentication routes
func (r *Router) setupAuthRoutes(api fiber.Router) {
	auth := api.Group("/auth")

	// Login endpoint
	auth.Post("/login", r.login)
	
	// Logout endpoint
	auth.Post("/logout", r.logout)
	
	// Get CSRF token
	auth.Get("/csrf", r.getCSRFToken)
	
	// Check auth status
	auth.Get("/me", r.middleware.OptionalAuth(), r.getAuthStatus)
}

// setupPublicRoutes configures public routes (optional auth)
func (r *Router) setupPublicRoutes(api fiber.Router, tenantController *controllers.TenantController, locationController *controllers.LocationController, customerController *controllers.CustomerController) {
	// Apply optional auth to all public routes
	public := api.Group("/", r.middleware.OptionalAuth())
	
	// Public tenant routes (read-only)
	tenants := public.Group("/tenants")
	tenants.Get("/", tenantController.GetTenants)
	tenants.Get("/:id", tenantController.GetTenant)
	
	// Public location routes (read-only)
	locations := public.Group("/tenants/:tenantId/locations")
	locations.Get("/", locationController.GetLocationsByTenant)
	locations.Get("/active", locationController.GetActiveLocationsByTenant)
	
	// Individual location routes
	location := public.Group("/locations")
	location.Get("/:id", locationController.GetLocation)
	
	// Public customer routes (read-only)
	customers := public.Group("/tenants/:tenantId/customers")
	customers.Get("/", customerController.GetCustomersByTenant)
	customers.Get("/search", customerController.SearchCustomers)
	
	// Individual customer routes
	customer := public.Group("/customers")
	customer.Get("/:id", customerController.GetCustomer)
}

// setupProtectedRoutes configures protected routes (auth required)
func (r *Router) setupProtectedRoutes(api fiber.Router, tenantController *controllers.TenantController, locationController *controllers.LocationController, customerController *controllers.CustomerController) {
	// Apply auth and CSRF protection to all protected routes
	protected := api.Group("/", r.middleware.RequireAuth(), r.middleware.SetupCSRF())
	
	// Protected tenant routes (write operations)
	tenants := protected.Group("/tenants")
	tenants.Post("/", tenantController.CreateTenant)
	tenants.Put("/:id", tenantController.UpdateTenant)
	tenants.Delete("/:id", tenantController.DeleteTenant)
	tenants.Post("/:id/activate", tenantController.ActivateTenant)
	tenants.Post("/:id/deactivate", tenantController.DeactivateTenant)
	
	// Protected location routes (write operations)
	locations := protected.Group("/tenants/:tenantId/locations")
	locations.Post("/", locationController.CreateLocation)
	
	// Individual location write operations
	location := protected.Group("/locations")
	location.Put("/:id", locationController.UpdateLocation)
	location.Delete("/:id", locationController.DeleteLocation)
	location.Post("/:id/activate", locationController.ActivateLocation)
	location.Post("/:id/deactivate", locationController.DeactivateLocation)
	
	// Protected customer routes (write operations)
	customers := protected.Group("/tenants/:tenantId/customers")
	customers.Post("/", customerController.CreateCustomer)
	
	// Individual customer write operations
	customer := protected.Group("/customers")
	customer.Put("/:id", customerController.UpdateCustomer)
	customer.Delete("/:id", customerController.DeleteCustomer)
	customer.Post("/:id/activate", customerController.ActivateCustomer)
	customer.Post("/:id/deactivate", customerController.DeactivateCustomer)
	customer.Post("/:id/tags", customerController.AddTag)
	customer.Delete("/:id/tags", customerController.RemoveTag)
	
	// Admin routes (require admin role)
	admin := protected.Group("/admin", r.middleware.RequireRole("admin"))
	admin.Get("/users", r.getUsers)
	admin.Delete("/users/:id", r.deleteUser)
}

// login handles user authentication
func (r *Router) login(c fiber.Ctx) error {
	// TODO: Implement actual authentication logic
	// For now, this is a placeholder that accepts any credentials
	
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid request body",
			},
		})
	}
	
	// Placeholder authentication (replace with real auth)
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "Invalid credentials",
			},
		})
	}
	
	// Create session
	userData := map[string]interface{}{
		"username": req.Username,
		"role":     "user", // Default role
	}
	
	if err := r.middleware.Login(c, req.Username, userData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Failed to create session",
			},
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": userData,
		},
	})
}

// helloWorld provides a simple hello world endpoint to test the API
func (r *Router) helloWorld(c fiber.Ctx) error {
	name := c.Query("name", "World")
	
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Hello, " + name + "!",
		"api":     "Parallax Workplace Management API",
		"version": "1.0.0",
		"middleware": fiber.Map{
			"cors":        "enabled",
			"csrf":       "ready",
			"auth":       "available",
			"rate_limit": "active",
		},
	})
}

// logout handles user logout
func (r *Router) logout(c fiber.Ctx) error {
	if err := r.middleware.Logout(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Failed to logout",
			},
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}

// getCSRFToken returns the CSRF token for the client
func (r *Router) getCSRFToken(c fiber.Ctx) error {
	token := r.middleware.GetCSRFToken(c)
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"csrf_token": token,
		},
	})
}

// getAuthStatus returns the current authentication status
func (r *Router) getAuthStatus(c fiber.Ctx) error {
	authenticated := c.Locals("authenticated")
	if authenticated == true {
		session := c.Locals("session")
		userID := c.Locals("user_id")
		
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"authenticated": true,
				"user_id":       userID,
				"session":       session,
			},
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"authenticated": false,
		},
	})
}

// getUsers returns list of users (admin only)
func (r *Router) getUsers(c fiber.Ctx) error {
	// TODO: Implement user listing
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"users": []interface{}{},
		},
	})
}

// deleteUser deletes a user (admin only)
func (r *Router) deleteUser(c fiber.Ctx) error {
	userID := c.Params("id")
	
	// TODO: Implement user deletion
	return c.JSON(fiber.Map{
		"success": true,
		"message": "User " + userID + " deleted successfully",
	})
}

// healthCheck provides a simple health check endpoint
func (r *Router) healthCheck(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "parallax-workplace-management-api",
		"version": "1.0.0",
	})
}