package engine

import (
	"context"
	"net/http"
)

func Language(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the request language
		lng := "en"
		ck, err := r.Cookie("lng")
		if err == nil {
			lng = ck.Value
		}

		ctx := context.WithValue(r.Context(), ContextLanguage, lng)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
