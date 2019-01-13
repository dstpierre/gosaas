package gosaas

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dstpierre/gosaas/data"
)

// Server is the starting point of the backend.
//
// Responsible for routing requests to handlers.
type Server struct {
	DB              *data.DB
	Logger          func(http.Handler) http.Handler
	Authenticator   func(http.Handler) http.Handler
	Throttler       func(http.Handler) http.Handler
	RateLimiter     func(http.Handler) http.Handler
	StaticDirectory string
	Routes          map[string]*Route
}

// NewServer returns a production server with all available middlewares.
// Only the top level routes needs to be passed as parameter.
func NewServer(routes map[string]*Route) *Server {
	return &Server{
		Logger:          Logger,
		Authenticator:   Authenticator,
		Throttler:       Throttler,
		RateLimiter:     RateLimiter,
		StaticDirectory: "/public/",
		Routes:          routes,
	}
}

// ServeHTTP is where the top level routes get matched with the map[string]*gosaas.Route
// received from the call to NewServer. Middleware are applied based on the found route properties.
//
// If no route can be found an error is returned.
//
// Static files are served from the "/public/" directory by default. To change this
// you may set the StaticDirectory after creating the server like this:
//
// 	mux := gosaas.NewServer(routes)
// 	mux.StaticDirectory = "/files/"
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, s.StaticDirectory) {
		http.ServeFile(w, r, r.URL.Path[1:])
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, ContextOriginalPath, r.URL.Path)

	if s.DB.CopySession {
		s.DB.Users.RefreshSession(s.DB.Connection, s.DB.DatabaseName)
		s.DB.Webhooks.RefreshSession(s.DB.Connection, s.DB.DatabaseName)

		defer func() {
			s.DB.Users.Close()
			s.DB.Webhooks.Close()
		}()
	}

	ctx = context.WithValue(ctx, ContextDatabase, s.DB)

	var next *Route
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	if r, ok := s.Routes[head]; ok {
		next = r
	} else {
		next = NewError(fmt.Errorf("path not found"), http.StatusNotFound)
	}

	ctx = context.WithValue(ctx, ContextMinimumRole, next.MinimumRole)

	// make sure we are authenticating all calls
	next.Handler = s.Authenticator(next.Handler)

	if next.Logger {
		next.Handler = s.Logger(next.Handler)
	}

	if next.EnforceRateLimit {
		next.Handler = s.RateLimiter(next.Handler)
		next.Handler = s.Throttler(next.Handler)
	}

	next.Handler.ServeHTTP(w, r.WithContext(ctx))
}
