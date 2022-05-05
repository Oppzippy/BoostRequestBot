package routes

import (
	"log"
	"net/http"

	"github.com/oppzippy/BoostRequestBot/api/v3/models"

	"github.com/oppzippy/BoostRequestBot/api/responder"

	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/api/context_key"
	"github.com/oppzippy/BoostRequestBot/api/json_unmarshaler"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type StealCreditsPatchRequest struct {
	Credits   *int    `json:"credits" validate:"required"`
	Operation *string `json:"operation" validate:"required"`
}

type StealCreditsPatch struct {
	repo        repository.Repository
	unmarshaler *json_unmarshaler.Unmarshaler
}

func NewStealCreditsPatchHandler(repo repository.Repository) *StealCreditsPatch {
	return &StealCreditsPatch{
		repo:        repo,
		unmarshaler: json_unmarshaler.New(),
	}
}

func (h *StealCreditsPatch) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	guildID := ctx.Value(context_key.GuildID).(string)
	userID := vars["userID"]

	body := StealCreditsPatchRequest{}
	err := h.unmarshaler.UnmarshalReader(r.Body, &body)
	if err != nil {
		if validationError, ok := err.(json_unmarshaler.ValidationError); ok {
			responder.RespondDetailedError(rw, http.StatusBadRequest, validationError.TranslatedErrors)
		} else {
			responder.RespondError(rw, http.StatusBadRequest, "Failed to parse request body. Please check the documentation.")
		}
		return
	}
	operation, ok := repository.OperationFromString(*body.Operation)
	if !ok {
		responder.RespondError(rw, http.StatusBadRequest, "Invalid operation. Options are add (+), subtract (-), multiply (*), divide (/), and set (=). Use the symbol.")
		return
	}

	err = h.repo.AdjustStealCreditsForUser(guildID, userID, operation, *body.Credits)
	if err != nil {
		log.Printf("Error adjusting steal credits: %v", err)
		responder.RespondError(rw, http.StatusInternalServerError, "No changes were made.")
		return
	}

	credits, err := h.repo.GetStealCreditsForUser(guildID, userID)
	if err != nil {
		log.Printf("Error fetching steal credits for user: %v", err)
		responder.RespondJSON(rw, models.OK{
			StatusCode: http.StatusOK,
			Message:    "OK",
		})
		return
	}

	responder.RespondJSON(rw, StealCreditsGetResponse{
		GuildID: guildID,
		UserID:  userID,
		Credits: credits,
	})
}
