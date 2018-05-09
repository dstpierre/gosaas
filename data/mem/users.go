package mem

import (
	"github.com/dstpierre/gosaas/data/model"
)

type Users struct {
	store []model.User
}

func (u *Users) GetDetail(id model.Key) (*model.User, error) {
	var user model.User
	for _, usr := range u.store {
		if usr.ID == id {
			user = usr
			break
		}
	}
	return &user, nil
}

func (u *Users) RefreshSession(conn *bool, dbName string) {
	u.store = append(u.store, model.User{ID: 1, Email: "test@domain.com"})
}
