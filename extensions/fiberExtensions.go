package extensions

import (
	"github.com/gofiber/fiber/v2"
)

func SendSuccess(c *fiber.Ctx, message string) error {
	return c.Status(200).JSON(&fiber.Map{
		"message": message,
		"success": true,
	})
}

func SendBadRequest(c *fiber.Ctx, message string) error {
	return c.Status(400).JSON(&fiber.Map{
		"message": message,
		"success": false,
	})
}
