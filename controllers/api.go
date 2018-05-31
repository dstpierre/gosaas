package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dstpierre/gosaas/data"
	"github.com/dstpierre/gosaas/data/model"
	"github.com/dstpierre/gosaas/engine"
)

// API is the starting point of our API.
// Responsible for routing the request to the correct handler
type API struct {
	DB            *data.DB
	Logger        func(http.Handler) http.Handler
	Authenticator func(http.Handler) http.Handler
	User          *engine.Route
}

// NewAPI returns a production API with all middlewares
func NewAPI() *API {
	return &API{
		Logger:        engine.Logger,
		Authenticator: engine.Authenticator,
	}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, engine.ContextOriginalPath, r.URL.Path)

	if a.DB.CopySession {
		fmt.Println("copy mongo session")
		a.DB.Users.RefreshSession(a.DB.Connection, a.DB.DatabaseName)
	}

	ctx = context.WithValue(ctx, engine.ContextDatabase, a.DB)

	var next *engine.Route
	var head string
	head, r.URL.Path = engine.ShiftPath(r.URL.Path)
	if head == "user" {
		next = newUser()
	} else {
		next = newError(fmt.Errorf("path not found"), http.StatusNotFound)
	}

	ctx = context.WithValue(ctx, engine.ContextMinimumRole, next.MinimumRole)

	// make sure we are authenticating all calls
	next.Handler = a.Authenticator(next.Handler)

	if next.Logger {
		next.Handler = a.Logger(next.Handler)
	}

	next.Handler.ServeHTTP(w, r.WithContext(ctx))
}

func newError(err error, statusCode int) *engine.Route {
	return &engine.Route{
		Logger:      true,
		MinimumRole: model.RoleUser,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			engine.Respond(w, r, statusCode, err)
		}),
	}
}
