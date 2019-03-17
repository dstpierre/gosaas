package postgres

import (
	"database/sql"

	"github.com/dstpierre/gosaas/model"
)

type Webhooks struct {
	DB *sql.DB
}

func (wh *Webhooks) Add(accountID int64, events, url string) error {
	_, err := wh.DB.Exec(`
		insert into webhooks(account_id, events, url)
		VALUES($1, $2, $3)
	`, accountID, events, url)
	return err
}

func (wh *Webhooks) List(accountID int64) ([]model.Webhook, error) {
	rows, err := wh.DB.Query("SELECT * FROM webhooks WHERE account_id = $1", accountID)
	if err != nil {
		return nil, err
	}

	var hooks []model.Webhook
	for rows.Next() {
		var hook model.Webhook
		if err := wh.scan(rows, &hook); err != nil {
			return nil, err
		}

		hooks = append(hooks, hook)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return hooks, nil
}

func (wh *Webhooks) Delete(accountID int64, event, url string) error {
	_, err := wh.DB.Exec(`
		DELETE FROM webhooks
		WHERE account_id = $1 AND
					events = $2 AND
					url = $3
	`, accountID, event, url)
	return err
}

func (wh *Webhooks) AllSubscriptions(event string) ([]model.Webhook, error) {
	rows, err := wh.DB.Query("SELECT * FROM webhooks WHERE events = $1")
	if err != nil {
		return nil, err
	}

	var hooks []model.Webhook
	for rows.Next() {
		var hook model.Webhook
		if err := wh.scan(rows, &hook); err != nil {
			return nil, err
		}

		hooks = append(hooks, hook)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return hooks, nil
}

func (wh *Webhooks) scan(rows *sql.Rows, hook *model.Webhook) error {
	return rows.Scan(hook.ID,
		hook.AccountID,
		hook.EventName,
		hook.TargetURL,
		hook.IsActive,
		hook.Created,
	)
}
