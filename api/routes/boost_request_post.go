package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/oppzippy/BoostRequestBot/api/context_key"
	"github.com/oppzippy/BoostRequestBot/api/json_unmarshaler"
	"github.com/oppzippy/BoostRequestBot/api/models"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestPost struct {
	repo        repository.Repository
	unmarshaler *json_unmarshaler.Unmarshaler
	brm         *boost_request_manager.BoostRequestManager
}

func NewBoostRequestPostHandler(repo repository.Repository, brm *boost_request_manager.BoostRequestManager) *BoostRequestPost {
	return &BoostRequestPost{
		repo:        repo,
		brm:         brm,
		unmarshaler: json_unmarshaler.New(),
	}
}

func (h *BoostRequestPost) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildID := ctx.Value(context_key.GuildID).(string)

	body := models.BoostRequestPartial{}
	err := h.unmarshaler.UnmarshalReader(r.Body, &body)
	if err != nil {
		badRequest(rw, r, "Failed to parse request body. Please check the documentation.")
		return
	}

	// TODO check to make sure the channel is actually in the specified guild
	brc := &repository.BoostRequestChannel{
		FrontendChannelID: "",
		GuildID:           guildID,
		BackendChannelID:  body.BackendChannelID,
		UsesBuyerMessage:  false,
		SkipsBuyerDM:      false,
	}
	err = h.repo.InsertBoostRequestChannel(brc)
	if err != nil {
		log.Printf("error inserting internal boost request channel (no frontend channel): %v", err)
		internalServerError(rw, r, "")
		return
	}

	br, err := h.brm.CreateBoostRequest(brc, boost_request_manager.BoostRequestPartial{
		RequesterID:            body.RequesterID,
		Type:                   body.Type,
		Message:                body.Message,
		PreferredAdvertiserIDs: body.PreferredAdvertiserIDs,
		Price:                  body.Price,
		AdvertiserCut:          body.AdvertiserCut,
	})
	if err != nil {
		log.Printf("Error creating boost request via api: %v", err)
		internalServerError(rw, r, "")
		return
	}

	var advertiserSelectedAt string
	if !br.ResolvedAt.IsZero() {
		advertiserSelectedAt = br.ResolvedAt.Format(time.RFC3339)
	}

	responseJSON, err := json.Marshal(models.BoostRequest{
		ID:                     br.ExternalID.String(), // Since we created the boost request after the UUID update, this will never be null
		RequesterID:            br.RequesterID,
		IsAdvertiserSelected:   br.IsResolved,
		AdvertiserID:           br.AdvertiserID,
		BackendChannelID:       br.Channel.BackendChannelID,
		BackendMessageID:       br.BackendMessageID,
		Type:                   br.Type,
		Message:                br.Message,
		Price:                  br.Price,
		Discount:               br.Discount,
		AdvertiserCut:          br.AdvertiserCut,
		PreferredAdvertiserIDs: br.PreferredAdvertiserIDs,
		CreatedAt:              br.CreatedAt.Format(time.RFC3339),
		AdvertiserSelectedAt:   advertiserSelectedAt,
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
