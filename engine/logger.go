package engine

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dstpierre/gosaas/cache"
	"github.com/dstpierre/gosaas/data/model"
)

// Logger middleware that log request information
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ContextRequestStart, time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func logRequest(r *http.Request, statusCode int) {
	ctx := r.Context()
	v := ctx.Value(ContextOriginalPath)
	path, ok := v.(string)
	if !ok {
		path = r.URL.Path
	}

	v = ctx.Value(ContextRequestStart)
	if v == nil {
		return
	}

	if s, ok := v.(time.Time); ok {
		log.Println(time.Since(s), statusCode, r.Method, path)
	}

	keys, ok := ctx.Value(ContextAuth).(Auth)
	if !ok {
		return
	}

	lr := model.APIRequest{
		AccountID:  keys.AccountID,
		Requested:  time.Now(),
		StatusCode: statusCode,
		URL:        path,
		UserID:     keys.UserID,
	}

	go func(lr model.APIRequest) {
		if err := cache.LogRequest(lr); err != nil {
			// TODO: this should be reported somewhere else as well
			log.Println("error while logging request to Redis", err)
		}
	}(lr)
}
