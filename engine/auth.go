package engine

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dstpierre/gosaas/cache"
	"github.com/dstpierre/gosaas/data"
	"github.com/dstpierre/gosaas/data/model"
)

// Auth represents an authenticated user
type Auth struct {
	AccountID model.Key
	UserID    model.Key
	Email     string
	Role      model.Roles
}

// Authenticator middleware used to authenticate requests
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		mr := ctx.Value(ContextMinimumRole).(model.Roles)

		key, pat, err := extractKeyFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ca := &cache.Auth{}

		// do we have this key on cache already
		var a Auth
		if err := ca.Exists(key, &a); err != nil {
			log.Println("error while trying to get cache auth", err)
		}

		if len(a.Email) > 0 {
			ctx = context.WithValue(ctx, ContextAuth, a)
		} else {
			fmt.Println("database call for auth", key)
			db := ctx.Value(ContextDatabase).(*data.DB)

			id, _ := model.ParseToken(key)
			acct, user, err := db.Users.Auth(id, key, pat)
			if err != nil {
				http.Error(w, "invalid token key", http.StatusUnauthorized)
				return
			}

			a.AccountID = acct.ID
			a.Email = user.Email
			a.UserID = user.ID

			// save it to cache
			ca.Set(key, a, 30*time.Second)

			ctx = context.WithValue(ctx, ContextAuth, a)
		}

		// we authorize the request
		if mr < a.Role {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractKeyFromRequest(r *http.Request) (key string, pat bool, err error) {
	// first let's look if the X-API-KEY is present in the HTTP header
	key = r.Header.Get("X-API-KEY")
	if len(key) > 0 {
		return
	}

	// check the query string
	key = r.URL.Query().Get("key")
	if len(key) > 0 {
		return
	}

	// check if we are supplying basic auth
	authorization := r.Header.Get("Authorization")
	s := strings.SplitN(authorization, " ", 2)
	if len(s) != 2 {
		err = fmt.Errorf("invalid basic authentication format: %s - you must provide Basic base64token", authorization)
		return
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		err = fmt.Errorf("invalid basic authentication format: %s - you must provide Basic base64token", authorization)
		return
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		err = fmt.Errorf("invalid basic authentication, your token should be _:access_token - got %s", string(b))
		return
	}

	key = pair[1]
	pat = true

	return
}
