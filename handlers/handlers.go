package handlers

import (
	"database/sql"
	"go-login/db"
	"go-login/helpers"
	"go-login/model"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Handlers() {

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		registerModel := model.Register{}
		var id = 0
		if err := c.BodyParser(&registerModel); err != nil {
			return err
		}

		if helpers.IsEmpty(registerModel.Email) || helpers.IsEmpty(registerModel.Password) {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": "You must fill in all fields",
			})
		}

		if !helpers.IsEmail(registerModel.Email) {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": "Please enter valid email",
			})
		}

		var passwordValidation = helpers.IsPasswordStrong(registerModel.Password)
		if !passwordValidation.Success {
			return c.Status(400).JSON(passwordValidation)
		}

		row, err := db.DB().Query("Select * from user where email = ?", registerModel.Email)
		if err != nil {
			log.Fatal(err)
		}
		defer row.Close()

		for row.Next() {
			switch err := row.Scan(&id, &registerModel.Email, &registerModel.Password); err {
			case sql.ErrNoRows:
				return c.Status(500).JSON(&fiber.Map{
					"success": false,
					"message": err,
				})
			case nil:
				return c.Status(400).JSON(&fiber.Map{
					"success": false,
					"message": "There is a user with this email please log in",
				})
			default:
				return c.Status(500).JSON(&fiber.Map{
					"success": false,
					"message": err,
				})
			}
		}

		hashPassword, err := helpers.HashPassword(registerModel.Password)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.DB().Exec("INSERT INTO user(email,password) VALUES(?,?)", registerModel.Email, hashPassword)
		if err != nil {
			log.Fatal(err)
		}

		return c.Status(200).JSON(registerModel)
	})

	app.Get("/login", func(c *fiber.Ctx) error {

		return c.JSON("Login")
	})

	log.Fatal(app.Listen(":80"))
}
