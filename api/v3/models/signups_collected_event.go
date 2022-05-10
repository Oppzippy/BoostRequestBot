package models

type SignupsCollectedEvent struct {
	BoostRequest *BoostRequest    `json:"boostRequest"`
	Signups      []SignupWithRoll `json:"signups"`
}

type SignupWithRoll struct {
	UserID string  `json:"userId"`
	Roll   float64 `json:"roll"`
}
