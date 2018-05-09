// +build integration
// +build !mgo
// +build !mem

package data

import (
	"testing"
)

func Test_DB_Open(t *testing.T) {
	db := DB{}
	if err := db.Open("postgres", "postgres://postgres:dbpwd@localhost/gosaas?sslmode=disable"); err != nil {
		t.Fatal("unable to connect to postgres", err)
	}
}
