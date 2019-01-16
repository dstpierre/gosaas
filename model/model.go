package model

import (
	"time"
)

// Account represents the basic information for an account.
type Account struct {
	ID             Key       `bson:"_id" json:"id"`
	Email          string    `bson:"email" json:"email"`
	StripeID       string    `bson:"stripeId" json:"stripeId"`
	SubscriptionID string    `bson:"subId" json:"subscriptionId"`
	Plan           string    `bson:"plan" json:"plan"`
	IsYearly       bool      `bson:"isYearly" json:"isYearly"`
	SubscribedOn   time.Time `bson:"subscribed" json:"subscribed"`
	Seats          int       `bson:"seats" json:"seats"`
	TrialInfo      Trial     `bson:"trial" json:"trial"`

	Users []User `bson:"users" json:"users"`
}

// IsPaid returns if this account is a paying customer.
func (a *Account) IsPaid() bool {
	return len(a.StripeID) > 0 && len(a.SubscriptionID) > 0
}

// Trial represents the trial information for an account.
type Trial struct {
	IsTrial  bool      `bson:"trial" json:"trial"`
	Plan     string    `bson:"plan" json:"plan"`
	Start    time.Time `bson:"start" json:"start"`
	Extended int       `bson:"extended" json:"extended"`
}

// Roles are used with user access control and authorization. You may add custom roles in-between the
// default ones.
type Roles int

const (
	// RolePublic for publicly accessible routes.
	RolePublic Roles = 0
	// RoleFree for free user.
	RoleFree = 10
	// RoleUser for standard user.
	RoleUser = 20
	// RoleAdmin for admins.
	RoleAdmin = 99
)

// User represents a user.
type User struct {
	ID           Key           `bson:"_id" json:"id"`
	AccountID    Key           `bson:"accountId" json:"accountId"`
	Email        string        `bson:"email" json:"email"`
	Password     string        `bson:"pw" json:"-"`
	Token        string        `bson:"tok" json:"token"`
	Role         Roles         `bson:"role" json:"role"`
	AccessTokens []AccessToken `bson:"pat" json:"accessTokens"`
}

// AccessToken represents access tokens.
type AccessToken struct {
	ID     Key    `bson:"_id" json:"id"`
	UserID Key    `bson:"userId" json:"userId"`
	Name   string `bson:"name" json:"name"`
	Token  string `bson:"tok" json:"token"`
}

// APIRequest represents a single API call.
type APIRequest struct {
	ID         Key       `bson:"_id" json:"id"`
	AccountID  Key       `bson:"accountId" json:"accountId"`
	UserID     Key       `bson:"userId" json:"userId"`
	URL        string    `bson:"url" json:"url"`
	Requested  time.Time `bson:"reqon" json:"requested"`
	StatusCode int       `bson:"sc" json:"statusCode"`
	RequestID  string    `bson:"reqid" json:"reqId"`
}

// Webhook represents a webhook subscription.
type Webhook struct {
	ID        Key       `bson:"_id" json:"id"`
	AccountID Key       `bson:"accountId" json:"accountId"`
	EventName string    `bson:"event" json:"event"`
	TargetURL string    `bson:"url" json:"url"`
	IsActive  bool      `bson:"active" json:"active"`
	Created   time.Time `bson:"created" json:"created"`
}
