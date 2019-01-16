package gosaas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dstpierre/gosaas/data"
	"github.com/dstpierre/gosaas/model"
)

func post(url string, data interface{}, result interface{}, headers map[string]string) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return err
	}

	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "applicaiton/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("failed post, error: ", err)
		} else {
			log.Println("failed post, recv: ", string(b))
		}
		return fmt.Errorf("error requesting %s returned %s", url, resp.Status)
	}

	if result != nil {
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(result); err != nil {
			return err
		}
	}
	return nil
}

// SendWebhook posts data to all subscribers of an event.
func SendWebhook(wh data.WebhookServices, event string, data interface{}) {
	defer func() {
		fmt.Println("closing webhooks connection")
		wh.Close()
	}()

	headers := make(map[string]string)
	headers["X-Webhook-Event"] = event

	subscribers, err := wh.AllSubscriptions(event)
	if err != nil {
		log.Println("unable to get webhook subscribers for ", event)
		return
	}

	for _, sub := range subscribers {
		go func(sub model.Webhook, headers map[string]string) {
			if err := post(sub.TargetURL, data, nil, headers); err != nil {
				log.Println("error calling URL", sub.TargetURL, err)
			}
		}(sub, headers)
	}
}

// Webhook handles everything related to the /webhooks requests
//
// POST /webhooks -> subscribe to events
// GET /webhooks -> get the list of subscriptions for current user
// POST /webhooks/unsub -> remove a subscription
type Webhook struct{}

func newWebhook() *Route {
	var wh interface{} = Webhook{}
	return &Route{
		Logger:      true,
		MinimumRole: model.RoleFree,
		Handler:     wh.(http.Handler),
	}
}

func (wh Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	if head == "" {
		if r.Method == "POST" {
			wh.subscribe(w, r)
		} else if r.Method == "GET" {
			wh.list(w, r)
		}
	} else if head == "unsub" && r.Method == "POST" {
		wh.delete(w, r)
	}
}

type addSubscriber struct {
	Events string `json:"events"`
	URL    string `json:"url"`
}

func (wh *Webhook) subscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(ContextAuth).(Auth)
	db := ctx.Value(ContextDatabase).(*data.DB)

	var data addSubscriber
	if err := ParseBody(r.Body, &data); err != nil {
		Respond(w, r, http.StatusBadRequest, err)
		return
	}

	if err := db.Webhooks.Add(keys.AccountID, data.Events, data.URL); err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusCreated, true)
}

func (wh *Webhook) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(ContextAuth).(Auth)
	db := ctx.Value(ContextDatabase).(*data.DB)

	subs, err := db.Webhooks.List(keys.AccountID)
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, subs)
}

type DeleteSubscription struct {
	Event string `json:"event"`
	URL   string `json:"url"`
}

func (wh *Webhook) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(ContextAuth).(Auth)
	db := ctx.Value(ContextDatabase).(*data.DB)

	var data DeleteSubscription
	if err := ParseBody(r.Body, &data); err != nil {
		Respond(w, r, http.StatusBadRequest, err)
		return
	}

	if err := db.Webhooks.Delete(keys.AccountID, data.Event, data.URL); err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, true)
}
