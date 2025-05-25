package views

import (
	"github.com/a-h/templ"
	"github.com/cloudparallax/parallax/web/templates"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/csrf"
)

func Render(c fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	csrfToken := csrf.TokenFromContext(c)
	componentHandler := templ.Handler(templates.Layout(component, csrfToken))
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}
