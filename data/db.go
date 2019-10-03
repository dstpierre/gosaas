package data

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/dstpierre/gosaas/model"
)

// DB is a database agnostic abstraction that contains a reference to
// the database connection.
//
// At this moment Postgres and an in-memory data provider are supported.
type DB struct {
	// DatabaseName is the name of the database used.
	DatabaseName string
	// Connection is the reference to the database connection.
	Connection *sql.DB

	// Users contains the data access functions related to account, user and billing.
	Users UserServices
	// Webhooks contains the data access functions related to managing Webhooks.
	Webhooks WebhookServices
}

// UserServices is an interface that contians all functions related to account, user and billing.
type UserServices interface {
	SignUp(email, password string) (*model.Account, error)
	ChangePassword(id, accountID int64, passwd string) error
	AddToken(accountID, userID int64, name string) (*model.AccessToken, error)
	RemoveToken(accountID, userID, tokenID int64) error
	Auth(accountID int64, token string, pat bool) (*model.Account, *model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetDetail(id int64) (*model.Account, error)
	GetByStripe(stripeID string) (*model.Account, error)
	SetSeats(id int64, seats int) error
	ConvertToPaid(id int64, stripeID, subID, plan string, yearly bool, seats int) error
	ChangePlan(id int64, plan string, yearly bool) error
	Cancel(id int64) error
}

// AdminServices TODO: investigate this...
type AdminServices interface {
	LogRequests(reqs []model.APIRequest) error
}

// WebhookServices is an interface that contains all functions to manage webhook.
type WebhookServices interface {
	Add(accountID int64, events, url string) error
	List(accountID int64) ([]model.Webhook, error)
	Delete(accountID int64, event, url string) error
	AllSubscriptions(event string) ([]model.Webhook, error)
}

// NewID returns a per second unique string based on account and user ids.
func NewID(accountID, userID int64) string {
	n := time.Now()
	i, _ := strconv.Atoi(
		fmt.Sprintf("%d%d%d%d%d%d%d%d",
			accountID,
			userID,
			n.Year()-2000,
			int(n.Month()),
			n.Day(),
			n.Hour(),
			n.Minute(),
			n.Second()))
	return fmt.Sprintf("%x", i)
}
