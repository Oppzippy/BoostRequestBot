package api

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type contextKey string

func NewWebAPI(repo repository.Repository, listenAddress string) *http.Server {
	server := http.Server{
		Addr:         listenAddress,
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
	v1.HandleFunc("/users/{userID:[0-9]+}/stealCredits", adjustStealCreditsHandler).Methods("PATCH")

	v1.Use(contentTypeMiddleware("application/json"))
	v1.Use(apiKeyMiddleware(repo))
	v1.Use(requireAuthorizationMiddleware())

	return r
}

func notFoundHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	resp, err := json.Marshal(ErrorResponse{
		StatusCode: http.StatusNotFound,
		Error:      "Not Found",
		Message:    "The requested API endpoint does not exist.",
	})
	if err != nil {
		log.Printf("Error marshalling error response: %v", err)
	}
	rw.Write(resp)
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
			key := r.Header.Get("X-API-Key")
			if key == "" {
				key = r.URL.Query().Get("api_key")
			}
			if key != "" {
				apiKey, err := repo.GetAPIKey(key)
				if err == nil {
					ctx := context.WithValue(r.Context(), contextKey("isAuthorized"), true)
					ctx = context.WithValue(ctx, contextKey("guildID"), apiKey.GuildID)

					*r = *r.WithContext(ctx)
					next.ServeHTTP(rw, r)
					return
				}
			}
			rw.WriteHeader(http.StatusUnauthorized)
			resp, err := json.Marshal(ErrorResponse{
				StatusCode: http.StatusUnauthorized,
				Error:      "Unauthorized",
				Message:    "You must specify an api key with the header X-API-Key: your_api_key",
			})
			if err != nil {
				log.Printf("Error marshalling error response: %v", err)
			}
			rw.Write(resp)
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
			resp, err := json.Marshal(ErrorResponse{
				StatusCode: http.StatusUnauthorized,
				Error:      "Unauthorized",
				Message:    "You must authenticate with the HTTP header 'Authorization: Bearer YOUR_API_KEY'",
			})
			if err != nil {
				log.Printf("Error marshalling error response: %v", err)
			}
			rw.Write(resp)
		})
	}
}
