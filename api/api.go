package api

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/api/context_key"
	"github.com/oppzippy/BoostRequestBot/api/middleware"
	"github.com/oppzippy/BoostRequestBot/api/v3/routes"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func NewWebAPI(repo repository.Repository, brm *boost_request_manager.BoostRequestManager, listenAddress string) *http.Server {
	server := http.Server{
		Addr:         listenAddress,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router(repo, brm),
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.WithValue(context.Background(), context_key.Repository, repo)
			return ctx
		},
	}
	return &server
}

func router(repo repository.Repository, brm *boost_request_manager.BoostRequestManager) http.Handler {
	r := mux.NewRouter()

	v3 := r.PathPrefix("/v3").Subrouter()
	v3.HandleFunc("/", routes.NotFoundHandler)
	v3.Handle("/users/{userID:[0-9]+}/stealCredits", routes.NewStealCreditsGetHandler(repo)).Methods("GET")
	v3.Handle("/users/{userID:[0-9]+}/stealCredits", routes.NewStealCreditsPatchHandler(repo)).Methods("PATCH")

	v3.Handle("/boostRequests/{boostRequestID}", routes.NewBoostRequestGetHandler(repo)).Methods("GET")
	v3.Handle("/boostRequests", routes.NewBoostRequestPostHandler(repo, brm)).Methods("POST")

	v3.Use(middleware.APIKeyMiddleware(repo))
	v3.Use(middleware.RequireAuthorizationMiddleware())

	return r
}
