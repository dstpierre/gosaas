// +build integration
// +build mgo

package mongo

import (
	"log"
	"os"
	"testing"

	"github.com/dstpierre/gosaas/data/model"
	mgo "gopkg.in/mgo.v2"
)

const (
	dbName = "gosaas"
)

var (
	db *mgo.Session
)

func TestMain(m *testing.M) {
	conn, err := model.Open("mongo", "127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}

	db = conn
	defer conn.Close()
	os.Exit(m.Run())
}
