package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/api/context_key"
	"github.com/oppzippy/BoostRequestBot/api/v3/models"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func APIKeyMiddleware(repo repository.Repository) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key == "" {
				key = r.URL.Query().Get("api_key")
			}
			if key != "" {
				apiKey, err := repo.GetAPIKey(key)
				if err == nil {
					ctx := context.WithValue(r.Context(), context_key.IsAuthorized, true)
					ctx = context.WithValue(ctx, context_key.GuildID, apiKey.GuildID)

					*r = *r.Clone(ctx)
					next.ServeHTTP(rw, r)
					return
				}
			}

			resp := models.GenericResponse{
				StatusCode: http.StatusUnauthorized,
				Error:      "Unauthorized",
				Message:    "You must specify an api key with the header X-API-Key: your_api_key",
			}
			ctx := context.WithValue(r.Context(), MiddlewareJsonResponse, resp)
			*r = *r.Clone(ctx)
		})
	}
}
