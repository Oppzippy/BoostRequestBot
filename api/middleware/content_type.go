package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

func ContentTypeMiddleware(contentType string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Add("Content-Type", contentType)
			next.ServeHTTP(rw, r)
		})
	}
}
