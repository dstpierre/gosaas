package gosaas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dstpierre/gosaas/data"
	"github.com/dstpierre/gosaas/engine"
)

// Server is the starting point of the backend.
// Responsible for routing requests to handlers.
type Server struct {
	DB            *data.DB
	Logger        func(http.Handler) http.Handler
	Authenticator func(http.Handler) http.Handler
	Throttler     func(http.Handler) http.Handler
	RateLimiter   func(http.Handler) http.Handler
	Routes        map[string]*engine.Route
}

// NewAPI returns a production API with all middlewares
func NewServer(routes map[string]*engine.Route) *Server {
	return &Server{
		Logger:        engine.Logger,
		Authenticator: engine.Authenticator,
		Throttler:     engine.Throttler,
		RateLimiter:   engine.RateLimiter,
		Routes:        routes,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, engine.ContextOriginalPath, r.URL.Path)

	if s.DB.CopySession {
		s.DB.Users.RefreshSession(s.DB.Connection, s.DB.DatabaseName)
		s.DB.Webhooks.RefreshSession(s.DB.Connection, s.DB.DatabaseName)

		defer func() {
			s.DB.Users.Close()
			s.DB.Webhooks.Close()
		}()
	}

	ctx = context.WithValue(ctx, engine.ContextDatabase, s.DB)

	var next *engine.Route
	var head string
	head, r.URL.Path = engine.ShiftPath(r.URL.Path)
	if r, ok := s.Routes[head]; ok {
		next = r
	} else {
		next = engine.NewError(fmt.Errorf("path not found"), http.StatusNotFound)
	}

	ctx = context.WithValue(ctx, engine.ContextMinimumRole, next.MinimumRole)

	// make sure we are authenticating all calls
	next.Handler = s.Authenticator(next.Handler)

	if next.Logger {
		next.Handler = s.Logger(next.Handler)
	}

	next.Handler = s.RateLimiter(next.Handler)
	next.Handler = s.Throttler(next.Handler)

	next.Handler.ServeHTTP(w, r.WithContext(ctx))
}
