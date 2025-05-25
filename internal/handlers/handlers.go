package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/cloudparallax/parallax/web/templates"
)

func LoadHandlers(app *fiber.App) {
	app.Get("/", func(c fiber.Ctx) error {
		component := templates.Splash()
		return templates.Layout(component).Render(c.Context(), c.Response().BodyWriter())
	})

	loadCounterHandler(app)
}