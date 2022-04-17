package handlers

import (
	"go-login/db"
	"go-login/extensions"
	"go-login/model"

	"github.com/gofiber/fiber/v2"
)

func AddToFavorite(c *fiber.Ctx) error {
	favoriteModel := model.ModelFav{}

	if err := c.BodyParser(&favoriteModel); err != nil {
		return extensions.SendBadRequest(c, err.Error())
	}

	_, err := db.DB().Exec("INSERT INTO favorites (user_id,favorite_id) Values(?,?)", favoriteModel.User_id, favoriteModel.Favorite_id)

	if err != nil {
		return extensions.SendBadRequest(c, err.Error())
	}
	return c.Status(200).JSON(&fiber.Map{
		"message": "İşlem basarili",
		"success": true,
	})
}

func AddToRejected(c *fiber.Ctx) error {
	favoriteModel := model.ModelFav{}

	if err := c.BodyParser(&favoriteModel); err != nil {
		return err
	}

	_, err := db.DB().Exec("INSERT INTO rejected (user_id,rejected_id) Values(?,?)", favoriteModel.User_id, favoriteModel.Favorite_id)

	if err != nil {
		return extensions.SendBadRequest(c, err.Error())
	}
	return c.Status(200).JSON(&fiber.Map{
		"message": "İşlem basarili",
		"success": true,
	})
}
