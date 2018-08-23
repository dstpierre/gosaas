// +build !mem

package mongo

import (
	"strings"
	"time"

	"github.com/dstpierre/gosaas/data/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Webhooks struct {
	DB *mgo.Database
}

// Add inserts a new webhook
func (wh *Webhooks) Add(accountID model.Key, events, url string) error {
	var hooks []model.Webhook
	en := strings.Split(events, ",")
	for _, e := range en {
		hooks = append(hooks, model.Webhook{
			ID:        bson.NewObjectId(),
			AccountID: accountID,
			EventName: strings.Trim(e, " "),
			TargetURL: url,
			IsActive:  true,
			Created:   time.Now(),
		})
	}
	return wh.DB.C("webhooks").Insert(hooks)
}

// Delete removes all matching target webhook url
func (wh *Webhooks) Delete(accountID model.Key, event, url string) error {
	where := bson.M{"accountId": accountID, "event": event, "url": url}
	if _, err := wh.DB.C("webhooks").RemoveAll(where); err != nil {
		return err
	}
	return nil
}

// List returns the webhook entries for an account
func (wh *Webhooks) List(accountID model.Key) ([]model.Webhook, error) {
	var results []model.Webhook
	where := bson.M{"accountId": accountID}
	if err := wh.DB.C("webhooks").Find(where).All(&results); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return results, nil
}

// AllSubscriptions returns all webhooks for an event
func (wh *Webhooks) AllSubscriptions(event string) ([]model.Webhook, error) {
	var results []model.Webhook
	where := bson.M{"event": event}
	if err := wh.DB.C("webhooks").Find(where).All(&results); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return results, nil
}

func (wh *Webhooks) RefreshSession(s *mgo.Session, dbName string) {
	wh.DB = s.Copy().DB(dbName)
}

func (wh *Webhooks) Close() {
	wh.DB.Session.Close()
}
