package engine

import (
	"context"
	"log"
	"net/http"
	"time"
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
}
