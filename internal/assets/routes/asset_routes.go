package routes

import (
	"WITS/internal/assets/handler"
	"WITS/package/middleware"

	"github.com/gofiber/fiber/v2"
)

func AssetRoutes(app *fiber.App, handler *handler.AssetHandler) {
	asset := app.Group("/api/v1")

	asset.Get("/assets", handler.GetAssets)
	asset.Get("/assets/:id", handler.GetAssetByID)

	protected := asset.Group("/")
	protected.Use(middleware.JWTProtected())

	protected.Post("/assets", handler.CreateAsset)
	protected.Put("/assets/:id", handler.UpdateAsset)
	protected.Delete("/assets/:id", handler.DeleteAsset)
	protected.Post("/assets/:id/assign", handler.AssignAsset)
	protected.Get("/assets/my", handler.GetActiveAssets)
	protected.Patch("/assets/:id/status", handler.UpdateAssetStatus)
	protected.Post("/assets/:id/return", handler.ReturnAsset)
	protected.Get("/assets/employee/:employeeId", handler.GetAssetsByEmployeeID)
	protected.Get("/assets/:id/assignments", handler.GetMyAssignments)
	protected.Get("/assets/:id/maintenance", handler.GetMaintenanceRecords)
	protected.Post("/assets/:id/maintenance", handler.CreateMaintenanceRecord)
	protected.Patch("/assets/:id/maintenance/:maintenanceId", handler.UpdateMaintenanceRecord)
	protected.Patch("/assets/assignments/:id/acknowledge", handler.UpdateAssignment)
	protected.Get("/assets/admin/reports", handler.GenerateReport)
	protected.Post("/assets/import", handler.ImportAssets)
}
