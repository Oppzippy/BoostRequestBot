package webhook

const (
	AdvertiserChosenEvent   = "advertiserChosen"
	AdvertiserChosenEventV2 = "advertiserChosenV2"
	AdvertiserChosenEventV3 = "advertiserChosenV3"
)

type WebhookEvent struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}
