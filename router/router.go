package router

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"go-login/handlers"
	"log"
)

func Routers() {

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋! \n please try /login or /register")
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		return handlers.Register(c)
	})

	app.Get("/appVersion", func(c *fiber.Ctx) error {
		return handlers.GetAppVersion(c)
	})

	app.Post("/addFavorite", func(c *fiber.Ctx) error {
		return handlers.AddToFavorite(c)
	})

	app.Post("/addRejected", func(c *fiber.Ctx) error {
		return handlers.AddToRejected(c)
	})

	app.Post("/imageUpload", func(c *fiber.Ctx) error {
		return handlers.ImageUpload(c)
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return handlers.Login(c)
	})

	app.Static("/images", "./images")

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("my_secret_key"),
	}))

	app.Get("/getUsers", func(c *fiber.Ctx) error {
		return handlers.GetUsers(c)
	})

	log.Fatal(app.Listen(":80"))
}
