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

	return c.JSON(model.StandardResponse{
		Success: true,
		Data:    assets,
		Meta: &model.ResponseMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
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

	return c.JSON(model.StandardResponse{
		Success: true,
		Data:    asset,
	})
}

func (h *AssetHandler) UpdateAsset(c *fiber.Ctx) error {
	employeeID := middleware.CurrentEmployeeID(c)
	if employeeID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "employee id is required",
		})
	}

	isHR, err := h.service.IsHREmployee(c.Context(), employeeID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "forbidden",
		})
	}

	if !isHR {
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

	return c.Status(fiber.StatusOK).JSON(model.StandardResponse{
		Success: true,
		Data:    asset,
	})
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
		var appErr *model.AppError
		if errors.As(err, &appErr) {
			return c.Status(appErr.Status).JSON(fiber.Map{
				"error": appErr.Message,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to delete asset",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
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
