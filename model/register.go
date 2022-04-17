package model

type Register struct {
	Id            int    `json:"id"`
	First_name    string `json:"first_name"`
	Last_name     string `json:"last_name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Gender        int16  `json:"gender"`
	Interested_in int16  `json:"interested_in"`
	Phone_number  string `json:"phone_number"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	Birthday      string `json:"birthday"`
}
