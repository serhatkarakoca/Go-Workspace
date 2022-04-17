package handlers

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"go-login/db"
	"go-login/extensions"
	"go-login/helpers"
	"go-login/model"
	"image/png"
	"log"
	"math"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func handlers() {

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {

		return c.SendString("Hello, World 👋! \n please try /login or /register")
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

		row, err := db.DB().Query("Select * from users where email = ?", registerModel.Email)
		if err != nil {
			log.Fatal(err)
		}
		defer row.Close()

		for row.Next() {
			switch err := row.Scan(&id, &registerModel.First_name, &registerModel.Last_name,
				&registerModel.Email, &registerModel.Password, &registerModel.Phone_number, &registerModel.Gender,
				&registerModel.Latitude, &registerModel.Longitude, &registerModel.Interested_in); err {
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

		_, err = db.DB().Exec("INSERT INTO users(first_name,last_name,email,password,gender,interested_in) VALUES(?,?,?,?,?,?)", registerModel.First_name, registerModel.Last_name, registerModel.Email, hashPassword, registerModel.Gender, registerModel.Interested_in)
		if err != nil {
			log.Fatal(err)
		}

		return c.Status(200).JSON(registerModel)
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return Login(c)
	})

	app.Get("/appVersion", func(c *fiber.Ctx) error {
		_, err := db.DB().Query("Select version from app")
		if err != nil {
			return extensions.SendBadRequest(c, err.Error())
		}
		return extensions.SendSuccess(c, "OK")
	})

	app.Get("/getUsers", func(c *fiber.Ctx) error {

		registerModel := model.Register{}
		var userArray []model.Users
		user := model.Users{}

		favorites := model.ModelFav{}
		var favoritesArray []string
		var rejectedArray []string
		var distance = 0.0

		favoritesRows, err := db.DB().Query("SELECT favorite_id FROM favorites where user_id = ?", c.FormValue("username"))
		if err != nil {
			return c.JSON("")
		}
		for favoritesRows.Next() {
			err := favoritesRows.Scan(&favorites.Favorite_id)

			if err == nil {
				favoritesArray = append(favoritesArray, favorites.Favorite_id)
			}

		}
		rejectedRows, err := db.DB().Query("select rejected_id from rejected where user_id = ?", c.FormValue("username"))
		if err != nil {
			return c.JSON("")
		}
		for rejectedRows.Next() {
			err := rejectedRows.Scan(&favorites.Favorite_id)

			if err == nil {
				rejectedArray = append(rejectedArray, favorites.Favorite_id)
			}

		}
		favoritesArray = append(favoritesArray, rejectedArray...)
		favoritesArray = append(favoritesArray, c.FormValue("username"))
		var argument []string
		argument = append(argument, c.FormValue("latitude"), c.FormValue("latitude"), c.FormValue("longitude"))
		argument = append(argument, favoritesArray...)
		argument = append(argument, c.FormValue("gender"), c.FormValue("distance"))
		args := make([]interface{}, len(argument))
		for i, id := range argument {
			args[i] = id
		}
		smtp, err := db.DB().Prepare("SELECT *, 6371 * 2 * ASIN(SQRT(POWER(SIN((? - abs(latitude)) * pi()/180 / 2), 2)+ COS(? * pi()/180 ) * COS(abs(latitude) * pi()/180)* POWER(SIN((? - longitude) * pi()/180 / 2), 2) )) as  distance FROM users where email not in (?" + strings.Repeat(",?", len(favoritesArray)-1) + ") and gender = ? HAVING distance < ? ORDER BY Rand() limit 0,20")
		if err != nil {
			return err
		}
		rows, err := smtp.Query(args...)
		if err != nil {
			return err
		}
		for rows.Next() {
			err := rows.Scan(&registerModel.Id, &registerModel.First_name, &registerModel.Last_name,
				&registerModel.Email, &registerModel.Password, &registerModel.Phone_number, &registerModel.Gender,
				&registerModel.Latitude, &registerModel.Longitude, &registerModel.Interested_in, &distance)

			if err != nil {
				return err
			}

			user.Email = registerModel.Email
			//user.Images = nil
			if int(math.Round(distance)) < 1 {
				distance = 1.0
			}
			//	user.Distance = strconv.Itoa(int(math.Round(distance)))
			userArray = append(userArray, user)
		}
		return c.Status(200).JSON(&fiber.Map{
			"size":  len(userArray),
			"users": userArray,
		})
	})

	app.Post("/addFavorite", func(c *fiber.Ctx) error {
		favoriteModel := model.ModelFav{}

		if err := c.BodyParser(&favoriteModel); err != nil {
			return err
		}

		_, err := db.DB().Exec("INSERT INTO favorites (user_id,favorite_id) Values(?,?)", favoriteModel.User_id, favoriteModel.Favorite_id)

		if err != nil {
			return c.JSON("Hata")
		}
		return c.Status(200).JSON(&fiber.Map{
			"message": "İşlem başarili",
			"success": true,
		})
	})

	app.Post("/addRejected", func(c *fiber.Ctx) error {
		favoriteModel := model.ModelFav{}

		if err := c.BodyParser(&favoriteModel); err != nil {
			return err
		}

		_, err := db.DB().Exec("INSERT INTO rejected (user_id,rejected_id) Values(?,?)", favoriteModel.User_id, favoriteModel.Favorite_id)

		if err != nil {
			return c.JSON("Hata")
		}
		return c.Status(200).JSON(&fiber.Map{
			"message": "İşlem başarili",
			"success": true,
		})
	})

	app.Post("/imageUpload", func(c *fiber.Ctx) error {
		var bodyList string

		err := c.BodyParser(&bodyList)
		if err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"message": err,
				"success": false,
			})
		}

		unbased, err := base64.StdEncoding.DecodeString(bodyList)
		if err != nil {
			panic("Cannot decode b64")
		}
		r := bytes.NewReader(unbased)
		im, err := png.Decode(r)
		if err != nil {
			panic("Bad png")
		}

		f, err := os.OpenFile("example.png", os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic("Cannot open file")
		}

		png.Encode(f, im)

		return c.SendString("TODO")
	})

	app.Get("/getImages", func(c *fiber.Ctx) error {
		var images []string
		var image string

		rows, err := db.DB().Query("Select image_path from images where user_id = ?", c.FormValue("username"))

		if err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"message": "Something went wrong",
				"success": false,
			})
		}

		for rows.Next() {
			err := rows.Scan(&image)

			if err != nil {
				return c.Status(400).JSON(&fiber.Map{
					"message": "Something went wrong",
					"success": false,
				})
			}
			images = append(images, image)

		}
		return c.Status(200).JSON(&fiber.Map{
			"images":    images,
			"imageSize": len(images),
			"username":  c.FormValue("username"),
		})
	})
	log.Fatal(app.Listen(":80"))
}
