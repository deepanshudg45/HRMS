package handler

import (
	"context"
	"errors"
	"strconv"

	"WITS/internal/assets/model"
	"WITS/internal/assets/service"

	"github.com/gofiber/fiber/v2"
)

type AssetService interface {
    CreateAsset(ctx context.Context, req model.CreateAssetRequest) (*model.AssetDTO, error)
    GetAssets(ctx context.Context, filters *model.AssetFilter) ([]model.AssetListDTO, int, error)
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
        return c.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    asset, err := h.service.CreateAsset(c.Context(), req)
    if err != nil {
        if errors.Is(err, service.ErrMissingRequiredFields) || errors.Is(err, service.ErrInvalidAssetType) {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": err.Error(),
            })
        }

        return c.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(201).JSON(asset)
}

// GetAssets handles GET /assets with query filters
func (h *AssetHandler) GetAssets(c *fiber.Ctx) error {
    ctx := context.Background()

    // Parse query parameters
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

    // Build filter
    filters := &model.AssetFilter{
        Status:   status,
        Type:     assetType,
        Category: category,
        Search:   search,
        Page:     page,
        Limit:    limit,
        Offset:   offset,
    }

    // Get assets from service
    assets, total, err := h.service.GetAssets(ctx, filters)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "error":   err.Error(),
        })
    }

    if assets == nil {
        assets = []model.AssetListDTO{}
    }

    // Calculate pagination metadata
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
