package handlers

import (
	"go-login/db"
	"go-login/extensions"

	"github.com/gofiber/fiber/v2"
)

func GetAppVersion(c *fiber.Ctx) error {
	var version = 0
	row, err := db.DB().Query("Select version from app")
	for row.Next() {
		err := row.Scan(&version)
		if err != nil {
			return extensions.SendBadRequest(c, err.Error())
		}
	}

	if err != nil {
		return extensions.SendBadRequest(c, err.Error())
	}

	if version != 1 {
		return extensions.SendBadRequest(c, "Access Denied")
	}

	return extensions.SendSuccess(c, "OK")
}
