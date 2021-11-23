package models

type BoostRequest struct {
	ID                     string            `json:"id"`
	RequesterID            string            `json:"requesterId"`
	IsAdvertiserSelected   bool              `json:"isAdvertiserSelected"`
	AdvertiserID           string            `json:"advertiserId,omitempty"`
	BackendChannelID       string            `json:"backendChannelId"`
	BackendMessageID       string            `json:"backendMessageID"`
	Message                string            `json:"message"`
	Price                  int64             `json:"price,string,omitempty"`
	Discount               int64             `json:"discount,string,omitempty"`
	AdvertiserCut          int64             `json:"advertiserCut,string,omitempty"`
	AdvertiserRoleCuts     map[string]string `json:"advertiserRoleCuts,omitempty"`
	PreferredAdvertiserIDs []string          `json:"preferredAdvertiserIds,omitempty"`
	CreatedAt              string            `json:"createdAt"`
	AdvertiserSelectedAt   string            `json:"advertiserSelectedAt,omitempty"`
}

type BoostRequestPartial struct {
	RequesterID            string            `json:"requesterId" validate:"required"`
	BackendChannelID       string            `json:"backendChannelId" validate:"required"`
	Message                string            `json:"message" validate:"required"`
	Price                  int64             `json:"price,string,omitempty"`
	AdvertiserCut          int64             `json:"advertiserCut,string,omitempty"`
	AdvertiserRoleCuts     map[string]string `json:"advertiserRoleCuts,omitempty"`
	Discount               int64             `json:"discount,string,omitempty"`
	PreferredAdvertiserIDs []string          `json:"preferredAdvertiserIds,omitempty"`
}
