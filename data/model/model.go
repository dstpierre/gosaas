package model

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
