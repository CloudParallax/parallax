package handlers

import (
	"encoding/json" // Added for JSON marshaling
	"fmt"           // Added for string formatting
	"math/rand"     // Added for chart data generation
	"os"
	"os/exec"
	"time" // Added for chart data generation

	"github.com/cloudparallax/parallax/internal/services"
	"github.com/cloudparallax/parallax/web/templates/components"
	"github.com/cloudparallax/parallax/web/templates/layouts"
	"github.com/cloudparallax/parallax/web/templates/pages"
	"github.com/cloudparallax/parallax/web/templates/responses"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
)

// ChartData structures for dynamic charting
type PerformanceDataPoint struct {
	Date    string `json:"date"`
	CPU     int    `json:"cpu"`
	Memory  int    `json:"memory"`
	Network int    `json:"network"`
}

type ResourceDataPoint struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	Color string `json:"color"`
}

// AllChartData bundles all data for the charts on the homepage
type AllChartData struct {
	PerformanceData []PerformanceDataPoint `json:"performanceData"`
	ResourceData    []ResourceDataPoint    `json:"resourceData"`
	TrafficData     []ResourceDataPoint    `json:"trafficData"` // Reusing ResourceDataPoint for Traffic Distribution
}

var blogService *services.BlogService

// chartDataHandler provides data for the homepage charts
func chartDataHandler(c fiber.Ctx) error {
	// Generate sample data
	// In a real application, this data would come from a database or other services.
	performanceData := make([]PerformanceDataPoint, 7)
	for i := 0; i < 7; i++ {
		performanceData[i] = PerformanceDataPoint{
			Date:    time.Now().AddDate(0, 0, -(6 - i)).Format("2006-01-02"),
			CPU:     rand.Intn(50) + 25, // Random values between 25 and 74
			Memory:  rand.Intn(50) + 25,
			Network: rand.Intn(50) + 25,
		}
	}

	resourceData := []ResourceDataPoint{
		{Name: "CPU", Value: rand.Intn(70) + 30, Color: "#3B82F6"},     // Random values between 30 and 99
		{Name: "Memory", Value: rand.Intn(70) + 30, Color: "#10B981"},
		{Name: "Storage", Value: rand.Intn(60) + 20, Color: "#8B5CF6"},  // Random values between 20 and 79
		{Name: "Network", Value: rand.Intn(70) + 30, Color: "#F59E0B"},
	}

	trafficData := []ResourceDataPoint{ // Reusing ResourceDataPoint for Traffic Distribution
		{Name: "Web Traffic", Value: rand.Intn(40) + 20, Color: "#3B82F6"}, // Random values between 20 and 59
		{Name: "API Calls", Value: rand.Intn(30) + 10, Color: "#10B981"},   // Random values between 10 and 39
		{Name: "Database", Value: rand.Intn(15) + 5, Color: "#8B5CF6"},    // Random values between 5 and 19
		{Name: "Cache", Value: rand.Intn(10) + 5, Color: "#F59E0B"},       // Random values between 5 and 14
	}

	allData := AllChartData{
		PerformanceData: performanceData,
		ResourceData:    resourceData,
		TrafficData:     trafficData,
	}

	jsonData, err := json.Marshal(allData)
	if err != nil {
		// In a real application, you should log this error.
		// log.Printf("Error marshaling chart data: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error preparing chart data script.")
	}

	// Construct the script tag that HTMX will swap in.
	// This script calls the JavaScript function `updateAllCharts` with the fetched data.
	scriptContent := fmt.Sprintf(`
<script type="text/javascript">
  if (typeof window.updateAllCharts === 'function') {
    window.updateAllCharts(%s);
  } else {
    // If updateAllCharts is not yet defined, wait for DOMContentLoaded.
    // This can happen if the script defining updateAllCharts is loaded after this HTMX response is processed.
    document.addEventListener('DOMContentLoaded', function() {
      if (typeof window.updateAllCharts === 'function') {
        window.updateAllCharts(%s);
      } else {
        console.error('updateAllCharts function not found even after DOMContentLoaded.');
      }
    });
  }
</script>`, string(jsonData), string(jsonData))

	// Set Content-Type to text/html as we are sending a script tag
	return c.Type("html").SendString(scriptContent)
}

