package webhook

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type WebhookManager struct {
	repo             repository.Repository
	trigger          chan struct{}
	triggerDestroy   chan struct{}
	retryLoop        chan struct{}
	retryLoopDestroy chan struct{}
	httpClient       *http.Client
}

func NewWebhookManager(repo repository.Repository) *WebhookManager {
	wm := &WebhookManager{
		repo:             repo,
		trigger:          make(chan struct{}, 1),
		triggerDestroy:   make(chan struct{}),
		retryLoop:        make(chan struct{}),
		retryLoopDestroy: make(chan struct{}),
		httpClient:       &http.Client{},
	}

	go func() {
		defer close(wm.triggerDestroy)
		for {
			_, ok := <-wm.trigger
			if !ok {
				return
			}
			wm.sendQueuedWebhooks()
		}
	}()

	go func() {
		defer close(wm.retryLoopDestroy)
		for {
			select {
			case wm.trigger <- struct{}{}:
			default:
			}

			select {
			case <-time.After(time.Minute * 15):
			case <-wm.retryLoop:
				return
			}
		}
	}()

	return wm
}

func (wm *WebhookManager) Destroy() {
	close(wm.retryLoop)
	<-wm.retryLoopDestroy
	close(wm.trigger)
	<-wm.triggerDestroy
}

func (wm *WebhookManager) QueueToSend(guildID string, event *WebhookEvent) error {
	webhook, err := wm.repo.GetWebhook(guildID)
	if err == repository.ErrNoResults {
		return nil
	}
	if err != nil {
		return err
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = wm.repo.InsertQueuedWebhook(webhook, string(body))
	if err != nil {
		return err
	}
	select {
	case wm.trigger <- struct{}{}:
	default:
	}
	return nil
}

func (wm *WebhookManager) sendQueuedWebhooks() error {
	queuedWebhooks, err := wm.repo.GetQueuedWebhooks()
	if err != nil {
		return err
	}
	for _, queuedWebhook := range queuedWebhooks {
		if queuedWebhook.LatestAttempt == nil || queuedWebhook.LatestAttempt.Add(time.Hour).After(time.Now()) {
			wm.sendQueuedWebhook(queuedWebhook)
		}
	}
	return nil
}

func (wm *WebhookManager) sendQueuedWebhook(queuedWebhook *repository.QueuedWebhookRequest) {
	startTime := time.Now()
	resp, err := wm.httpClient.Post(
		queuedWebhook.Webhook.URL,
		"application/json",
		bytes.NewBufferString(queuedWebhook.Body),
	)

	var statusCode int
	if err != nil {
		log.Printf("webhook: post request failed: %v", err)
	} else {
		statusCode = resp.StatusCode
		err = resp.Body.Close()
		if err != nil {
			log.Printf("webhook: failed to close body")
		}
	}

	wm.repo.InsertWebhookAttempt(repository.WebhookAttempt{
		QueuedWebhookRequest: *queuedWebhook,
		StatusCode:           statusCode,
		CreatedAt:            startTime,
	})
}
