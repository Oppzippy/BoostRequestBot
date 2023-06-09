package routes

import (
	"log"
	"net/http"

	"github.com/oppzippy/BoostRequestBot/api/responder"

	"github.com/oppzippy/BoostRequestBot/api/context_key"
	"github.com/oppzippy/BoostRequestBot/api/json_unmarshaler"
	"github.com/oppzippy/BoostRequestBot/api/v3/models"
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
	// TODO check to make sure the channel is actually in the specified guild
	ctx := r.Context()

	guildID := ctx.Value(context_key.GuildID).(string)

	body := models.BoostRequestPartial{}
	err := h.unmarshaler.UnmarshalReader(r.Body, &body)
	if err != nil {
		if validationError, ok := err.(json_unmarshaler.ValidationError); ok {
			responder.RespondDetailedError(rw, http.StatusBadRequest, validationError.TranslatedErrors)
		} else {
			responder.RespondError(rw, http.StatusBadRequest, "Failed to parse request body. Please check the documentation.")
		}
		return
	}

	brPartial, err := boost_request_manager.FromModelBoostRequestPartial(&body)
	if err != nil {
		responder.RespondError(rw, http.StatusBadRequest, "Failed to parse request body. Please check the documentation.")
		return
	}
	brPartial.GuildID = guildID

	br, err := h.brm.CreateBoostRequest(nil, brPartial)
	if err != nil {
		log.Printf("Error creating boost request via api: %v", err)
		responder.RespondError(rw, http.StatusInternalServerError, "")
		return
	}

	br, err = h.repo.GetBoostRequestById(br.GuildID, *br.ExternalID)
	if err != nil {
		log.Printf("Error fetching boost request: %v", err)
		responder.RespondError(rw, http.StatusInternalServerError, "")
		return
	}

	response := models.FromRepositoryBoostRequest(br)
	responder.RespondJSON(rw, response)
}
