package data

import (
	"github.com/dstpierre/gosaas/model"
)

// DB is a database agnostic abstraction that contains a reference to
// the database connection.
//
// At this moment MongoDB and an in-memory data provider are supported.
type DB struct {
	// DatabaseName is the name of the database used.
	DatabaseName string
	// Connection is the reference to the database connection.
	Connection *model.Connection
	// CopySession indicates if the database connection should be copy of each requests (used for MongoDB).
	CopySession bool

	// Users contains the data access functions related to account, user and billing.
	Users UserServices
	// Webhooks contains the data access functions related to managing Webhooks.
	Webhooks WebhookServices
}

// SessionRefresher is an interface that contains a function to copy the database session.
type SessionRefresher interface {
	RefreshSession(*model.Connection, string)
}

// SessionCloser is an interface that contains a function to close the database session.
type SessionCloser interface {
	Close()
}

// UserServices is an interface that contians all functions related to account, user and billing.
type UserServices interface {
	SessionRefresher
	SessionCloser
	SignUp(email, password string) (*model.Account, error)
	AddToken(accountID, userID model.Key, name string) (*model.AccessToken, error)
	RemoveToken(accountID, userID, tokenID model.Key) error
	Auth(accountID model.Key, token string, pat bool) (*model.Account, *model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetDetail(id model.Key) (*model.Account, error)
	GetByStripe(stripeId string) (*model.Account, error)
	SetSeats(id model.Key, seats int) error
	ConvertToPaid(id model.Key, stripeID, subID, plan string, yearly bool, seats int) error
	ChangePlan(id model.Key, plan string, yearly bool) error
	Cancel(id model.Key) error
}

// AdminServices TODO: investigate this...
type AdminServices interface {
	SessionRefresher
	SessionCloser
	LogRequests(reqs []model.APIRequest) error
}

// WebhookServices is an interface that contains all functions to manage webhook.
type WebhookServices interface {
	SessionRefresher
	SessionCloser
	Add(accountID model.Key, events, url string) error
	List(accountID model.Key) ([]model.Webhook, error)
	Delete(accountID model.Key, event, url string) error
	AllSubscriptions(event string) ([]model.Webhook, error)
}
