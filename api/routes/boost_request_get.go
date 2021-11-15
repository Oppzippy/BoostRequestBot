package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/api/context_key"
	"github.com/oppzippy/BoostRequestBot/api/json_unmarshaler"
	"github.com/oppzippy/BoostRequestBot/api/middleware"
	"github.com/oppzippy/BoostRequestBot/api/models"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestGet struct {
	repo        repository.Repository
	unmarshaler *json_unmarshaler.Unmarshaler
}

func NewBoostRequestGetHandler(repo repository.Repository) *BoostRequestGet {
	return &BoostRequestGet{
		repo:        repo,
		unmarshaler: json_unmarshaler.New(),
	}
}

func (h *BoostRequestGet) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	guildID := ctx.Value(context_key.GuildID).(string)

	boostRequestID, err := uuid.Parse(vars["boostRequestID"])
	if err != nil {
		badRequest(rw, r, "Invalid uuid")
		return
	}

	br, err := h.repo.GetBoostRequestById(guildID, boostRequestID)
	if err == repository.ErrNoResults {
		notFound(rw, r, "That boost request does not exist.")
		return
	}
	if err != nil {
		internalServerError(rw, r, "")
		return
	}

	var advertiserSelectedAt string
	if !br.ResolvedAt.IsZero() {
		advertiserSelectedAt = br.ResolvedAt.Format(time.RFC3339)
	}

	result := &models.BoostRequest{
		ID:                     br.ExternalID.String(),
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
	}
	ctx = context.WithValue(ctx, middleware.MiddlewareJsonResponse, result)
	*r = *r.Clone(ctx)
}
