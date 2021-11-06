package webhook

const (
	AdvertiserChosenEvent = "advertiserChosen"
)

type WebhookEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}
