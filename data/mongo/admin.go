// +build !mem

package mongo

import (
	"github.com/dstpierre/gosaas/data/model"
	mgo "gopkg.in/mgo.v2"
)

type Admin struct {
	DB *mgo.Database
}

func (a *Admin) LogRequest(reqs []model.APIRequest) error {
	return a.DB.C("requests").Insert(reqs)
}

func (a *Admin) RefreshSession(s *mgo.Session, dbName string) {
	a.DB = s.Copy().DB(dbName)
}

func (a *Admin) Close() {
	a.DB.Session.Close()
}
