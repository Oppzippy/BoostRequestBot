package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type middlewareJsonResponseType int

const MiddlewareJsonResponse = middlewareJsonResponseType(iota)

func JsonResponseMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(rw, r)
			response := r.Context().Value(MiddlewareJsonResponse)
			if response != nil {
				responseJSON, err := json.Marshal(response)

				if err != nil {
					log.Printf("Error marshalling GET steal credits response: %v", err)
					rw.WriteHeader(500)
					return
				}

				_, err = rw.Write(responseJSON)
				if err != nil {
					log.Printf("Error sending http response: %v", err)
				}
			}
		})
	}
}
