package gosaas

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/dstpierre/gosaas/cache"
	"github.com/dstpierre/gosaas/model"
	uuid "github.com/satori/go.uuid"
)

// Logger is a middleware that log requests information to stdout.
//
// If the request failed with a status code >= 300, a dump of the
// request will be saved into the cache store. You can investigate and replay
// the request in a development environment using this tool https://github.com/dstpierre/httpreplay.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ContextRequestStart, time.Now())
		ctx = context.WithValue(ctx, ContextRequestID, uuid.NewV4())

		dr, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println("unable to dump request", err)
		} else {
			ctx = context.WithValue(ctx, ContextRequestDump, dr)
		}

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

	v = ctx.Value(ContextRequestID)
	reqID, ok := v.(string)
	if !ok {
		reqID = "failed"
	}

	v = ctx.Value(ContextRequestDump)
	dr, ok := v.([]byte)
	if !ok {
		log.Println(path, "unable to retrieve the dump request data")
	} else {
		if statusCode >= http.StatusBadRequest {
			// we don't want to log 404 not found
			if statusCode != http.StatusNotFound {
				if dr != nil {
					if err := cache.LogWebRequest(reqID, dr); err != nil {
						log.Println("unable to save failed request", err)
					}
				}
			}
		}
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
		RequestID:  reqID,
	}

	go func(lr model.APIRequest) {
		if err := cache.LogRequest(lr); err != nil {
			// TODO: this should be reported somewhere else as well
			log.Println("error while logging request to Redis", err)
		}
	}(lr)
}
