package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cloudparallax/parallax/internal/handlers"
	"github.com/cloudparallax/parallax/internal/middleware"
	"github.com/cloudparallax/parallax/web/static"
	"github.com/gin-gonic/gin"
)

func LoadApp() {
	port := os.Getenv("SERVER_PORT")
	router := gin.Default()
	if os.Getenv("ENV") == "production" {
		router.StaticFS("/static", http.FS(static.StaticFS))
		gin.SetMode(gin.ReleaseMode)
	} else {
		router.Static("/static", "./web/static")
		gin.SetMode(gin.DebugMode)
	}

	middleware.LoadMiddleware(router)
	handlers.LoadHandlers(router)
	fmt.Printf("Staring Server at :%s", port)

	router.Run(fmt.Sprintf(":%s", port))
}
