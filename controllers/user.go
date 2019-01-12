package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dstpierre/gosaas"
	"github.com/dstpierre/gosaas/data"
	"github.com/dstpierre/gosaas/model"
)

// User handles everything related to the /user requests
type User struct{}

func newUser() *gosaas.Route {
	var u interface{} = User{}
	return &gosaas.Route{
		Logger:      true,
		MinimumRole: model.RoleAdmin,
		Handler:     u.(http.Handler),
	}
}

func (u User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = gosaas.ShiftPath(r.URL.Path)
	if head == "profile" {
		u.profile(w, r)
		return
	} else if head == "detail" {
		u.detail(w, r)
		return
	}
	gosaas.NewError(fmt.Errorf("path not found"), http.StatusNotFound).Handler.ServeHTTP(w, r)
}

func (u User) profile(w http.ResponseWriter, r *http.Request) {
	gosaas.Respond(w, r, http.StatusOK, "viewing detail")
}

func (u User) detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(gosaas.ContextAuth).(gosaas.Auth)
	db := ctx.Value(gosaas.ContextDatabase).(*data.DB)

	var result = new(struct {
		ID    model.Key `json:"userId"`
		Email string    `json:"email"`
		Time  time.Time `json:"time"`
	})

	user, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		gosaas.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	/*
		var wh data.WebhookServices
		switch v := db.Webhooks.(type) {
		case *model.Webhooks:
			wh = &model.Webhooks{}
		default:
			log.Println("unhandled data type")
			wh = v
		}
		wh.RefreshSession(db.Connection, db.DatabaseName)
		go sendWebhook(wh, gosaas.WebhookEventUserDetail, user)
	*/

	result.ID = user.ID
	result.Email = user.Email
	result.Time = time.Now()

	gosaas.Respond(w, r, http.StatusOK, result)
}
