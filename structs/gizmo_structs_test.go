package structs

import (
	"testing"
)

var test_db = "./test_db.db"

func TestInitializeStructs(t *testing.T) {
	t.Logf("TestInitializeStructs SUCCESS")

	game := Game{teams: nil, time: 1}
	t.Logf("game = %+v\n", game)
}
