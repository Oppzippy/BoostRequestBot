package routes

import (
	"context"
	"net/http"

	"github.com/oppzippy/BoostRequestBot/api/middleware"
	"github.com/oppzippy/BoostRequestBot/api/v1/models"
)

func respondOK(rw http.ResponseWriter, r *http.Request) {
	response := models.GenericResponse{
		StatusCode: http.StatusOK,
		Message:    "OK",
	}
	ctx := context.WithValue(r.Context(), middleware.MiddlewareJsonResponse, response)
	*r = *r.Clone(ctx)
}

func internalServerError(rw http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "An unexpected error has occurred."
	}
	response := models.GenericResponse{
		StatusCode: http.StatusInternalServerError,
		Error:      "Internal Server Error",
		Message:    message,
	}
	ctx := context.WithValue(r.Context(), middleware.MiddlewareJsonResponse, response)
	*r = *r.Clone(ctx)
	rw.WriteHeader(http.StatusInternalServerError)
}

func badRequest(rw http.ResponseWriter, r *http.Request, message string) {
	response := models.GenericResponse{
		StatusCode: http.StatusBadRequest,
		Error:      "Bad Request",
		Message:    message,
	}
	ctx := context.WithValue(r.Context(), middleware.MiddlewareJsonResponse, response)
	*r = *r.Clone(ctx)
	rw.WriteHeader(http.StatusBadRequest)
}

func notFound(rw http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "The specified item could not be found."
	}
	response := models.GenericResponse{
		StatusCode: http.StatusNotFound,
		Error:      "Not Found",
		Message:    message,
	}
	ctx := context.WithValue(r.Context(), middleware.MiddlewareJsonResponse, response)
	*r = *r.Clone(ctx)
	rw.WriteHeader(http.StatusNotFound)
}
