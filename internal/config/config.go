package config

import (
	"github.com/joho/godotenv"
)

func LoadEnvConfig() {
	godotenv.Load(".env")

	// if err != nil {
	// log.Fatalf("The config module cannot load the .env file")
	// }
}
