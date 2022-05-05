package responder

import (
	"net/http"

	"github.com/oppzippy/BoostRequestBot/api/v3/models"
)

func RespondError(rw http.ResponseWriter, statusCode int, message string) {
	if message == "" {
		message = "An unknown error has occurred"
	}
	rw.WriteHeader(statusCode)
	resp := &models.Error{
		StatusCode: statusCode,
		Error:      http.StatusText(statusCode),
		Message:    message,
	}
	RespondJSON(rw, resp)
}

func RespondDetailedError(rw http.ResponseWriter, statusCode int, errorDetails any) {
	rw.WriteHeader(statusCode)
	resp := &models.DetailedError{
		StatusCode: statusCode,
		Error:      http.StatusText(statusCode),
		Details:    errorDetails,
	}
	RespondJSON(rw, resp)
}
