package main

import (
	"database/sql"
	database_utils "pkg/database_utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, _ := sql.Open("sqlite3", "./new_db.db")
	database_utils.InitializeDatabase(db)
}
