package engine

import (
	"net/http"

	"github.com/dstpierre/gosaas/data/model"
)

// Route represent a web handler with optional middlewares
type Route struct {
	// middleware
	Logger bool

	// authorization
	MinimumRole model.Roles

	Handler http.Handler
}
