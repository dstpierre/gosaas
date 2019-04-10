package gosaas

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/dstpierre/gosaas/data"
	"github.com/dstpierre/gosaas/internal/config"
	"github.com/dstpierre/gosaas/model"
	"golang.org/x/crypto/bcrypt"
)

// User handles everything related to the /user requests
type User struct{}

func newUser() *Route {
	var u interface{} = User{}
	return &Route{
		AllowCrossOrigin: true,
		Logger:           true,
		MinimumRole:      model.RolePublic,
		WithDB:           true,
		Handler:          u.(http.Handler),
	}
}

func (u User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	if head == "signup" {
		if r.Method == http.MethodGet {
			u.signup(w, r)
		} else if r.Method == http.MethodPost {
			u.create(w, r)
		}
	} else if head == "login" {
		if r.Method == http.MethodGet {
			u.login(w, r)
		} else if r.Method == http.MethodPost {
			u.signin(w, r)
		}
	} else if head == "profile" && r.Method == http.MethodGet {
		u.profile(w, r)
	} else {
		notFound(w)
	}
}

func (u User) signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ServePage(w, r, config.Current.SignUpTemplate, CreateViewData(ctx, nil, nil))
}

func (u User) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value(ContextDatabase).(*data.DB)
	isJSON := ctx.Value(ContextContentIsJSON).(bool)

	var data = new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})

	if isJSON {
		if err := ParseBody(r.Body, &data); err != nil {
			Respond(w, r, http.StatusBadRequest, err)
			return
		}
	} else {
		r.ParseForm()

		data.Email = r.Form.Get("email")
	}

	if len(data.Password) == 0 {
		data.Password = randStringRunes(7)
		fmt.Println("TODO: remove this, temporary password", data.Password)
	}

	b, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		if isJSON {
			Respond(w, r, http.StatusInternalServerError, err)
		} else {
			http.Redirect(w, r, config.Current.SignUpErrorRedirect, http.StatusSeeOther)
		}
		return
	}

	acct, err := db.Users.SignUp(data.Email, string(b))
	if err != nil {
		if isJSON {
			Respond(w, r, http.StatusInternalServerError, err)
		} else {
			http.Redirect(w, r, config.Current.SignUpErrorRedirect, http.StatusSeeOther)
		}
		return
	}

	if config.Current.SignUpSendEmailValidation {
		u.sendEmail(acct.Email, data.Password)
	}

	if isJSON {
		Respond(w, r, http.StatusCreated, acct)
	} else {
		ck := &http.Cookie{
			Name:  "X-API-KEY",
			Path:  "/",
			Value: acct.Users[0].Token,
		}

		// we set the cookie so they will have their authentication token on the next request.
		http.SetCookie(w, ck)

		http.Redirect(w, r, config.Current.SignUpSuccessRedirect, http.StatusSeeOther)
	}
}

func (u User) sendEmail(email, pass string) {
	//TODO: implement this
	//queue.Enqueue(queue.TaskEmail, queue.SendEmailParameter{})
}

func (u User) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ServePage(w, r, config.Current.SignInTemplate, CreateViewData(ctx, nil, nil))
}
func (u User) signin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value(ContextDatabase).(*data.DB)
	isJSON := ctx.Value(ContextContentIsJSON).(bool)

	var data = new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})

	fmt.Println("isJson", isJSON)

	if isJSON {
		if err := ParseBody(r.Body, &data); err != nil {
			Respond(w, r, http.StatusBadRequest, err)
			return
		}
	} else {
		r.ParseForm()
		data.Email = r.Form.Get("email")
		data.Password = r.Form.Get("password")
	}

	user, err := db.Users.GetUserByEmail(data.Email)
	if err != nil {
		if isJSON {
			Respond(w, r, http.StatusInternalServerError, err)
		} else {
			http.Redirect(w, r, config.Current.SignInErrorRedirect, http.StatusSeeOther)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		Respond(w, r, http.StatusNotFound, fmt.Errorf("invalid email/password."))
		return
	}

	if isJSON {
		Respond(w, r, http.StatusOK, user)
	} else {
		ck := &http.Cookie{
			Name:  "X-API-KEY",
			Path:  "/",
			Value: user.Token,
		}

		// we set the cookie so they will have their authentication token on the next request.
		http.SetCookie(w, ck)

		http.Redirect(w, r, config.Current.SignInSuccessRedirect, http.StatusSeeOther)
	}
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
