package models

type BoostRequest struct {
	ID                     string   `json:"id"`
	RequesterID            string   `json:"requesterId"`
	IsAdvertiserSelected   bool     `json:"isAdvertiserSelected"`
	AdvertiserID           string   `json:"advertiserId,omitempty"`
	BackendChannelID       string   `json:"backendChannelId"`
	BackendMessageID       string   `json:"backendMessageID"`
	Message                string   `json:"message"`
	Price                  int64    `json:"price,omitempty"`
	AdvertiserCut          int64    `json:"advertiserCut,omitempty"`
	PreferredAdvertiserIDs []string `json:"preferredAdvertiserIds"`
	CreatedAt              string   `json:"createdAt"`
	AdvertiserSelectedAt   string   `json:"advertiserSelectedAt,omitempty"`
}

type BoostRequestPartial struct {
	RequesterID            string   `json:"requesterId" validate:"required"`
	BackendChannelID       string   `json:"BackendChannelId" validate:"required"`
	Message                string   `json:"message" validate:"required"`
	Price                  int64    `json:"price" validate:"required"`
	AdvertiserCut          int64    `json:"advertiserCut" validate:"required"`
	PreferredAdvertiserIDs []string `json:"preferredAdvertiserIds"`
}
