// +build integration
// +build !mgo
// +build !mem

package pg

import (
	"testing"
)

func Test_DB_Users_GetDetail(t *testing.T) {
	users := Users{DB: db}
	user, err := users.GetDetail(1)
	if err != nil {
		t.Error("unable to get user details", err)
	} else if user == nil {
		t.Error("expected user to have a value got nil")
	} else if user.Email != "test@domain.com" {
		t.Error("expected email to be test@domain.com got ", user.Email)
	}
}
