package app

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudparallax/parallax/internal/handlers"
	"github.com/cloudparallax/parallax/internal/middleware"
	staticfs "github.com/cloudparallax/parallax/web/static"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func LoadApp() {
	port := os.Getenv("SERVER_PORT")

	app := fiber.New()

	// Initialize default config
	app.Use(logger.New())

	if os.Getenv("ENV") == "production" {
		app.Use("/static*", static.New("", static.Config{
			FS:     staticfs.StaticFS,
			Browse: false,
		}))
	} else {
		app.Get("/static*", static.New("", static.Config{
			FS:     os.DirFS("web/static"),
			Browse: false,
		}))
	}

	middleware.LoadMiddleware(app)
	handlers.LoadHandlers(app)
	fmt.Printf("Staring Server at :%s", port)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
