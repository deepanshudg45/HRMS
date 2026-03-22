package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/wil/hrms/pkg/db"
	"github.com/wil/hrms/pkg/middleware"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := gin.Default()

	// Global middleware
	r.Use(middleware.ErrorHandler())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	// db connection
	r.GET("/api/v1/health", func(c *gin.Context) {
	err := db.GetPool().Ping(c)

	if err != nil {
		c.JSON(500, gin.H{
			"status": "DOWN",
			"db":     "not reachable",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "UP",
		"db":     "connected",
	})
})
	logger.Info("Server started on :8080")
	r.Run(":8080")
}
