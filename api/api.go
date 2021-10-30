package api

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/api/context_key"
	"github.com/oppzippy/BoostRequestBot/api/middleware"
	"github.com/oppzippy/BoostRequestBot/api/routes"
	"github.com/oppzippy/BoostRequestBot/boost_request"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func NewWebAPI(repo repository.Repository, brm *boost_request.BoostRequestManager, listenAddress string) *http.Server {
	server := http.Server{
		Addr:         listenAddress,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router(repo),
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.WithValue(context.Background(), context_key.Repository, repo)
			ctx = context.WithValue(ctx, context_key.BooostRequestManager, brm)
			return ctx
		},
	}
	return &server
}

func router(repo repository.Repository) http.Handler {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/", routes.NotFoundHandler)
	v1.Handle("/users/{userID:[0-9]+}/stealCredits", routes.NewStealCreditsGetHandler(repo)).Methods("GET")
	v1.Handle("/users/{userID:[0-9]+}/stealCredits", routes.NewStealCreditsPatchHandler(repo)).Methods("PATCH")

	// v1.Handle("/boostRequests/{boostRequestID}", routes.NewBoostRequestGetHandler(repo)).Methods("GET")
	v1.Handle("/boostRequests", routes.NewBoostRequestPostHandler(repo)).Methods("POST")

	v1.Use(middleware.ContentTypeMiddleware("application/json"))
	v1.Use(middleware.JsonResponseMiddleware())
	v1.Use(middleware.APIKeyMiddleware(repo))
	v1.Use(middleware.RequireAuthorizationMiddleware())

	return r
}
