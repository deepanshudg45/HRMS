package handler

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
	"strconv"
	"strings"

	"WITS/internal/assets/model"
	"WITS/package/middleware"

	"github.com/gofiber/fiber/v2"
)

type AssetService interface {
	CreateAsset(ctx context.Context, req model.CreateAssetRequest) (*model.AssetDTO, error)
	GetAssets(ctx context.Context, filters *model.AssetFilter) ([]model.AssetListDTO, int, error)
	GetAssetByID(ctx context.Context, assetID string) (*model.AssetDetailDTO, error)
	UpdateAsset(ctx context.Context, assetID string, req model.UpdateAssetRequest) (*model.AssetDetailDTO, error)
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

func (h *AssetHandler) CreateAsset(c *fiber.Ctx) error {
	var req model.CreateAssetRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	asset, err := h.service.CreateAsset(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to create asset",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(asset)
}

func (h *AssetHandler) GetAssets(c *fiber.Ctx) error {
	ctx := context.Background()

	status := c.Query("status", "")
	assetType := c.Query("type", "")
	category := c.Query("category", "")
	search := c.Query("search", "")

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	filters := &model.AssetFilter{
		Status:   status,
		Type:     assetType,
		Category: category,
		Search:   search,
		Page:     page,
		Limit:    limit,
		Offset:   offset,
	}

	assets, total, err := h.service.GetAssets(ctx, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to fetch assets",
		})
	}

	if assets == nil {
		assets = []model.AssetListDTO{}
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"success": true,
		"data":    assets,
		"meta": fiber.Map{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

func (h *AssetHandler) GetAssetByID(c *fiber.Ctx) error {
	assetID := c.Params("id")
	if assetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "asset id is required",
		})
	}

	asset, err := h.service.GetAssetByID(c.Context(), assetID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "asset not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    asset,
	})
}

func (h *AssetHandler) UpdateAsset(c *fiber.Ctx) error {
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

	var req model.UpdateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	asset, err := h.service.UpdateAsset(c.Context(), assetID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to update asset",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    asset,
	})
}

func (h *AssetHandler) AssignAsset(c *fiber.Ctx) error {
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

	var req model.AssignRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := h.service.AssignAsset(c.Context(), assetID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to assign asset",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
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

func (h *AssetHandler) DeleteAsset(c *fiber.Ctx) error {
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

	err := h.service.DeleteAsset(c.Context(), assetID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to delete asset",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AssetHandler) UpdateAssetStatus(c *fiber.Ctx) error {
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

	var req model.StatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	asset, err := h.service.UpdateAssetStatus(c.Context(), assetID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to update asset status",
		})
	}

	return c.Status(fiber.StatusOK).JSON(asset)
}

func (h *AssetHandler) GetActiveAssets(c *fiber.Ctx) error {
	employeeID := middleware.CurrentEmployeeID(c)

	assets, err := h.service.GetActiveAssets(c.Context(), employeeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to fetch assets",
		})
	}

	if assets == nil {
		assets = []model.MyAssetDTO{}
	}

	return c.Status(fiber.StatusOK).JSON(assets)
}

func (h *AssetHandler) ReturnAsset(c *fiber.Ctx) error {
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

	var req model.ReturnRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := h.service.ReturnAsset(c.Context(), assetID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to return asset",
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *AssetHandler) UpdateAssignment(c *fiber.Ctx) error {
	assignmentID := c.Params("id")
	if assignmentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "assignment id is required",
		})
	}

	employeeID := middleware.CurrentEmployeeID(c)

	result, err := h.service.UpdateAssignment(c.Context(), assignmentID, employeeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to update assignment",
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *AssetHandler) GetMyAssignments(c *fiber.Ctx) error {
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

	assignments, err := h.service.GetMyAssignments(c.Context(), assetID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to fetch assignments",
		})
	}

	if assignments == nil {
		assignments = []model.AssignmentHistoryDTO{}
	}

	return c.Status(fiber.StatusOK).JSON(assignments)
}

