package controllers

import (
	"net/http"
	"net/http/httptest"
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	api := &API{}

	rec := httptest.NewRecorder()

	api.ServeHTTP(rec, req)
	return rec
}
