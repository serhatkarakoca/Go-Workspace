package helpers

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Username string `json:"email"`
	jwt.StandardClaims
}

type jwtModel struct {
	Access_token  string
	Refresh_token string
}

func GetToken(username string) jwtModel {
	expirationTime := time.Now().Add(1000 * time.Minute)
	var jwtKey = []byte("my_secret_key")
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Fatal(err)
	}
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	rt, err := refreshToken.SignedString([]byte("my_secret_key"))
	if err != nil {
		return jwtModel{}
	}

	return jwtModel{
		Access_token:  tokenString,
		Refresh_token: rt,
	}
}
