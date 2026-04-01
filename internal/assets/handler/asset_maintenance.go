package handler

import (
	"context"

	"WITS/internal/assets/model"
	"WITS/package/middleware"

	"github.com/gofiber/fiber/v2"
)

type AssetService interface {
	CreateAsset(ctx context.Context, req model.CreateAssetRequest) (*model.AssetDTO, error)
	GetAssets(ctx context.Context, filters *model.AssetFilter) ([]model.AssetListDTO, int, error)
	GetAssetByID(ctx context.Context, assetID string) (*model.AssetDetailDTO, error)
	UpdateAsset(ctx context.Context, assetID string, req model.UpdateAssetRequest) (*model.AssetDTO, error)
	GetMaintenanceRecordsByAssetID(ctx context.Context, assetID string) ([]model.MaintenanceDTO, error)
	CreateMaintenanceRecord(ctx context.Context, assetID string, req model.MaintenanceRequest) (*model.MaintenanceDTO, error)
	UpdateMaintenanceRecord(ctx context.Context, assetID string, maintenanceID string, req model.UpdateMaintenanceRequest) (*model.MaintenanceDTO, error)
	DeleteAsset(ctx context.Context, assetID string) error
	AssignAsset(ctx context.Context, assetID string, req model.AssignRequest) (*model.AssignAssetDTO, error)
	UpdateAssetStatus(ctx context.Context, assetID string, req model.StatusRequest) (*model.AssetDTO, error)
	GetActiveAssets(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error)
	ReturnAsset(ctx context.Context, assetID string, req model.ReturnRequest) (*model.ReturnDTO, error)
	UpdateAssignment(ctx context.Context, assignmentID string, employeeID string) (*model.AssignmentDTO, error)
	GetMyAssignments(ctx context.Context, assetID string) ([]model.AssignmentHistoryDTO, error)
	GetAssetsByEmployeeID(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error)
	GenerateReport(ctx context.Context, filters *model.AssetFilter) ([]model.AssetReportDTO, error)
	ImportAssets(ctx context.Context, requests []model.CreateAssetRequest) (*model.ImportResultDTO, error)
}

type AssetHandler struct {
	service AssetService
}

func NewAssetHandler(service AssetService) *AssetHandler {
	return &AssetHandler{service: service}
}

func (h *AssetHandler) GetMaintenanceRecords(c *fiber.Ctx) error {
	if !middleware.IsHR(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "forbidden",
		})
	}

	assetID := c.Params("id")
	if assetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "asset id is required",
		})
	}

	records, err := h.service.GetMaintenanceRecordsByAssetID(c.Context(), assetID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to fetch maintenance records",
		})
	}

	if records == nil {
		records = []model.MaintenanceDTO{}
	}

	return c.Status(fiber.StatusOK).JSON(records)
}

func (h *AssetHandler) CreateMaintenanceRecord(c *fiber.Ctx) error {
	if !middleware.IsHR(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "forbidden",
		})
	}

	assetID := c.Params("id")
	if assetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "asset id is required",
		})
	}

	var req model.MaintenanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	record, err := h.service.CreateMaintenanceRecord(c.Context(), assetID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to create maintenance record",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(record)
}

func (h *AssetHandler) UpdateMaintenanceRecord(c *fiber.Ctx) error {
	if !middleware.IsHR(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "forbidden",
		})
	}

	assetID := c.Params("id")
	maintenanceID := c.Params("maintenanceId")
	if assetID == "" || maintenanceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "asset id and maintenance id are required",
		})
	}

	var req model.UpdateMaintenanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	record, err := h.service.UpdateMaintenanceRecord(c.Context(), assetID, maintenanceID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to update maintenance record",
		})
	}

	return c.Status(fiber.StatusOK).JSON(record)
}
