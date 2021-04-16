package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	ResponseCode int    `json:"responseCode"`
	Error        string `json:"error"`
	Message      string `json:"message"`
}

func internalServerError(message string) []byte {
	if message == "" {
		message = "An unexpected error has occurred."
	}
	response := ErrorResponse{
		ResponseCode: http.StatusInternalServerError,
		Error:        "Internal Server Error",
		Message:      message,
	}
	resp, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling internal server error json: %v", err)
	}
	return resp
}
