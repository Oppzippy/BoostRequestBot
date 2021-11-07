package webhook

const (
	AdvertiserChosenEvent = "advertiserChosen"
)

type WebhookEvent struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}
