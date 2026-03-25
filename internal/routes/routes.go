package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"hrms-1/internal/assets"
	"hrms-1/pkg/db"
)


func AssetsRoutes(app *fiber.App){
	repo := &assets.Repository{
		DB: db.GetPool(),
	}
	// fmt.Println("DB pool:", db.GetPool())
	handler := &assets.Handler{
		Service: &assets.Service{},
		Repo:    repo,
	}
	app.Get("/assets_handler",handler.GetAssets)
}