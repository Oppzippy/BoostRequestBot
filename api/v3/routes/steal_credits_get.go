package routes

import (
	"log"
	"net/http"

	"github.com/oppzippy/BoostRequestBot/api/responder"

	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/api/context_key"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type StealCreditsGetResponse struct {
	GuildID string `json:"guildId"`
	UserID  string `json:"userId"`
	Credits int    `json:"credits"`
}

type StealCreditsGet struct {
	repo repository.Repository
}

func NewStealCreditsGetHandler(repo repository.Repository) *StealCreditsGet {
	return &StealCreditsGet{
		repo: repo,
	}
}

func (h *StealCreditsGet) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	guildID := ctx.Value(context_key.GuildID).(string)
	userID := vars["userID"]

	credits, err := h.repo.GetStealCreditsForUser(guildID, userID)
	if err != nil {
		log.Printf("Error fetching steal credits for user: %v", err)
		responder.RespondError(rw, http.StatusInternalServerError, "")
		return
	}

	responder.RespondJSON(rw, StealCreditsGetResponse{
		GuildID: guildID,
		UserID:  userID,
		Credits: credits,
	})
}
