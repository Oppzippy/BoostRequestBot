package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

type GenericResponse struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error,omitempty"`
	Message    string `json:"message"`
}

func (er *GenericResponse) Write(rw http.ResponseWriter) {
	marshaledResp, err := json.Marshal(er)
	if err != nil {
		log.Printf("Error marshalling internal server error json: %v", err)
	}
	rw.WriteHeader(er.StatusCode)
	rw.Write(marshaledResp)
}

func respondOK(rw http.ResponseWriter) {
	response := GenericResponse{
		StatusCode: http.StatusOK,
		Message:    "OK",
	}
	response.Write(rw)
}

func internalServerError(rw http.ResponseWriter, message string) {
	if message == "" {
		message = "An unexpected error has occurred."
	}
	response := GenericResponse{
		StatusCode: http.StatusInternalServerError,
		Error:      "Internal Server Error",
		Message:    message,
	}
	response.Write(rw)
}

func badRequest(rw http.ResponseWriter, message string) {
	resp := GenericResponse{
		StatusCode: http.StatusBadRequest,
		Error:      "Bad Request",
		Message:    message,
	}
	resp.Write(rw)
}
