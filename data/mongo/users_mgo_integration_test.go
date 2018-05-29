// +build integration
// +build !mem

package mongo

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func Test_DB_Users_SignUp(t *testing.T) {
	users := Users{}
	users.RefreshSession(db, dbName)

	acct, err := users.SignUp("unit@test.com", "unittest")
	if err != nil {
		t.Error("unable to create new account", err)
	} else if acct.Email != "unit@test.com" {
		t.Error("account email is not correct expected test@unit.com got", acct.Email)
	} else if len(acct.Users) != 1 {
		t.Error("account has no users expected 1 got", len(acct.Users))
	} else if acct.Users[0].Email != "unit@test.com" {
		t.Error("user email is not correct expected test@unit.com got", acct.Users[0].Email)
	}
}

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
