package routes

import (
	"WITS/internal/assets/handler"

	"github.com/gofiber/fiber/v2"
)

func AssetRoutes(app *fiber.App, handler *handler.AssetHandler) {
	asset := app.Group("/api/v1")

	asset.Get("/assets", handler.GetAssets)
	asset.Post("/assets", handler.CreateAsset)
	asset.Get("/assets/:id", handler.GetAssetByID)
	asset.Put("/assets/:id", handler.UpdateAsset)
	asset.Delete("/assets/:id", handler.DeleteAsset)
	asset.Post("/assets/:id/assign", handler.AssignAsset)
	asset.Get("/assets/my", handler.GetActiveAssets)
	asset.Patch("/assets/:id/status", handler.UpdateAssetStatus)
	asset.Post("/assets/:id/return", handler.ReturnAsset)
	asset.Get("/assets/employee/:employeeId", handler.GetAssetsByEmployeeID)
	asset.Get("/assets/:id/assignments", handler.GetMyAssignments)
	asset.Get("/assets/:id/maintenance", handler.GetMaintenanceRecords)
	asset.Post("/assets/:id/maintenance", handler.CreateMaintenanceRecord)
	asset.Patch("/assets/:id/maintenance/:maintenanceId", handler.UpdateMaintenanceRecord)
	asset.Patch("/assets/assignments/:id/acknowledge", handler.UpdateAssignment)
	asset.Get("/assets/admin/reports", handler.GenerateReport)
	asset.Post("/assets/import", handler.ImportAssets)
}