// LoadHandlers sets up all application routes and handlers
func LoadHandlers(app *fiber.App) {
	// Initialize services
	blogService = services.NewBlogService("content/blog")

	// Register route groups
	setupMainRoutes(app)
	setupAPIRoutes(app)
	setupDemoRoutes(app)
	
	// Load additional handlers
	loadCounterHandler(app)
}

// setupMainRoutes configures the main application routes
func setupMainRoutes(app *fiber.App) {
	// Home page
	app.Get("/", func(c fiber.Ctx) error {
		csrfToken := csrf.TokenFromContext(c)
		component := pages.HomePage()

		if isHTMXRequest(c) {
			return component.Render(c.Context(), c.Response().BodyWriter())
		}

		return layouts.BaseLayout("Home", csrfToken, component).Render(c.Context(), c.Response().BodyWriter())
	})

	// Blog routes
	app.Get("/blog", func(c fiber.Ctx) error {
		csrfToken := csrf.TokenFromContext(c)
		posts, err := blogService.GetAllPosts()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error loading blog posts")
		}

		component := pages.BlogListPage(posts)

		if isHTMXRequest(c) {
			return component.Render(c.Context(), c.Response().BodyWriter())
		}

		return layouts.BaseLayout("Blog", csrfToken, component).Render(c.Context(), c.Response().BodyWriter())
	})

	app.Get("/blog/:slug", func(c fiber.Ctx) error {
		csrfToken := csrf.TokenFromContext(c)
		slug := c.Params("slug")

		post, err := blogService.GetPostBySlug(slug)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Blog post not found")
		}

		component := pages.BlogPostPage(*post)

		if isHTMXRequest(c) {
			return component.Render(c.Context(), c.Response().BodyWriter())
		}

		return layouts.BaseLayout(post.Title, csrfToken, component).Render(c.Context(), c.Response().BodyWriter())
	})

	// About page
	app.Get("/about", func(c fiber.Ctx) error {
		csrfToken := csrf.TokenFromContext(c)
		component := pages.AboutPage()

		if isHTMXRequest(c) {
			return component.Render(c.Context(), c.Response().BodyWriter())
		}

		return layouts.BaseLayout("About", csrfToken, component).Render(c.Context(), c.Response().BodyWriter())
	})

	// Contact page
	app.Get("/contact", func(c fiber.Ctx) error {
		csrfToken := csrf.TokenFromContext(c)
		component := pages.ContactPage()

		if isHTMXRequest(c) {
			return component.Render(c.Context(), c.Response().BodyWriter())
		}

		return layouts.BaseLayout("Contact", csrfToken, component).Render(c.Context(), c.Response().BodyWriter())
	})
}

// setupAPIRoutes configures API endpoints
func setupAPIRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Form handler for CSRF demonstration
	api.Post("/test-form", func(c fiber.Ctx) error {
		message := c.FormValue("message")
		if message == "" {
			message = "No message provided"
		}

		component := responses.FormSuccessResponse(message)
		return component.Render(c.Context(), c.Response().BodyWriter())
	})

	// CSS rebuild endpoint for development
	api.Post("/rebuild-css", func(c fiber.Ctx) error {
		if os.Getenv("ENV") == "production" {
			return c.Status(fiber.StatusForbidden).SendString("Not available in production")
		}

		cmd := exec.Command("tailwindcss", "-i", "./web/static/main.css", "-o", "./web/static/dist/output.css")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			errorMsg := `<div class="text-red-700 bg-red-100 border border-red-200 rounded p-2">
				<strong>❌ CSS rebuild failed!</strong><br/>
				Error: ` + err.Error() + `<br/>
				<pre class="text-xs mt-2 text-red-600">` + string(output) + `</pre>
			</div>`
			return c.Status(fiber.StatusInternalServerError).Type("text/html").SendString(errorMsg)
		}

		component := responses.BuildSuccessResponse(string(output))
		return component.Render(c.Context(), c.Response().BodyWriter())
	})

	// Chart data endpoint
	api.Get("/chart-data", chartDataHandler)
}

// setupDemoRoutes configures demo routes for HTMX examples
func setupDemoRoutes(app *fiber.App) {
	demo := app.Group("/api/demo")

	// Dynamic content for HTMX demonstrations
	demo.Get("/content/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		component := components.DemoContent(id)
		
		return component.Render(c.Context(), c.Response().BodyWriter())
	})
}

// isHTMXRequest checks if the request comes from HTMX
func isHTMXRequest(c fiber.Ctx) bool {
	return c.Get("HX-Request") == "true"
}