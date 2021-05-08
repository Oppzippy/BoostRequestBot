package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/api/context_key"
	"github.com/oppzippy/BoostRequestBot/api/routes"
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
					ctx := context.WithValue(r.Context(), context_key.K("isAuthorized"), true)
					ctx = context.WithValue(ctx, context_key.K("guildID"), apiKey.GuildID)

					*r = *r.WithContext(ctx)
					next.ServeHTTP(rw, r)
					return
				}
			}

			resp := routes.GenericResponse{
				StatusCode: http.StatusUnauthorized,
				Error:      "Unauthorized",
				Message:    "You must specify an api key with the header X-API-Key: your_api_key",
			}
			resp.Write(rw)
		})
	}
}
