package models

type DetailedError struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Details    any    `json:"details,omitempty"`
}
