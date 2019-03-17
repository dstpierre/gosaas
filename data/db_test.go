package data

import (
	"testing"
)

func Test_DB_Open(t *testing.T) {
	db := DB{}
	ds := "user=postgres password=postgres dbname=postgres sslmode=disable"

	if err := db.Open("postgres", ds); err != nil {
		t.Fatal("unable to connect to postgres", err)
	}
}
