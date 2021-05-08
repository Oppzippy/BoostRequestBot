package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type StealCreditsGetResponse struct {
	GuildID string `json:"guildId"`
	UserID  string `json:"userId"`
	Credits int    `json:"credits"`
}

type StealCreditsPatchRequest struct {
	Credits   *int    `json:"credits" validate:"required"`
	Operation *string `json:"operation" validate:"required"`
}

const okResponse = `{"statusCode": 200, "message": "ok"}`

var validate = validator.New()

func getStealCreditsHandler(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	guildID := ctx.Value(contextKey("guildID")).(string)
	userID := vars["userID"]
	repo := ctx.Value(contextKey("repo")).(repository.Repository)

	credits, err := repo.GetStealCreditsForUser(guildID, userID)
	if err != nil {
		log.Printf("Error fetching steal credits for user: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(internalServerError(""))
		return
	}

	responseJSON, err := json.Marshal(StealCreditsGetResponse{
		GuildID: guildID,
		UserID:  userID,
		Credits: credits,
	})

	if err != nil {
		log.Printf("Error marshalling GET steal credits response: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(internalServerError(""))
		return
	}
	rw.Write(responseJSON)
}

func adjustStealCreditsHandler(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	guildID := ctx.Value(contextKey("guildID")).(string)
	userID := vars["userID"]
	repo := ctx.Value(contextKey("repo")).(repository.Repository)

	body := StealCreditsPatchRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		badRequest(rw, "Failed to parse request body. Please check the documentation.")
		return
	}
	err = validate.Struct(body)
	if err != nil {
		badRequest(rw, "Failed to parse request body. Please check the documentation.")
		return
	}
	operation, ok := repository.OperationFromString(*body.Operation)
	if !ok {
		resp := ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Bad Request",
			Message:    "Invalid operation. Options are add (+), subtract (-), multiply (*), divide (/), and set (=). Use the symbol.",
		}
		marshaledResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error marshalling bad request: %v", err)
		}
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(marshaledResp)
		return
	}

	err = repo.AdjustStealCreditsForUser(guildID, userID, operation, *body.Credits)
	if err != nil {
		log.Printf("Error adjusting steal credits: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(internalServerError("No changes were made."))
		return
	}
	rw.Write([]byte(okResponse))
}
