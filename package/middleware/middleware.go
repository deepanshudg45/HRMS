package middleware

import "github.com/gofiber/fiber/v2"

func HealthCheck() fiber.Handler {
    return func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "ok",
        })
    }
}
