package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CurrentEmployeeID(c *fiber.Ctx) string {
	employeeID := strings.TrimSpace(c.Get("X-Employee-ID"))
	if employeeID == "" {
		employeeID = strings.TrimSpace(c.Query("employeeId"))
	}
	return employeeID
}

func CurrentRole(c *fiber.Ctx) string {
	return strings.ToUpper(strings.TrimSpace(c.Get("X-Role")))
}

func IsHR(c *fiber.Ctx) bool {
	return CurrentRole(c) == "HR"
}

func CanAccessEmployee(c *fiber.Ctx, employeeID string) bool {
	employeeID = strings.TrimSpace(employeeID)
	if employeeID == "" {
		return false
	}

	if IsHR(c) {
		return true
	}

	return CurrentEmployeeID(c) == employeeID
}
