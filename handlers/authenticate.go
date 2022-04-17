package handlers

import (
	"database/sql"
	"go-login/db"
	"go-login/extensions"
	"go-login/helpers"
	"go-login/model"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	registerModel := model.Register{}
	//var id = 0

	if helpers.IsEmpty(c.FormValue("email")) || helpers.IsEmpty(c.FormValue("password")) {
		return c.Status(400).JSON(&fiber.Map{
			"message": "You must fill in all fields",
			"success": false,
		})
	}
	rows, err := db.DB().Query("Select * from users where email = ?", c.FormValue("email"))
	if err != nil {
		log.Fatal(err)
		return c.Status(400).JSON(&fiber.Map{
			"message": err,
			"success": false,
		})
	}

	for rows.Next() {
		switch err := rows.Scan(&registerModel.Id, &registerModel.First_name, &registerModel.Last_name,
			&registerModel.Email, &registerModel.Password, &registerModel.Phone_number, &registerModel.Gender,
			&registerModel.Latitude, &registerModel.Longitude, &registerModel.Interested_in, &registerModel.Birthday); err {
		case sql.ErrNoRows:
			return c.Status(500).JSON("ErrNoRows")
		case nil:
			if helpers.CheckPasswordHash(c.FormValue("password"), registerModel.Password) {
				access_token := helpers.GetToken(registerModel.Email)

				return c.Status(201).JSON(&fiber.Map{
					"access_token":  access_token.Access_token,
					"refresh_token": access_token.Refresh_token,
					"createdTime":   time.Now(),
					"expires_time":  time.Now().Add(15 * time.Minute),
				})
			} else {
				return c.Status(400).JSON(&fiber.Map{
					"success":           false,
					"message":           "E-posta ve ya şifreniz yanlış",
					"error":             "BadCredentials",
					"error_description": "BadCredentials",
				})
			}

		default:
			return c.Status(500).JSON("Error")
		}
	}
	return c.Status(400).JSON("Account not exist")
}

func Register(c *fiber.Ctx) error {
	registerModel := model.Register{}
	var id = 0
	if err := c.BodyParser(&registerModel); err != nil {
		return extensions.SendBadRequest(c, err.Error())
	}

	if helpers.IsEmpty(registerModel.Email) || helpers.IsEmpty(registerModel.Password) {
		return extensions.SendBadRequest(c, "Please fill in the all fields !")
	}

	if !helpers.IsEmail(registerModel.Email) {
		return extensions.SendBadRequest(c, "Please enter valid email")
	}

	var passwordValidation = helpers.IsPasswordStrong(registerModel.Password)
	if !passwordValidation.Success {
		return c.Status(400).JSON(passwordValidation)
	}

	row, err := db.DB().Query("Select * from users where email = ?", registerModel.Email)
	if err != nil {
		return extensions.SendBadRequest(c, "There is no user with this email")
	}

	for row.Next() {
		switch err := row.Scan(&id, &registerModel.First_name, &registerModel.Last_name,
			&registerModel.Email, &registerModel.Password, &registerModel.Phone_number, &registerModel.Gender,
			&registerModel.Latitude, &registerModel.Longitude, &registerModel.Interested_in, &registerModel.Birthday); err {
		case sql.ErrNoRows:
			return c.Status(500).JSON(&fiber.Map{

				"success": false,
				"message": err,
			})
		case nil:
			return c.Status(400).JSON(&fiber.Map{
				"success":           false,
				"message":           "There is a user with this email please log in",
				"error":             "DuplicateUidError",
				"error_description": "DuplicateUidError",
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
		return extensions.SendBadRequest(c, err.Error())
	}

	_, err = db.DB().Exec("INSERT INTO users(first_name,last_name,email,password,phone_number,gender,latitude,longitude,interested_in,birthday) VALUES(?,?,?,?,?,?,?,?,?,?)", registerModel.First_name, registerModel.Last_name, registerModel.Email, hashPassword, registerModel.Phone_number, registerModel.Gender, registerModel.Latitude, registerModel.Longitude, registerModel.Interested_in, registerModel.Birthday)
	if err != nil {
		return extensions.SendBadRequest(c, err.Error())
	}

	return extensions.SendSuccess(c, "registered successfully")
}
