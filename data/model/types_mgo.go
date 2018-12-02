// +build !mem

package model

import (
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Connection = mgo.Session
type Key = bson.ObjectId

func Open(options ...string) (*mgo.Session, error) {
	conn, err := mgo.Dial(options[1])
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func keyToString(id Key) string {
	return id.Hex()
}
