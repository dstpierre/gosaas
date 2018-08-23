package model

import (
	"time"
)

type Account struct {
	ID    Key    `bson:"_id" json:"id"`
	Email string `bson:"email" json:"email"`

	Users []User `bson:"users" json:"users"`
}

type Roles int

const (
	RoleAdmin Roles = iota
	RoleUser
)

type User struct {
	ID           Key           `bson:"_id" json:"id"`
	Email        string        `bson:"email" json:"email"`
	Password     string        `bson:"pw" json:"-"`
	Token        string        `bson:"tok" json:"token"`
	Role         Roles         `bson:"role" json:"role"`
	AccessTokens []AccessToken `bson:"pat" json:"accessTokens"`
}

type AccessToken struct {
	ID    Key    `bson:"_id" json:"id"`
	Name  string `bson:"name" json:"name"`
	Token string `bson:"tok" json:"token"`
}

// APIRequest represents a single API call
type APIRequest struct {
	ID         Key       `bson:"_id" json:"id"`
	AccountID  Key       `bson:"accountId" json:"accountId"`
	UserID     Key       `bson:"userId" json:"userId"`
	URL        string    `bson:"url" json:"url"`
	Requested  time.Time `bson:"reqon" json:"requested"`
	StatusCode int       `bson:"sc" json:"statusCode"`
	RequestID  string    `bson:"reqid" json:"reqId"`
}

type Webhook struct {
	ID        Key       `bson:"_id" json:"id"`
	AccountID Key       `bson:"accountId" json:"accountId"`
	EventName string    `bson:"event" json:"event"`
	TargetURL string    `bson:"url" json:"url"`
	IsActive  bool      `bson:"active" json:"active"`
	Created   time.Time `bson:"created" json:"created"`
}
