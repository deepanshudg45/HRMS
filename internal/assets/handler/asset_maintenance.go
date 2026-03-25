package handler

import (
    "context"
    "errors"

    "WITS/internal/assets/model"
    "WITS/internal/assets/service"

    "github.com/gofiber/fiber/v2"
)

type AssetService interface {
    CreateAsset(ctx context.Context, req model.CreateAssetRequest) (*model.AssetDTO, error)
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
