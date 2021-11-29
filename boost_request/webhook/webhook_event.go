package webhook

const (
	AdvertiserChosenEvent   = "advertiserChosen"
	AdvertiserChosenEventV2 = "advertiserChosenV2"
)

type WebhookEvent struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}
