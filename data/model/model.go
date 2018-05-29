package model

type Account struct {
	ID    Key    `bson:"_id" json:"id"`
	Email string `bson:"email" json:"email"`

	Users []User `bson:"users" json:"users"`
}

type User struct {
	ID       Key    `bson:"_id" json:"id"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"pw" json:"-"`
}
