package engine

import (
	"net/http"

	"github.com/dstpierre/gosaas/data/model"
)

// Route represent a web handler with optional middlewares
type Route struct {
	// middleware
	WithDB           bool
	Logger           bool
	EnforceRateLimit bool

	// authorization
	MinimumRole model.Roles

	Handler http.Handler
}

func NewError(err error, statusCode int) *Route {
	return &Route{
		Logger:      true,
		MinimumRole: model.RolePublic,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Respond(w, r, statusCode, err)
		}),
	}
}
