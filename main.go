package main

import (
	"go-login/db"
	"go-login/router"
)

func main() {
	db.Connection()
	router.Routers()
}
