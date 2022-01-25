package main

import (
	"go-login/db"
	"go-login/handlers"
)

func main() {
	db.Connection()
	handlers.Handlers()
}
