package router

import (
	"doanbackend/config"
	"doanbackend/handler"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes() *gin.Engine {
	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.App.AllowedOrigins}, //strings.Split(config.App.AllowedOrigins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	public := engine.Group("/api/v1")
	{
		public.POST("/predict", handler.HandleData)
	}

	return engine
}
