package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
		i, err := strconv.Atoi(head)
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
	id := ctx.Value(engine.ContextUserID)

	var result = new(struct {
		ID   int       `json:"userId"`
		Time time.Time `json:"time"`
	})

	result.ID = id.(int)
	result.Time = time.Now()

	engine.Respond(w, r, http.StatusOK, result)
}
