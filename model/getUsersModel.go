package model

type Users struct {
	Email         string  `json:"email"`
	First_name    string  `json:"first_name"`
	Last_name     string  `json:"last_name"`
	Gender        int16   `json:"gender"`
	Interested_in int16   `json:"interested_in"`
	Phone_number  string  `json:"phone_number"`
	Latitude      string  `json:"latitude"`
	Longitude     string  `json:"longitude"`
	Birthday      string  `json:"birthday"`
	Images        []Image `json:"images"`
	Distance      string  `json:"distance"`
}

type Images struct {
	Url string `json:"image_path"`
}
