package main

import (
	"doanbackend/config"
	"doanbackend/router"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	config.LoadConfig()
	r := router.SetUpRoutes()
	port := fmt.Sprintf(":%s", config.App.Port)
	r.Run(port)
}
