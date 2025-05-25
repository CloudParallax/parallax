package handlers

import (
	"github.com/cloudparallax/parallax/internal/views"
	"github.com/cloudparallax/parallax/web/templates"
	"github.com/gofiber/fiber/v3"
)

func LoadHandlers(app *fiber.App) {
	app.Get("/", func(c fiber.Ctx) error {
		component := templates.Splash()
		return templates.Layout(component).Render(c.Context(), c.Response().BodyWriter())
	})

	app.Get("/test", func(c fiber.Ctx) error {
		return views.Render(c, templates.Splash())
	})

	loadCounterHandler(app)
}
