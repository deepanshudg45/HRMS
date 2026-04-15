package handler

import (
	"WITS/internal/assets/model"
	"WITS/package/middleware"

	"github.com/gofiber/fiber/v2"
)

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
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (h *AssetHandler) GetActiveAssets(c *fiber.Ctx) error {
	employeeID := middleware.CurrentEmployeeID(c)

	if employeeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "employee id is required",
		})
	}

	assets, err := h.service.GetActiveAssets(c.Context(), employeeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to fetch assets",
		})
	}

	if assets == nil {
		assets = []model.MyAssetDTO{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    assets,
		"count":   len(assets),
	})
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
