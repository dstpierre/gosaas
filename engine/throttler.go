package engine

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dstpierre/gosaas/cache"
)

// Throttler middleware used to throttle and apply rate limit to requests.
func Throttler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var keys Auth

		ctx := r.Context()
		v := ctx.Value(ContextAuth)
		if v == nil {
			keys = Auth{}
		} else {
			a, ok := v.(Auth)
			if ok {
				keys = a
			}
		}

		key := fmt.Sprintf("%v", keys.AccountID)

		// TODO: Make this configurable
		count, err := cache.Throttle(key, 24*time.Hour)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: Make this configurable
		if count >= 9999 {
			// we get the expiration duration of this key so we can notify the user
			d, err := cache.GetThrottleExpiration(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if d.Seconds() > 0 {
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(d.Seconds())))
			}
			http.Error(w, fmt.Sprintf("you've reached your daily limit, retry in %v", d), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
