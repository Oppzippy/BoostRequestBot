package models

type OK struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message,omitempty"`
}
