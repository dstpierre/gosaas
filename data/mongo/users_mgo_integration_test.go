// +build integration
// +build mgo

package mongo

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func Test_DB_Users_GetDetail(t *testing.T) {
	users := Users{}
	users.RefreshSession(db, dbName)

	user, err := users.GetDetail(bson.ObjectIdHex("5af1afc905bf597042d47d90"))
	if err != nil {
		t.Error("unable to get user details", err)
	} else if user == nil {
		t.Error("expected user to have a value got nil")
	} else if user.Email != "test@domain.com" {
		t.Error("expected email to be test@domain.com got ", user.Email)
	}
}
