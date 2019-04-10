package model

import (
	"time"
)

// Account represents the basic information for an account.
type Account struct {
	ID             int64     `json:"id"`
	Email          string    `json:"email"`
	StripeID       string    `json:"stripeId"`
	SubscriptionID string    `json:"subscriptionId"`
	Plan           string    ` json:"plan"`
	IsYearly       bool      `json:"isYearly"`
	SubscribedOn   time.Time `json:"subscribed"`
	Seats          int       ` json:"seats"`
	TrialInfo      Trial     ` json:"trial"`
	IsActive       bool      ` json:"active"`

	Users []User ` json:"users"`
}

// IsPaid returns if this account is a paying customer.
func (a *Account) IsPaid() bool {
	return len(a.StripeID) > 0 && len(a.SubscriptionID) > 0
}

// Trial represents the trial information for an account.
type Trial struct {
	IsTrial  bool      ` json:"trial"`
	Plan     string    ` json:"plan"`
	Start    time.Time ` json:"start"`
	Extended int       ` json:"extended"`
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
	ID           int64         `json:"id"`
	AccountID    int64         ` json:"accountId"`
	Email        string        `json:"email"`
	Password     string        ` json:"-"`
	Token        string        ` json:"token"`
	Role         Roles         ` json:"role"`
	AccessTokens []AccessToken ` json:"accessTokens"`
}

// AccessToken represents access tokens.
type AccessToken struct {
	ID     int64  ` json:"id"`
	UserID int64  ` json:"userId"`
	Name   string ` json:"name"`
	Token  string ` json:"token"`
}

// APIRequest represents a single API call.
type APIRequest struct {
	ID         int64     ` json:"id"`
	AccountID  int64     ` json:"accountId"`
	UserID     int64     ` json:"userId"`
	URL        string    `json:"url"`
	Requested  time.Time ` json:"requested"`
	StatusCode int       ` json:"statusCode"`
	RequestID  string    ` json:"reqId"`
}

// Webhook represents a webhook subscription.
type Webhook struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"accountId"`
	EventName string    `json:"event"`
	TargetURL string    `json:"url"`
	IsActive  bool      `json:"active"`
	Created   time.Time `json:"created"`
}
