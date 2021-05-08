package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`
}

func internalServerError(message string) []byte {
	if message == "" {
		message = "An unexpected error has occurred."
	}
	response := ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Error:      "Internal Server Error",
		Message:    message,
	}
	resp, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling internal server error json: %v", err)
	}
	return resp
}

func badRequest(rw http.ResponseWriter, message string) {
	resp := ErrorResponse{
		StatusCode: http.StatusBadRequest,
		Error:      "Bad Request",
		Message:    message,
	}
	marshaledResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshaling bad request: %v", err)
	}
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write(marshaledResp)
}
