package model

type ImageModel struct {
	Id         int    `json:"id"`
	User_id    string `json:"user_id"`
	Image_path string `json:"image_path"`
	Type       int    `json:"type"`
}

type ImageUrl struct {
	Url string `json:"image"`
}

type Image struct {
	Image_path string `json:"image_path"`
	Type       string `json:"type"`
}
