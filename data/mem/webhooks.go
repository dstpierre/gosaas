// +build mem

package mem

import (
	"strings"
	"time"

	"github.com/dstpierre/gosaas/data/model"
)

type Webhooks struct {
	webhooks []model.Webhook
}

func (wh *Webhooks) Add(accountID model.Key, events, url string) error {
	id := len(wh.webhooks)
	acctID := id * 250

	en := strings.Split(events, ",")
	for _, e := range en {
		wh.webhooks = append(wh.webhooks, model.Webhook{
			ID:        id,
			AccountID: acctID,
			EventName: strings.Trim(e, " "),
			TargetURL: url,
			IsActive:  true,
			Created:   time.Now(),
		})

		id++
	}

	return nil
}

func (wh *Webhooks) List(accountID model.Key) ([]model.Webhook, error) {
	var filtered []model.Webhook
	for _, w := range wh.webhooks {
		if w.AccountID == accountID {
			filtered = append(filtered, w)
		}
	}

	return filtered, nil
}

func (wh *Webhooks) Delete(accountID model.Key, event, url string) error {
	// no need to implement
	return nil
}

func (wh *Webhooks) AllSubscriptions(event string) ([]model.Webhook, error) {
	var filtered []model.Webhook
	for _, w := range wh.webhooks {
		if w.EventName == event {
			filtered = append(filtered, w)
		}
	}
	return filtered, nil
}

func (wh *Webhooks) RefreshSession(conn *bool, dbName string) {
}

func (wh *Webhooks) Close() {
}
