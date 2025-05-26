package app

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudparallax/parallax/internal/handlers"
	"github.com/cloudparallax/parallax/internal/middleware"
	staticfs "github.com/cloudparallax/parallax/web/static"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func LoadApp() {
	port := os.Getenv("SERVER_PORT")

	app := fiber.New()

	// Initialize default config
	app.Use(logger.New())

	// Initialize default config
	app.Use(compress.New())

	// Initialize default config
	app.Use(recoverer.New())

	if os.Getenv("ENV") == "production" {
		app.Use("/static*", static.New("", static.Config{
			FS:     staticfs.StaticFS,
			Browse: false,
		}))

		app.Use(favicon.New(favicon.Config{
			File:       "favicon.ico",
			FileSystem: staticfs.StaticFS,
		}))
	} else {
		app.Get("/static*", static.New("", static.Config{
			FS:     os.DirFS("web/static"),
			Browse: false,
		}))
		app.Use(favicon.New(favicon.Config{
			File: "./web/static/favicon.ico",
			URL:  "/favicon.ico",
		}))
	}

	middleware.LoadMiddleware(app)
	handlers.LoadHandlers(app)
	fmt.Printf("Staring Server at :%s", port)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}

func GetEnv(s1, s2 string) string {
	if value, exists := os.LookupEnv(s1); exists {
		return value
	}
	return s2
}
