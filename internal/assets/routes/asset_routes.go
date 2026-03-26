package routes

import (
	"WITS/internal/assets/handler"

	"github.com/gofiber/fiber/v2"
)

func AssetRoutes(app *fiber.App, handler *handler.AssetHandler) {
    asset := app.Group("/api/v1")
    
    // Asset inventory routes
    asset.Get("/assets", handler.GetAssets)
    asset.Post("/assets", handler.CreateAsset)
    
    // Other commented endpoints for future implementation
    // asset.Delete("/assets/:id", handler.DeleteAsset)
    // asset.Patch("/assets/:id/status", handler.UpdateAssetStatus)
    // asset.Get("/assets/my",handler.GetActiveAssets)
    // asset.Post("/assets/:id/return", handler.ReturnAsset)
    // asset.Patch("/assets/assignments/:id/acknowledge",handler.UpdateAssignment)
    // asset.Get("/assets/:id/assignments", handler.GetMyAssignments)
    // asset.Get("/assets/:id/maintenance", handler.GetMaintenanceRecords)
    // asset.Post("/assets/import", handler.ImportAssets)
    // asset.Get("/assets/admin/reports", handler.GenerateReport)
    // asset.Get("/assets/employee/:employeeId", handler.GetAssetsByEmployeeID)
}