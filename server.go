package gosaas

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dstpierre/gosaas/data"
	"github.com/dstpierre/gosaas/internal/config"
)

func init() {
	if err := config.LoadFromFile(); err != nil {
		log.Println(err)
	}

	if len(config.Current.StripeKey) > 0 {
		SetStripeKey(config.Current.StripeKey)
	}

	if len(config.Current.Plans) > 0 {
		for _, p := range config.Current.Plans {
			if p.Params == nil {
				p.Params = make(map[string]interface{})
			}
			data.AddPlan(p)
		}
	}
}

// Server is the starting point of the backend.
//
// Responsible for routing requests to handlers.
type Server struct {
	DB              *data.DB
	Logger          func(http.Handler) http.Handler
	Authenticator   func(http.Handler) http.Handler
	Throttler       func(http.Handler) http.Handler
	RateLimiter     func(http.Handler) http.Handler
	Cors            func(http.Handler) http.Handler
	StaticDirectory string
	Routes          map[string]*Route
}

// NewServer returns a production server with all available middlewares.
// Only the top level routes needs to be passed as parameter.
//
// There's three built-in implementations:
//
// 1. users: for user management (signup, signin, authentication, get detail, etc).
//
// 2. billing: for a fully functional billing process (converting from free to paid, changing plan, get invoices, etc).
//
// 3. webhooks: for allowing users to subscribe to events (you may trigger webhook via gosaas.SendWebhook).
//
// To override default inplementation you simply have to supply your own like so:
//
// 	routes := make(map[string]*gosaas.Route)
// 	routes["billing"] = &gosaas.Route{Handler: billing.Route}
//
// This would use your own billing implementation instead of the one supplied by gosaas.
func NewServer(routes map[string]*Route) *Server {
	// if users, billing and webhooks are not part
	// of the routes, we default to gosaas's implementation.
	if _, ok := routes["users"]; !ok {
		routes["users"] = newUser()
	}

	if _, ok := routes["billing"]; !ok {
		routes["billing"] = newBilling()
	}

	if _, ok := routes["webhooks"]; !ok {
		routes["webhooks"] = newWebhook()
	}

	return &Server{
		Logger:          Logger,
		Authenticator:   Authenticator,
		Throttler:       Throttler,
		RateLimiter:     RateLimiter,
		Cors:            Cors,
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
	defer r.Body.Close()

	if strings.HasPrefix(r.URL.Path, s.StaticDirectory) {
		http.ServeFile(w, r, r.URL.Path[1:])
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, ContextOriginalPath, r.URL.Path)

	isJSON := strings.ToLower(r.Header.Get("Content-Type")) == "application/json"
	ctx = context.WithValue(ctx, ContextContentIsJSON, isJSON)

	var route *Route
	var handler http.Handler
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	if r, ok := s.Routes[head]; ok {
		route = r
		handler = r.Handler
	} else if catchall, ok := s.Routes["__catchall__"]; ok {
		route = catchall
		handler = catchall.Handler
	} else {
		route = NewError(fmt.Errorf("path not found"), http.StatusNotFound)
		handler = route.Handler
	}

	if route.WithDB {
		ctx = context.WithValue(ctx, ContextDatabase, s.DB)
	}

	ctx = context.WithValue(ctx, ContextMinimumRole, route.MinimumRole)

	// make sure we are authenticating all calls
	handler = s.Authenticator(handler)

	if route.Logger {
		handler = s.Logger(handler)
	}

	if route.EnforceRateLimit {
		handler = s.RateLimiter(handler)
		handler = s.Throttler(handler)
	}

	// are we allowing cross-origin requests for this route
	if route.AllowCrossOrigin {
		handler = s.Cors(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}
