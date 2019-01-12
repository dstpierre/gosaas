package gosaas

import (
	"context"
	"net/http"
)

// Language is a middleware handling the language cookie named "lng".Language
//
// This is used in HTML templates and Go code when using the Translate
// function. You need to create a language files inside a directory named
// languagepack (i.e. en.json, fr.json).
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
