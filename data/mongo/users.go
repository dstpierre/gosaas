package mongo

import (
	"github.com/dstpierre/gosaas/data/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Users struct {
	DB *mgo.Database
}

func (u *Users) GetDetail(id model.Key) (*model.User, error) {
	var user model.User
	where := bson.M{"_id": id}
	if err := u.DB.C("users").Find(where).One(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *Users) RefreshSession(s *mgo.Session, dbName string) {
	u.DB = s.Copy().DB(dbName)
}
