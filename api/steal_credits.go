package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type StealCreditsGetResponse struct {
	GuildID string `json:"guildId"`
	UserID  string `json:"userId"`
	Credits int    `json:"credits"`
}

const okResponse = `{"statusCode": 200, "message": "ok"}`

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

	r.ParseForm()
	creditsStr := r.PostForm.Get("credits")
	credits, err := strconv.Atoi(creditsStr)
	if creditsStr == "" || err != nil {
		resp := ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "Bad Request",
			Message:    "You must specify an integer number of credits to add in the POST body.",
		}
		marshaledResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error marshaling bad request: %v", err)
		}
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(marshaledResp)
		return
	}
	operation, ok := repository.OperationFromString(r.PostForm.Get("operation"))
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

	err = repo.AdjustStealCreditsForUser(guildID, userID, operation, credits)
	if err != nil {
		log.Printf("Error adjusting steal credits: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(internalServerError("No changes were made."))
		return
	}
	rw.Write([]byte(okResponse))
}
