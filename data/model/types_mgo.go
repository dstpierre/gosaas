// +build mgo

package model

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
