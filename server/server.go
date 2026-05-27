package server

import (
	"net/http"
	"time"

	"agentcore/config"
	"agentcore/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Start() error {
	// Set the mode
	if config.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Setup routes
	setupRoutes(r)

	// Start the server
	port := config.ServerPort
	logrus.Infof("Starting server on port %s", port)
	return r.Run(":" + port)
}

func setupRoutes(r *gin.Engine) {
	// Example routes
	r.GET("/api/v1/users", handlers.GetUsers)
	r.GET("/api/v1/users/:id", handlers.GetUser)
	r.POST("/api/v1/users", handlers.CreateUser)
	r.PUT("/api/v1/users/:id", handlers.UpdateUser)
	r.DELETE("/api/v1/users/:id", handlers.DeleteUser)
	r.POST("/api/v1/query", handlers.Query)

	// Add more routes here as needed
}
