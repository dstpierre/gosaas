package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dstpierre/gosaas/data"
	"github.com/dstpierre/gosaas/data/model"
	"github.com/dstpierre/gosaas/engine"
)

// Webhook handles everything related to the /webhooks requests
type Webhook struct{}

func sendWebhook(wh data.WebhookServices, event string, data interface{}) {
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
			if err := engine.Post(sub.TargetURL, data, nil, headers); err != nil {
				log.Println("error calling URL", sub.TargetURL, err)
			}
		}(sub, headers)
	}
}

func newWebhook() *engine.Route {
	var wh interface{} = Webhook{}
	return &engine.Route{
		Logger:      true,
		MinimumRole: model.RoleAdmin,
		Handler:     wh.(http.Handler),
	}
}

func (wh Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = engine.ShiftPath(r.URL.Path)
	if head == "subscribe" && r.Method == "POST" {
		wh.subscribe(w, r)
	} else if head == "list" && r.Method == "GET" {
		wh.list(w, r)
	} else if head == "unsub" && r.Method == "POST" {
		wh.delete(w, r)
	}
}

type AddSubscriber struct {
	Events string `json:"events"`
	URL    string `json:"url"`
}

func (wh *Webhook) subscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(engine.ContextAuth).(engine.Auth)
	db := ctx.Value(engine.ContextDatabase).(*data.DB)

	var data AddSubscriber
	if err := engine.ParseBody(r.Body, &data); err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	if err := db.Webhooks.Add(keys.AccountID, data.Events, data.URL); err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}
	engine.Respond(w, r, http.StatusCreated, true)
}

func (wh *Webhook) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(engine.ContextAuth).(engine.Auth)
	db := ctx.Value(engine.ContextDatabase).(*data.DB)

	subs, err := db.Webhooks.List(keys.AccountID)
	if err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}
	engine.Respond(w, r, http.StatusOK, subs)
}

type DeleteSubscription struct {
	Event string `json:"event"`
	URL   string `json:"url"`
}

func (wh *Webhook) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(engine.ContextAuth).(engine.Auth)
	db := ctx.Value(engine.ContextDatabase).(*data.DB)

	var data DeleteSubscription
	if err := engine.ParseBody(r.Body, &data); err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	if err := db.Webhooks.Delete(keys.AccountID, data.Event, data.URL); err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}
	engine.Respond(w, r, http.StatusOK, true)
}
