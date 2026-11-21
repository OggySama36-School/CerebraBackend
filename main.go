package main

import (
	"doanbackend/config"
	"doanbackend/router"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	config.LoadConfig()
	r := router.SetUpRoutes()
	port := config.App.Port
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
