package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"go-login/helpers"
	"go-login/model"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ImageUpload(c *fiber.Ctx) error {

	url := model.ImageUrl{}
	b := make([]byte, 4) //equals 8 characters
	rand.Read(b)
	var imageName = hex.EncodeToString(b)

	err := c.BodyParser(&url)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"message": err,
			"success": false,
		})
	}
	var strImage = strings.Replace(url.Url, "data:image/png;base64,", "", -1)
	println(strImage)
	unbased, err := base64.RawStdEncoding.DecodeString(strImage)
	if err != nil {
		panic("Cannot decode b64")
	}
	r := bytes.NewReader(unbased)
	typeString := helpers.GetStringInBetween(url.Url, ":", ";")
	println(typeString)
	switch typeString {
	case "image/jpeg":
		im, err := jpeg.Decode(r)
		if err != nil {
			println(err)
		}

		f, err := os.OpenFile("images/"+imageName+".jpeg", os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic("Cannot open file")
		}

		jpeg.Encode(f, im, nil)

	case "image/png":
		im, err := png.Decode(r)
		if err != nil {
			println(err)
		}

		f, err := os.OpenFile("images/"+imageName+".png", os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic("Cannot open file")
		}

		png.Encode(f, im)
		println(f.Name())
	}

	return c.Status(200).JSON(&fiber.Map{
		"message": "Başarılı",
		"path":    "",
	})
}
