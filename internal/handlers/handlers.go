package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/cloudparallax/parallax/web/templates"
)

func LoadHandlers(router *gin.Engine) {
	router.GET("/", func(context *gin.Context) {
		component := templates.Splash()
		templates.Layout(component).Render(context, context.Writer)
	})

	loadCounterHandler(router)
}
