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

func Test_DB_Users_AddToken(t *testing.T) {
	users := Users{}
	users.RefreshSession(db, dbName)

	tok, err := users.AddToken(testAccount.ID, testAccount.Users[0].ID, "unit test")
	if err != nil {
		t.Error("error while creating access token", err)
	} else if tok == nil {
		t.Error("returned nil for the token")
	} else if tok.Name != "unit test" {
		t.Error("expected name to be `unit test` got", tok.Name)
	}
}

func Test_DB_Users_Auth(t *testing.T) {
	users := Users{}
	users.RefreshSession(db, dbName)

	a, u, err := users.Auth(testAccount.ID.Hex(), testAccount.Users[0].Token, false)
	if err != nil {
		t.Error("error during authentication", err)
	} else if u.ID != testAccount.Users[0].ID {
		t.Errorf("expected user to be %s got %s", testAccount.Users[0].Email, u.Email)
	} else if a.ID != testAccount.ID {
		t.Errorf("expected account to be %v got %v", testAccount.ID, a.ID)
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
