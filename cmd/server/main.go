package main

import (
	"context"
	"log"

	"WITS/config"
	"WITS/internal/assets/handler"
	"WITS/internal/assets/repository"
	"WITS/internal/assets/routes"
	"WITS/internal/assets/service"
	authHandlerPkg "WITS/internal/auth/handler"
	authRoutes "WITS/internal/auth/routes"
	"WITS/package/database"
	"WITS/package/migration"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	godotenv.Load()

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := migration.Run(context.Background(), db, "migration"); err != nil {
		log.Fatal(err)
	}

	log.Println("migrations completed successfully")

	app := fiber.New()
	app.Use(cors.New())

	authHandler := authHandlerPkg.NewAuthHandler(cfg.JWTSecret)
	assetRepo := repository.NewRepository(db)
	assetService := service.NewAssetService(assetRepo)
	assetHandler := handler.NewAssetHandler(assetService)

	authRoutes.AuthRoutes(app, authHandler)
	routes.AssetRoutes(app, assetHandler)

	warrantyCron := cron.New()
	_, err = warrantyCron.AddFunc("0 7 * * *", func() {
		if err := assetService.CheckWarrantyExpiry(context.Background()); err != nil {
			log.Println("warranty cron error:", err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	warrantyCron.Start()
	defer warrantyCron.Stop()

	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
