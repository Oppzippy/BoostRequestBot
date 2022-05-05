package middleware

import (
	"net/http"

	"github.com/oppzippy/BoostRequestBot/api/responder"

	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/api/context_key"
	"github.com/oppzippy/BoostRequestBot/api/v3/models"
)

// Requires API key middleware
func RequireAuthorizationMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			isAuthorized, ok := r.Context().Value(context_key.IsAuthorized).(bool)
			if ok && isAuthorized {
				next.ServeHTTP(rw, r)
				return
			}
			rw.WriteHeader(http.StatusUnauthorized)
			resp := models.Error{
				StatusCode: http.StatusUnauthorized,
				Error:      "Unauthorized",
				Message:    "You must authenticate with the HTTP header 'X-API-Key: YOUR_API_KEY'",
			}
			responder.RespondJSON(rw, resp)
		})
	}
}
