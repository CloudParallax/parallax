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
	blogRepo := repositories.NewMemoryBlogRepository()
	counterRepo := repositories.NewMemoryCounterRepository()

	// Initialize use cases
	blogUseCase := usecases.NewBlogUseCase(blogRepo)
	counterUseCase := usecases.NewCounterUseCase(counterRepo)

	// Initialize controllers
	blogController := controllers.NewBlogController(blogUseCase)
	counterController := controllers.NewCounterController(counterUseCase)

	// Setup API routes
	api := r.app.Group("/api/v1")

	// Health check (no auth required)
	api.Get("/health", r.healthCheck)

	// Hello World endpoint (no auth required)
	api.Get("/hello", r.helloWorld)

	// Auth routes (no auth required)
	r.setupAuthRoutes(api)

	// Public routes (optional auth)
	r.setupPublicRoutes(api, blogController, counterController)

	// Protected routes (auth required)
	r.setupProtectedRoutes(api, blogController, counterController)
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
func (r *Router) setupPublicRoutes(api fiber.Router, blogController *controllers.BlogController, counterController *controllers.CounterController) {
	// Apply optional auth to all public routes
	public := api.Group("/", r.middleware.OptionalAuth())
	
	blog := public.Group("/blog")
	
	// Public blog routes (read-only)
	blog.Get("/posts", blogController.GetPosts)
	blog.Get("/posts/published", blogController.GetPublishedPosts)
	blog.Get("/posts/search", blogController.SearchPosts)
	blog.Get("/posts/tags", blogController.GetPostsByTags)
	blog.Get("/posts/:id", blogController.GetPost)
	blog.Get("/posts/slug/:slug", blogController.GetPostBySlug)
	
	// Public counter routes (read-only)
	counter := public.Group("/counter")
	counter.Get("/", counterController.GetAllCounters)
	counter.Get("/:id", counterController.GetCounter)
}

// setupProtectedRoutes configures protected routes (auth required)
func (r *Router) setupProtectedRoutes(api fiber.Router, blogController *controllers.BlogController, counterController *controllers.CounterController) {
	// Apply auth and CSRF protection to all protected routes
	protected := api.Group("/", r.middleware.RequireAuth(), r.middleware.SetupCSRF())
	
	blog := protected.Group("/blog")
	
	// Protected blog routes (write operations)
	blog.Post("/posts", blogController.CreatePost)
	blog.Put("/posts/:id", blogController.UpdatePost)
	blog.Delete("/posts/:id", blogController.DeletePost)
	blog.Post("/posts/:id/publish", blogController.PublishPost)
	blog.Post("/posts/:id/unpublish", blogController.UnpublishPost)
	
	// Protected counter routes (write operations)
	counter := protected.Group("/counter")
	counter.Post("/", counterController.CreateCounter)
	counter.Delete("/:id", counterController.DeleteCounter)
	counter.Put("/:id/increment", counterController.IncrementCounter)
	counter.Put("/:id/decrement", counterController.DecrementCounter)
	counter.Put("/:id/value", counterController.SetCounterValue)
	counter.Put("/:id/reset", counterController.ResetCounter)
	
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
		"api":     "Parallax API",
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
		"service": "parallax-api",
		"version": "1.0.0",
	})
}