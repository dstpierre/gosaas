// +build !mem

package mongo

import (
	"github.com/dstpierre/gosaas/data/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Users struct {
	DB *mgo.Database
}

func (u *Users) SignUp(email, password string) (*model.Account, error) {
	accountID := bson.NewObjectId()

	acct := model.Account{ID: accountID, Email: email}
	acct.Users = append(acct.Users, model.User{
		ID:       bson.NewObjectId(),
		Email:    email,
		Password: password,
		Token:    model.NewToken(accountID),
	})
	if err := u.DB.C("users").Insert(acct); err != nil {
		return nil, err
	}
	return &acct, nil
}

func (u *Users) AddToken(accountID, userID model.Key, name string) (*model.AccessToken, error) {
	tok := model.AccessToken{
		ID:    bson.NewObjectId(),
		Name:  name,
		Token: model.NewToken(accountID),
	}

	where := bson.M{"_id": accountID, "users._id": userID}
	update := bson.M{"$push": bson.M{"users.$.pat": bson.M{"$push": tok}}}
	if err := u.DB.C("users").Update(where, update); err != nil {
		return nil, err
	}
	return &tok, nil
}

func (u *Users) RemoveToken(accountID, userID, tokenID model.Key) error {
	where := bson.M{"_id": accountID, "users._id": userID}
	update := bson.M{"$pull": bson.M{"users.$.pat": bson.M{"_id": tokenID}}}
	return u.DB.C("users").Update(where, update)
}

func (u *Users) GetDetail(id model.Key) (*model.Account, error) {
	var acct model.Account
	where := bson.M{"_id": id}
	if err := u.DB.C("users").Find(where).One(&acct); err != nil {
		return nil, err
	}
	return &acct, nil
}

func (u *Users) RefreshSession(s *mgo.Session, dbName string) {
	u.DB = s.Copy().DB(dbName)
}
