package routes

import (
	"encoding/json"
	"log"
	"net/http"

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

	guildID := ctx.Value(context_key.K("guildID")).(string)
	userID := vars["userID"]

	body := StealCreditsPatchRequest{}
	err := h.unmarshaler.UnmarshalReader(r.Body, &body)
	if err != nil {
		badRequest(rw, "Failed to parse request body. Please check the documentation.")
		return
	}
	operation, ok := repository.OperationFromString(*body.Operation)
	if !ok {
		badRequest(rw, "Invalid operation. Options are add (+), subtract (-), multiply (*), divide (/), and set (=). Use the symbol.")
		return
	}

	err = h.repo.AdjustStealCreditsForUser(guildID, userID, operation, *body.Credits)
	if err != nil {
		log.Printf("Error adjusting steal credits: %v", err)
		internalServerError(rw, "No changes were made.")
		return
	}

	credits, err := h.repo.GetStealCreditsForUser(guildID, userID)
	if err != nil {
		log.Printf("Error fetching steal credits for user: %v", err)
		respondOK(rw)
		return
	}

	responseJSON, err := json.Marshal(StealCreditsGetResponse{
		GuildID: guildID,
		UserID:  userID,
		Credits: credits,
	})

	if err != nil {
		log.Printf("Error marshalling GET steal credits response: %v", err)
		respondOK(rw)
		return
	}
	rw.Write(responseJSON)
}
