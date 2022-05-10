package models

type SignupsCollectedEvent struct {
	BoostRequest *BoostRequest    `json:"boostRequest"`
	Signups      []SignupWithRoll `json:"signups"`
}

type SignupWithRoll struct {
	UserID string
	Roll   float64
}
