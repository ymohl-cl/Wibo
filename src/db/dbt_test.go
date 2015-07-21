package db

import (
	//	"database/sql"
	"fmt"
	"github.com/Wibo/src/db"
	"testing"
)

func dbTest(t *testing.T) {
	var err error
	Db, err := db.OpenCo(err)
	ok, err := Db.Query("SELECT * FROM \"user\";")
	if err != nil {
		fmt.Println("pq error:", err)
	}
	fmt.Println(ok)
	t.Error(ok)
}
