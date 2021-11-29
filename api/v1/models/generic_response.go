package models

type GenericResponse struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error,omitempty"`
	Message    string `json:"message"`
}
