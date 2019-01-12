// +build mem

package controllers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dstpierre/gosaas"
	"github.com/dstpierre/gosaas/data"
)

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	db := &data.DB{}
	if err := db.Open("unit", "test"); err != nil {
		log.Fatal("error while creating mem data ", err)
	}

	routes := make(map[string]*gosaas.Route)
	routes["user"] = newUser()
	routes["billing"] = newBilling()

	mux := &gosaas.Server{
		DB:            db,
		Logger:        logger,
		Authenticator: authenticator,
		Routes:        routes,
	}

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	return rec
}

func Test_UserProfile_Handler(t *testing.T) {
	req, err := http.NewRequest("GET", "/user/profile", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := executeRequest(req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("returns status %v was expecting %v", status, http.StatusOK)
	}
}
