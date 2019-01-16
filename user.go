package gosaas

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/dstpierre/gosaas/data"
	"github.com/dstpierre/gosaas/model"
	"golang.org/x/crypto/bcrypt"
)

// User handles everything related to the /user requests
type User struct{}

func newUser() *Route {
	var u interface{} = User{}
	return &Route{
		Logger:      true,
		MinimumRole: model.RolePublic,
		WithDB:      true,
		Handler:     u.(http.Handler),
	}
}

func (u User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	if head == "signup" && r.Method == http.MethodPost {
		u.signup(w, r)
	} else if head == "signin" && r.Method == http.MethodPost {
		u.signin(w, r)
	} else if head == "profile" && r.Method == http.MethodGet {
		u.profile(w, r)
	} else {
		notFound(w)
	}
}

func (u User) signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value(ContextDatabase).(*data.DB)

	var data = new(struct {
		Email string `json:"email"`
	})

	if err := ParseBody(r.Body, &data); err != nil {
		Respond(w, r, http.StatusBadRequest, err)
		return
	}

	pw := randStringRunes(7)
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	acct, err := db.Users.SignUp(data.Email, string(b))
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusCreated, acct)
}

func (u User) signin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value(ContextDatabase).(*data.DB)

	var data = new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})

	if err := ParseBody(r.Body, &data); err != nil {
		Respond(w, r, http.StatusBadRequest, err)
		return
	}

	user, err := db.Users.GetUserByEmail(data.Email)
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		Respond(w, r, http.StatusNotFound, fmt.Errorf("invalid email/password."))
		return
	}

	Respond(w, r, http.StatusOK, user)
}

func (u User) profile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	keys := ctx.Value(ContextAuth).(Auth)
	db := ctx.Value(ContextDatabase).(*data.DB)

	acct, err := db.Users.GetDetail(keys.AccountID)
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	Respond(w, r, http.StatusOK, acct)
}

func randStringRunes(n int) string {
	letterRunes := []rune("abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ2345679")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
