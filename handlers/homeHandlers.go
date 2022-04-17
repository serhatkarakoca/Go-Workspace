package handlers

import (
	"go-login/db"
	"go-login/extensions"
	"go-login/model"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {

	//image := model.Images{}
	registerModel := model.Users{}
	var userArray []model.Users
	//user := model.Users{}
	var imagesArray []model.Image
	imageModel := model.Image{}
	favorites := model.ModelFav{}
	var favoritesArray []string
	var rejectedArray []string
	var distance = 0.0

	favoritesRows, err := db.DB().Query("SELECT favorite_id FROM favorites where user_id = ?", c.FormValue("username"))
	if err != nil {
		return extensions.SendBadRequest(c, err.Error())

	}
	for favoritesRows.Next() {
		err := favoritesRows.Scan(&favorites.Favorite_id)

		if err == nil {
			favoritesArray = append(favoritesArray, favorites.Favorite_id)
		}

	}
	rejectedRows, err := db.DB().Query("select rejected_id from rejected where user_id = ?", c.FormValue("username"))
	if err != nil {
		return extensions.SendBadRequest(c, err.Error())
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
	smtp, err := db.DB().Prepare("SELECT u.first_name,u.last_name,u.email,u.phone_number,u.gender,u.latitude,u.longitude,u.birthday,GROUP_CONCAT(b.image_path ORDER BY b.id) as image_path,GROUP_CONCAT(b.type ORDER BY b.id) as type, 6371 * 2 * ASIN(SQRT(POWER(SIN((? - abs(latitude)) * pi()/180 / 2), 2)+ COS(? * pi()/180 ) * COS(abs(latitude) * pi()/180)* POWER(SIN((? - longitude) * pi()/180 / 2), 2) )) as  distance FROM users as u left join images as b on b.user_id = u.id  where email not in (?" + strings.Repeat(",?", len(favoritesArray)-1) + ") and gender = ? Group by u.id HAVING distance < ? ORDER BY Rand() limit 0,20")
	if err != nil {
		return extensions.SendBadRequest(c, err.Error())
	}
	rows, err := smtp.Query(args...)
	if err != nil {
		return extensions.SendBadRequest(c, err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&registerModel.First_name, &registerModel.Last_name,
			&registerModel.Email, &registerModel.Phone_number, &registerModel.Gender,
			&registerModel.Latitude, &registerModel.Longitude, &registerModel.Birthday, &imageModel.Image_path, &imageModel.Type, &distance)
		if err != nil {
			return extensions.SendBadRequest(c, err.Error())
		}

		if int(math.Round(distance)) < 1 {
			distance = 1.0
		}
		registerModel.Distance = strconv.Itoa(int(math.Round(distance)))
		paths := strings.Split(imageModel.Image_path, ",")
		types := strings.Split(imageModel.Type, ",")

		for index, element := range paths {
			var img = model.Image{}
			img.Image_path = element
			img.Type = types[index]
			imagesArray = append(imagesArray, img)
		}
		registerModel.Images = imagesArray
		userArray = append(userArray, registerModel)
	}
	return c.Status(200).JSON(&fiber.Map{
		"size":  len(userArray),
		"users": userArray,
	})
}
