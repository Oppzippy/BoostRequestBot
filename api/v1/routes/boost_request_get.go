package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/api/context_key"
	"github.com/oppzippy/BoostRequestBot/api/json_unmarshaler"
	"github.com/oppzippy/BoostRequestBot/api/middleware"
	"github.com/oppzippy/BoostRequestBot/api/v1/models"
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
		log.Printf("Error fetching boost request via api: %v", err)
		internalServerError(rw, r, "")
		return
	}

	var result *models.BoostRequest = models.FromRepositoryBoostRequest(br)
	ctx = context.WithValue(ctx, middleware.MiddlewareJsonResponse, result)
	*r = *r.Clone(ctx)
}
