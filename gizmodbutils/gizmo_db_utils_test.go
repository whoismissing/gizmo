package gizmodbutils

import (
	"testing"
	"database/sql"
)

func TestInitializeDatabase(t *testing.T) {
	db, _ := sql.Open("sqlite3", "./test_db.db")

	InitializeDatabase(db)
}
