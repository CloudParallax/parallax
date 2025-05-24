package main

import (
	"github.com/cloudparallax/parallax/internal/config"
	"github.com/cloudparallax/parallax/web/app"
)

func main() {
	config.LoadEnvConfig()
	app.LoadApp()
}
