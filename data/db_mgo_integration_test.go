// +build mgo
// +build integration

package data

import (
	"testing"
)

func Test_DB_Open(t *testing.T) {
	db := DB{}
	if err := db.Open("mongo", "127.0.0.1"); err != nil {
		t.Fatal("unable to connect to mongo", err)
	}
}
