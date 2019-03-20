package postgres

import (
	"testing"

	"github.com/dstpierre/gosaas/model"
)

func createAccountAndUser(t *testing.T, users *Users, email, pass string) *model.Account {
	acct, err := users.SignUp(email, pass)
	if err != nil {
		t.Fatalf("error on signul: %v", err)
	} else if len(acct.Users) != 1 {
		t.Fatalf("expected 1 user got: %d", len(acct.Users))
	}

	return acct
}

func TestUsersSignUp(t *testing.T) {
	t.Parallel()

	users := &Users{DB: db}
	acct := createAccountAndUser(t, users, "unit@test.com", "1234")
	if acct.Email != "unit@test.com" {
		t.Errorf("expected unit@test.com as email and got: %s", acct.Email)
	}
}

func TestUsersAuth(t *testing.T) {
	t.Parallel()

	users := &Users{DB: db}
	acct := createAccountAndUser(t, users, "auth@unittest.com", "1234")

	id, tok := model.ParseToken(acct.Users[0].Token)
	_, u, err := users.Auth(id, tok, false)
	if err != nil {
		t.Error(err)
	} else if u.ID != acct.Users[0].ID {
		t.Errorf("expected %d as user id got: %d", acct.Users[0].ID, u.ID)
	}
}

func TestUsersGetByEmail(t *testing.T) {
	t.Parallel()

	users := &Users{DB: db}
	acct := createAccountAndUser(t, users, "mymail@unittest.com", "1234")

	u, err := users.GetUserByEmail(acct.Email)
	if err != nil {
		t.Error(err)
	} else if acct.Users[0].ID != u.ID {
		t.Errorf("expected user id %d got: %d", acct.Users[0].ID, u.ID)
	}
}

func TestUsersPaid(t *testing.T) {
	t.Parallel()

	users := &Users{DB: db}
	acct := createAccountAndUser(t, users, "stripe@unittest.com", "1234")

	err := users.ConvertToPaid(acct.ID, "stripe_id_here", "sub_id_here", "planA", false, 1)
	if err != nil {
		t.Fatal(err)
	}

	a, err := users.GetByStripe("stripe_id_here")
	if err != nil {
		t.Error(err)
	} else if a.ID != acct.ID {
		t.Errorf("expected account id %d got: %d", acct.ID, a.ID)
	}
}

func TestUsersCancel(t *testing.T) {
	t.Parallel()

	users := &Users{DB: db}
	acct := createAccountAndUser(t, users, "cancel@unittest.com", "1234")

	err := users.ConvertToPaid(acct.ID, "stripe", "sub", "p", true, 1)
	if err != nil {
		t.Fatal(err)
	}

	err = users.Cancel(acct.ID)
	if err != nil {
		t.Fatal(err)
	}

	check, err := users.GetDetail(acct.ID)
	if err != nil {
		t.Error(err)
	} else if check.SubscriptionID != "" {
		t.Errorf("expected sub id to be '' got: %s", check.SubscribedOn)
	}
}
