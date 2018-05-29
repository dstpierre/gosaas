package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dstpierre/gosaas/data"
	"github.com/dstpierre/gosaas/data/model"
	"github.com/dstpierre/gosaas/engine"
)

// User handles everything related to the /user requests
type User struct{}

func newUser() *engine.Route {
	var u interface{} = User{}
	return &engine.Route{
		Logger:  true,
		Handler: u.(http.Handler),
	}
}

func (u User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = engine.ShiftPath(r.URL.Path)
	if head == "profile" {
		u.profile(w, r)
		return
	} else if head == "detail" {
		head, _ := engine.ShiftPath(r.URL.Path)
		i, err := strconv.ParseInt(head, 10, 64)
		if err != nil {
			newError(err, http.StatusInternalServerError).Handler.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, engine.ContextUserID, i)
		u.detail(w, r.WithContext(ctx))
		return
	}
	newError(fmt.Errorf("path not found"), http.StatusNotFound).Handler.ServeHTTP(w, r)
}

func (u User) profile(w http.ResponseWriter, r *http.Request) {
	engine.Respond(w, r, http.StatusOK, "viewing detail")
}

func (u User) detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := ctx.Value(engine.ContextUserID).(model.Key)
	db := ctx.Value(engine.ContextDatabase).(*data.DB)

	var result = new(struct {
		ID    model.Key `json:"userId"`
		Email string    `json:"email"`
		Time  time.Time `json:"time"`
	})

	user, err := db.Users.GetDetail(id)
	if err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	result.ID = user.ID
	result.Email = user.Email
	result.Time = time.Now()

	engine.Respond(w, r, http.StatusOK, result)
}
