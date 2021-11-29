package api

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/api/context_key"
	"github.com/oppzippy/BoostRequestBot/api/middleware"
	routes_v1 "github.com/oppzippy/BoostRequestBot/api/v1/routes"
	"github.com/oppzippy/BoostRequestBot/api/v2/routes"
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

	v2 := r.PathPrefix("/v2").Subrouter()
	v2.HandleFunc("/", routes.NotFoundHandler)
	v2.Handle("/users/{userID:[0-9]+}/stealCredits", routes.NewStealCreditsGetHandler(repo)).Methods("GET")
	v2.Handle("/users/{userID:[0-9]+}/stealCredits", routes.NewStealCreditsPatchHandler(repo)).Methods("PATCH")

	v2.Handle("/boostRequests/{boostRequestID}", routes.NewBoostRequestGetHandler(repo)).Methods("GET")
	v2.Handle("/boostRequests", routes.NewBoostRequestPostHandler(repo, brm)).Methods("POST")

	v2.Use(middleware.ContentTypeMiddleware("application/json"))
	v2.Use(middleware.JsonResponseMiddleware())
	v2.Use(middleware.APIKeyMiddleware(repo))
	v2.Use(middleware.RequireAuthorizationMiddleware())

	// v1

	v1 := r.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/", routes_v1.NotFoundHandler)
	v1.Handle("/users/{userID:[0-9]+}/stealCredits", routes_v1.NewStealCreditsGetHandler(repo)).Methods("GET")
	v1.Handle("/users/{userID:[0-9]+}/stealCredits", routes_v1.NewStealCreditsPatchHandler(repo)).Methods("PATCH")

	v1.Handle("/boostRequests/{boostRequestID}", routes_v1.NewBoostRequestGetHandler(repo)).Methods("GET")
	v1.Handle("/boostRequests", routes_v1.NewBoostRequestPostHandler(repo, brm)).Methods("POST")

	v1.Use(middleware.ContentTypeMiddleware("application/json"))
	v1.Use(middleware.JsonResponseMiddleware())
	v1.Use(middleware.APIKeyMiddleware(repo))
	v1.Use(middleware.RequireAuthorizationMiddleware())

	return r
}
