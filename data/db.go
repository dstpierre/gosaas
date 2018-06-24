package data

import (
	"github.com/dstpierre/gosaas/data/model"
)

type DB struct {
	DatabaseName string
	Connection   *model.Connection
	CopySession  bool

	Users UserServices
}

type SessionRefresher interface {
	RefreshSession(*model.Connection, string)
}

type UserServices interface {
	SessionRefresher
	SignUp(email, password string) (*model.Account, error)
	AddToken(accountID, userID model.Key, name string) (*model.AccessToken, error)
	RemoveToken(accountID, userID, tokenID model.Key) error
	Auth(accountID, token string, pat bool) (*model.Account, *model.User, error)
	GetDetail(id model.Key) (*model.Account, error)
}

type AdminServices interface {
	SessionRefresher
	LogRequests(reqs []model.APIRequest) error
}
