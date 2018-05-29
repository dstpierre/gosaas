// +build integration
// +build !mem

package mongo

import (
	"log"
	"os"
	"testing"

	"github.com/dstpierre/gosaas/data/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

	// clean up
	if _, err = db.DB(dbName).C("users").RemoveAll(bson.M{"email": "unit@test.com"}); err != nil {
		log.Fatal("unable to remove integration test data: ", err)
	}

	os.Exit(m.Run())
}
