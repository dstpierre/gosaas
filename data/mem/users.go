// +build mem

package mem

import (
	"fmt"

	"github.com/dstpierre/gosaas/data/model"
)

type Users struct {
	users []model.Account
}

func (u *Users) SignUp(email, password string) (*model.Account, error) {
	acctID := len(u.users) + 1
	userID := acctID * 200

	acct := model.Account{
		ID:    acctID,
		Email: email,
		Users: []model.User{{
			Email:    email,
			ID:       userID,
			Password: password,
			Token:    model.NewToken(acctID),
		}},
	}

	u.users = append(u.users, acct)
	return &acct, nil
}

func (u *Users) AddToken(accountID, userID model.Key, name string) (*model.AccessToken, error) {
	tok := model.AccessToken{
		ID:    userID * 300,
		Name:  name,
		Token: model.NewToken(accountID),
	}

	for _, acct := range u.users {
		if acct.ID == accountID {
			for _, usr := range acct.Users {
				if usr.ID == userID {
					usr.AccessTokens = append(usr.AccessTokens, tok)
					return &tok, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("unable to find account %d and user %d", accountID, userID)
}

func (u *Users) RemoveToken(accountID, userID, tokenID model.Key) error {
	for _, acct := range u.users {
		if acct.ID == accountID {
			for _, usr := range acct.Users {
				if usr.ID == userID {
					// cheap but does not need more here :p
					usr.AccessTokens = make([]model.AccessToken, 0)
					break
				}
			}
		}
	}
	return nil
}

func (u *Users) GetDetail(id model.Key) (*model.Account, error) {
	var user model.Account
	for _, acct := range u.users {
		if acct.ID == id {
			user = acct
			break
		}
	}
	return &user, nil
}

func (u *Users) RefreshSession(conn *bool, dbName string) {
	u.users = append(u.users, model.Account{
		ID:    1,
		Email: "test@domain.com",
		Users: []model.User{{
			Email:    "test@domain.com",
			ID:       1,
			Password: "unittest",
			Token:    model.NewToken(1),
		}},
	})
}
