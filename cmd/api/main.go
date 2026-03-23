package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/wil/hrms/pkg/db"
	"github.com/wil/hrms/pkg/middleware"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	connStr := os.Getenv("DATABASE_URL")
	db.InitDB(connStr)

	r := gin.Default()

	// Global middleware
	r.Use(middleware.ErrorHandler())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	// db connection
	r.GET("/api/v1/health", func(c *gin.Context) {
	err := db.GetPool().Ping(context.Background())

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
