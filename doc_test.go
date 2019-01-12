package gosaas

import (
	"log"
	"net/http"

	"github.com/dstpierre/gosaas/model"
)

func ExampleNewServer() {
	routes := make(map[string]*Route)
	routes["task"] = &Route{
		Logger:      true,           // enable logging
		MinimumRole: model.RoleFree, // make sure only free user and up can access this route
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// normally you would use a custom type that implement ServerHTTP
			// and handle all sub-level route for this top-level route.
			Respond(w, r, http.StatusOK, "list of tasks....")
		}),
	}

	mux := NewServer(routes)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
