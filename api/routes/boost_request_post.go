package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/oppzippy/BoostRequestBot/api/context_key"
	"github.com/oppzippy/BoostRequestBot/api/json_unmarshaler"
	"github.com/oppzippy/BoostRequestBot/api/middleware"
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

	brPartial, err := boost_request_manager.FromModelBoostRequestPartial(&body)
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

	br, err := h.brm.CreateBoostRequest(brc, brPartial)
	if err != nil {
		log.Printf("Error creating boost request via api: %v", err)
		internalServerError(rw, r, "")
		return
	}

	br, err = h.repo.GetBoostRequestById(br.Channel.GuildID, *br.ExternalID)
	if err != nil {
		log.Printf("Error fetching boost request: %v", err)
		internalServerError(rw, r, "")
		return
	}

	var response *models.BoostRequest = models.FromRepositoryBoostRequest(br)
	ctx = context.WithValue(ctx, middleware.MiddlewareJsonResponse, response)
	*r = *r.Clone(ctx)
}
