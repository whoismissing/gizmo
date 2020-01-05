package main

import (
	"database/sql"
	db_utils "github.com/whoismissing/gizmo/gizmodbutils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, _ := sql.Open("sqlite3", "./gizmo.db")
	db_utils.InitializeDatabase(db)
}
