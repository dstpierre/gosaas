package model

type User struct {
	ID    Key    `bson:"_id" json:"id"`
	Email string `bson:"email" json:"email"`
}
