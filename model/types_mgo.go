// +build !mem

package model

import (
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Connection is the database connection type.
type Connection = mgo.Session

// Key is the primary key type.
type Key = bson.ObjectId

// Open creates a new MongoDB connection.
func Open(options ...string) (*mgo.Session, error) {
	conn, err := mgo.Dial(options[1])
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// KeyToString converts a Key to a string.
func KeyToString(id Key) string {
	return id.Hex()
}

// StringToKey converts a string to a Key.
func StringToKey(id string) Key {
	return bson.ObjectIdHex(id)
}

// NewID returns a new Key.
func NewID() Key {
	return bson.NewObjectId()
}
