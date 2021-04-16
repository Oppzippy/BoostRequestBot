package api

import (
	"context"
	"net"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type contextKey string

var authBearerRegex *regexp.Regexp = regexp.MustCompile("^Bearer (.*)")

func NewWebAPI(repo repository.Repository) *http.Server {
	server := http.Server{
		Addr:         os.Getenv("HTTP_LISTEN_ADDRESS"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router(repo),
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.WithValue(context.Background(), contextKey("repo"), repo)
			return ctx
		},
	}
	return &server
}

func router(repo repository.Repository) http.Handler {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/", notFoundHandler)
	v1.HandleFunc("/users/{userID:[0-9]+}/stealCredits", getStealCreditsHandler).Methods("GET")
	v1.HandleFunc("/users/{userID:[0-9]+}/stealCredits", setStealCreditsHandler).Methods("PUT")
	v1.HandleFunc("/users/{userID:[0-9]+}/stealCredits", adjustStealCreditsHandler).Methods("PATCH")

	v1.Use(contentTypeMiddleware("application/json"))
	v1.Use(apiKeyMiddleware(repo))
	v1.Use(requireAuthorizationMiddleware())

	return r
}

func notFoundHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte(`{"statusCode": 404,"error": "Not found", "message": "The requested API endpoint does not exist."}`))
}

func contentTypeMiddleware(contentType string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Add("Content-Type", contentType)
			next.ServeHTTP(rw, r)
		})
	}
}

func apiKeyMiddleware(repo repository.Repository) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			match := authBearerRegex.FindStringSubmatch(auth)
			if len(match) == 1 {
				key := match[0]
				apiKey, err := repo.GetAPIKey(key)

				if err != nil {
					ctx := context.WithValue(r.Context(), contextKey("isAuthorized"), true)
					ctx = context.WithValue(ctx, contextKey("guildID"), apiKey.GuildID)

					*r = *r.WithContext(ctx)
					next.ServeHTTP(rw, r)
					return
				}
			}
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte(`{"statusCode":401,"error":"Unauthorized","message":"You must specify an api key with the header Authorization: Bearer api_key"}`))
		})
	}
}

func requireAuthorizationMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			isAuthorized, ok := r.Context().Value(contextKey("isAuthorized")).(bool)
			if ok && isAuthorized {
				next.ServeHTTP(rw, r)
				return
			}
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte(`{"statusCode": 401, "error": "Unauthorized", "message": "You must authenticate with the HTTP header 'Authorization: Bearer YOUR_API_KEY'"}`))
		})
	}
}
