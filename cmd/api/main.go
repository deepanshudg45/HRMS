package main

import (
	// "fmt"
	"os"
	"github.com/gofiber/fiber/v2"
	"hrms-1/pkg/db"
	// "hrms-1/pkg/middleware"
	"hrms-1/internal/routes"
	"go.uber.org/zap"
	"hrms-1/pkg/config"
)

func main() {
	config.LoadEnv()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	connStr := os.Getenv("DATABASE_URL")
	db.InitDB(connStr)

	err := db.RunMigrations(db.GetPool())
	if err != nil {
		panic(err)
	}
	// fmt.Println("DB URL:", connStr)


	app :=fiber.New()
	routes.AssetsRoutes(app)

	// Global middleware
	// app.Use(middleware.ErrorHandler())

	app.Get("/health", func(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status": "ok",
	})
})
	
	logger.Info("Server started on :8000")
	app.Listen(":8000")
}
