// +build mem

package gosaas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/dstpierre/gosaas/data"
	"github.com/dstpierre/gosaas/model"
)

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		dr, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println("unable to dump request", err)
		} else {
			ctx = context.WithValue(ctx, ContextRequestDump, dr)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		key, _, err := extractKeyFromRequest(r)
		if err != nil {
			// we simulate the way the real authenticator middleware
			// act by simply doing nothing for public route
			mr := ctx.Value(ContextMinimumRole).(model.Roles)
			if mr == model.RolePublic {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			http.Error(w, fmt.Sprintf("error returned when extracting key %v", err), http.StatusInternalServerError)
			return
		}

		if key == "unit-test-token" {
			a := Auth{
				AccountID: 1,
				Email:     "unit@test.com",
				Role:      model.RoleAdmin,
				UserID:    1,
			}
			ctx = context.WithValue(ctx, ContextAuth, a)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func executeRequest(req *http.Request) (*httptest.ResponseRecorder, *Server) {
	req.Header.Set("Content-Type", "application/json")

	db := &data.DB{}
	if err := db.Open("unit", "test"); err != nil {
		log.Fatal("error while creating mem data ", err)
	}

	routes := make(map[string]*Route)
	// we use the built-in routes for tests
	routes["users"] = newUser()
	routes["billing"] = newBilling()
	routes["webhooks"] = newWebhook()

	mux := &Server{
		DB:              db,
		Logger:          logger,
		Authenticator:   authenticator,
		StaticDirectory: "/public/",
		Routes:          routes,
	}

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	if rec.Code >= http.StatusBadRequest {
		fmt.Printf("error while requesting %s\n%s\n\n", req.URL.Path, string(rec.Body.Bytes()))
	}
	return rec, mux
}

func Test_Users_SignUp(t *testing.T) {
	t.Parallel()

	var data = new(struct {
		Email string `json:"email"`
	})
	data.Email = "new@user.com"

	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/users/signup", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	rec, svr := executeRequest(req)
	if status := rec.Code; status != http.StatusCreated {
		t.Errorf("returns status %v was expecting %v", status, http.StatusCreated)
	}

	var user model.Account
	if err := ParseBody(ioutil.NopCloser(bytes.NewReader(rec.Body.Bytes())), &user); err != nil {
		t.Fatal(err)
	} else if user.Email != data.Email {
		t.Errorf("returns user's email %v differ from added %v", user.Email, data.Email)
	}

	// we validate that the new user has been saved
	acct, err := svr.DB.Users.GetDetail(user.ID)
	if err != nil {
		t.Fatal(err)
	} else if acct.Email != data.Email {
		t.Errorf("database email is %s and was expecting %s", acct.Email, data.Email)
	}
}

func Test_Users_SignIn(t *testing.T) {
	t.Parallel()

	var data = new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	data.Email = "test@domain.com"
	data.Password = "unit-test"

	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/users/login", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	rec, _ := executeRequest(req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("returns status %v was expecting %v", status, http.StatusOK)
	}

	var user model.User
	if err := ParseBody(ioutil.NopCloser(bytes.NewReader(rec.Body.Bytes())), &user); err != nil {
		t.Errorf("error while parsing returning JSON: %v", err)
	} else if user.Email != "test@domain.com" {
		t.Errorf("email was %s was expecting test@domain.com", user.Email)
	}
}

func Test_Users_Profile(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest("GET", "/users/profile?key=unit-test-token", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec, _ := executeRequest(req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("returns status %v was expecting %v", status, http.StatusOK)
	}
}
