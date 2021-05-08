package routes

import (
	"net/http"
)

func NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	resp := GenericResponse{
		StatusCode: http.StatusNotFound,
		Error:      "Not Found",
		Message:    "The requested API endpoint does not exist.",
	}
	resp.Write(rw)
}
