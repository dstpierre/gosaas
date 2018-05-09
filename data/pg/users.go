package pg

import (
	"database/sql"

	"github.com/dstpierre/gosaas/data/model"
)

type Users struct {
	DB *sql.DB
}

func (u *Users) GetDetail(id model.Key) (*model.User, error) {
	var user model.User
	if err := u.DB.QueryRow("SELECT * FROM users WHERE user_Id = $1", id).Scan(&user.ID, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *Users) RefreshSession(conn *sql.DB, dbName string) {
}
