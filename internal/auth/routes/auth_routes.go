package routes

import (
	"WITS/internal/auth/handler"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, handler *handler.AuthHandler) {
	auth := app.Group("/api/v1/auth")

	auth.Post("/login", handler.Login)
}