func (h *AssetHandler) GetAssetsByEmployeeID(c *fiber.Ctx) error {
	employeeID := c.Params("employeeId")
	if employeeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "employee id is required",
		})
	}

	if !middleware.CanAccessEmployee(c, employeeID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "forbidden",
		})
	}

	assets, err := h.service.GetAssetsByEmployeeID(c.Context(), employeeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to fetch employee assets",
		})
	}

	if assets == nil {
		assets = []model.MyAssetDTO{}
	}

	return c.Status(fiber.StatusOK).JSON(assets)
}

func (h *AssetHandler) GenerateReport(c *fiber.Ctx) error {
	if !middleware.IsHR(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "forbidden",
		})
	}

	filters := &model.AssetFilter{
		Status:   c.Query("status", ""),
		Type:     c.Query("type", ""),
		Category: c.Query("category", ""),
		Search:   c.Query("search", ""),
	}

	reports, err := h.service.GenerateReport(c.Context(), filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to generate report",
		})
	}

	if reports == nil {
		reports = []model.AssetReportDTO{}
	}

	return c.Status(fiber.StatusOK).JSON(reports)
}

func (h *AssetHandler) ImportAssets(c *fiber.Ctx) error {
	if !middleware.IsHR(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "forbidden",
		})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "csv file is required",
		})
	}

	requests, parseErrors, err := parseAssetImportCSV(fileHeader)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to read csv file",
		})
	}

	result, err := h.service.ImportAssets(c.Context(), requests)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to import assets",
		})
	}

	if len(parseErrors) > 0 {
		result.Skipped += len(parseErrors)
		result.Errors = append(result.Errors, parseErrors...)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func parseAssetImportCSV(fileHeader *multipart.FileHeader) ([]model.CreateAssetRequest, []model.RowError, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	header, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil, errors.New("csv file is empty")
		}
		return nil, nil, err
	}

	indexMap := map[string]int{}
	for i, col := range header {
		normalized := normalizeCSVHeader(col)
		indexMap[normalized] = i
	}

	requests := []model.CreateAssetRequest{}
	rowErrors := []model.RowError{}
	rowNumber := 1

	for {
		rowNumber++
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			rowErrors = append(rowErrors, model.RowError{
				Row:    rowNumber,
				Reason: err.Error(),
			})
			continue
		}

		req, rowErr := buildCreateAssetRequest(record, indexMap)
		if rowErr != nil {
			rowErrors = append(rowErrors, model.RowError{
				Row:    rowNumber,
				Reason: rowErr.Error(),
			})
			continue
		}

		requests = append(requests, req)
	}

	return requests, rowErrors, nil
}

func buildCreateAssetRequest(record []string, indexMap map[string]int) (model.CreateAssetRequest, error) {
	get := func(keys ...string) string {
		for _, key := range keys {
			if idx, ok := indexMap[key]; ok && idx < len(record) {
				return strings.TrimSpace(record[idx])
			}
		}
		return ""
	}

	req := model.CreateAssetRequest{
		AssetType:      get("assettype", "asset_type", "type"),
		AssetName:      get("assetname", "asset_name", "name"),
		Brand:          get("brand"),
		Model:          get("model"),
		AssetCategory:  get("assetcategory", "asset_category", "category"),
		SerialNo:       get("serialno", "serial_no"),
		PurchaseDate:   get("purchasedate", "purchase_date"),
		Vendor:         get("vendor"),
		WarrantyExpiry: get("warrantyexpiry", "warranty_expiry"),
		Location:       get("location"),
		Notes:          get("notes"),
	}

	if req.AssetType == "" || req.AssetName == "" || req.AssetCategory == "" {
		return model.CreateAssetRequest{}, errors.New("assetType, assetName and assetCategory are required")
	}

	costValue := get("purchasecostinr", "purchase_cost_inr")
	if costValue != "" {
		cost, err := strconv.ParseFloat(costValue, 64)
		if err != nil {
			return model.CreateAssetRequest{}, errors.New("invalid purchaseCostINR")
		}
		req.PurchaseCostINR = cost
	}

	return req, nil
}

func normalizeCSVHeader(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "_", "")
	value = strings.ReplaceAll(value, " ", "")
	return value
}
