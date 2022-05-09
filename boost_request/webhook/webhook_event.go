package webhook

type WebhookEventType string

const (
	AdvertiserChosenEvent   WebhookEventType = "advertiserChosen"
	AdvertiserChosenEventV2                  = "advertiserChosenV2"
	AdvertiserChosenEventV3                  = "advertiserChosenV3"
	SignupsCollectedEventV3                  = "signupsCollectedV3"
)

type WebhookEvent struct {
	Event   WebhookEventType `json:"event"`
	Payload interface{}      `json:"payload"`
}
