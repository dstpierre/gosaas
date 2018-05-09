// +build !mgo
// +build !mem

package data

import (
	"github.com/dstpierre/gosaas/data/model"
	"github.com/dstpierre/gosaas/data/pg"
)

func (db *DB) Open(driverName, dataSourceName string) error {
	conn, err := model.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}

	db.Users = &pg.Users{DB: conn}

	db.Connection = conn
	return nil
}
