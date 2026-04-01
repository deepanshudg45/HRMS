package main

import (
	"log"

	"WITS/config"
	"WITS/internal/assets/handler"
	"WITS/internal/assets/repository"
	"WITS/internal/assets/routes"
	"WITS/internal/assets/service"
	"WITS/package/database"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := fiber.New()

	assetRepo := repository.NewRepository(db)
	assetService := service.NewAssetService(assetRepo)
	assetHandler := handler.NewAssetHandler(assetService)

	routes.AssetRoutes(app, assetHandler)

	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
