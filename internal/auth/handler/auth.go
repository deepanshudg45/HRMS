package handler

import (
	"strings"

	authModel "WITS/internal/auth/model"
	"WITS/package/middleware"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	jwtSecret string
}

func NewAuthHandler(jwtSecret string) *AuthHandler {
	return &AuthHandler{jwtSecret: jwtSecret}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req authModel.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	req.EmployeeID = strings.TrimSpace(req.EmployeeID)
	req.Role = strings.ToUpper(strings.TrimSpace(req.Role))

	if req.EmployeeID == "" || req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "employeeId and role are required",
		})
	}

	token, err := middleware.GenerateJWT(req.EmployeeID, req.Role, h.jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(authModel.LoginResponse{
		Token: token,
		User: authModel.UserDTO{
			EmployeeID: req.EmployeeID,
			Role:       req.Role,
		},
	})
}
