package routes

import (
	"context"
	"net/http"

	"github.com/oppzippy/BoostRequestBot/api/middleware"
	"github.com/oppzippy/BoostRequestBot/api/v3/models"
)

func NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	resp := models.GenericResponse{
		StatusCode: http.StatusNotFound,
		Error:      "Not Found",
		Message:    "The requested API endpoint does not exist.",
	}
	ctx := context.WithValue(r.Context(), middleware.MiddlewareJsonResponse, resp)
	*r = *r.Clone(ctx)
}
