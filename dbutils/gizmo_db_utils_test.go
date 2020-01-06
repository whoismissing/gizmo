package dbutils

import (
	"testing"
	"database/sql"
	"log"
	"os"
)

var test_db = "./test_db.db"

func TestInitializeDatabaseSuccess(t *testing.T) {
	t.Logf("Test InitializeDatabase SUCCESS")

	db, err := sql.Open("sqlite3", test_db)

	if err != nil {
		log.Fatal(err)
	}

	InitializeDatabase(db)

	err = os.Remove(test_db)

	if err != nil {
		log.Fatal(err)
	}
}
