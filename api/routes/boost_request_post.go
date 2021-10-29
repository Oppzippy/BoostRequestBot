package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/oppzippy/BoostRequestBot/api/json_unmarshaler"
	"github.com/oppzippy/BoostRequestBot/api/models"
	"github.com/oppzippy/BoostRequestBot/boost_request"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestPostResponse struct {
	GuildID      string `json:"guildId"`
	UserID       string `json:"userId"`
	BoostRequest int    `json:"credits"`
}

type BoostRequestPost struct {
	repo        repository.Repository
	unmarshaler *json_unmarshaler.Unmarshaler
	brm         *boost_request.BoostRequestManager
}

func NewBoostRequestPostHandler(repo repository.Repository) *BoostRequestPost {
	return &BoostRequestPost{
		repo:        repo,
		unmarshaler: json_unmarshaler.New(),
	}
}

func (h *BoostRequestPost) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	// guildID := ctx.Value(context_key.K("guildID")).(string)

	body := models.BoostRequestPartial{}
	err := h.unmarshaler.UnmarshalReader(r.Body, &body)
	if err != nil {
		badRequest(rw, r, "Failed to parse request body. Please check the documentation.")
		return
	}

	br, err := h.brm.CreateBoostRequest(&repository.BoostRequestChannel{}, boost_request.BoostRequestPartial{
		RequesterID: body.RequesterID,
		Message:     body.Message,
	})

	responseJSON, err := json.Marshal(models.BoostRequest{
		Id:                   br.ExternalID.String(), // Since we created the boost request after the UUID update, this will never be null
		RequesterID:          br.RequesterID,
		IsAdvertiserSelected: br.IsResolved,
		AdvertiserID:         br.AdvertiserID,
		BackendChannelID:     br.Channel.BackendChannelID,
		BackendMessageID:     br.BackendMessageID,
		Message:              br.Message,
		// Price:                  br.Price,
		// AdvertiserCut:          br.AdvertiserCut,
		// PreferredAdvertiserIds: br.PreferredAdvertiserIds,
		CreatedAt:            br.CreatedAt.Format(time.RFC3339),
		AdvertiserSelectedAt: br.ResolvedAt.Format(time.RFC3339),
	})
	if err != nil {
		log.Printf("Error marshalling POST boost request response: %v", err)
		internalServerError(rw, r, "")
		return
	}

	_, err = rw.Write(responseJSON)
	if err != nil {
		log.Printf("Error sending http response: %v", err)
	}
}
