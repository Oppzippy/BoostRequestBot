package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/api/context_key"
	"github.com/oppzippy/BoostRequestBot/api/routes"
)

// Requires API key middleware
func RequireAuthorizationMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			isAuthorized, ok := r.Context().Value(context_key.K("isAuthorized")).(bool)
			if ok && isAuthorized {
				next.ServeHTTP(rw, r)
				return
			}
			rw.WriteHeader(http.StatusUnauthorized)
			resp, err := json.Marshal(routes.GenericResponse{
				StatusCode: http.StatusUnauthorized,
				Error:      "Unauthorized",
				Message:    "You must authenticate with the HTTP header 'X-API-Key: YOUR_API_KEY'",
			})
			if err != nil {
				log.Printf("Error marshalling error response: %v", err)
			}
			rw.Write(resp)
		})
	}
}
