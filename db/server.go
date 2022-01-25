package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Connection() {
	database, err := sql.Open("mysql", "root:12345678@/deneme")

	if err != nil {
		panic(err.Error())
	}
	fmt.Print("connected to database")
	db = database

}

func DB() *sql.DB {
	return db
}
