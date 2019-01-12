package gosaas

import (
	"net/http"

	"github.com/dstpierre/gosaas/model"
)

// Route represents a web handler with optional middlewares.
type Route struct {
	// middleware
	WithDB           bool
	Logger           bool
	EnforceRateLimit bool

	// authorization
	MinimumRole model.Roles

	Handler http.Handler
}

// NewError returns a new route that simply Respond with the error and status code.
func NewError(err error, statusCode int) *Route {
	return &Route{
		Logger:      true,
		MinimumRole: model.RolePublic,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Respond(w, r, statusCode, err)
		}),
	}
}
