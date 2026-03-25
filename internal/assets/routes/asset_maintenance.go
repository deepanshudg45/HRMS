package routes

import (
    "WITS/internal/assets/handler"

    "github.com/gofiber/fiber/v2"
)

func AssetRoutes(app *fiber.App, handler *handler.AssetHandler) {
    api := app.Group("/api/assets")

    assets := api.Group("/assets")
    assets.Post("", handler.CreateAsset)
}
