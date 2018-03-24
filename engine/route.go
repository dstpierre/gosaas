package engine

import "net/http"

// Route represent a web handler with optional middlewares
type Route struct {
	Logger  bool
	Handler http.Handler
}
