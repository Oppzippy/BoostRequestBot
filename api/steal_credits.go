package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type stealCreditsGetResponse struct {
	UserID  string `json:"userId"`
	Credits int    `json:"credits"`
}

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

	responseJSON, err := json.Marshal(stealCreditsGetResponse{
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
			ResponseCode: http.StatusBadRequest,
			Error:        "Bad Request",
			Message:      "You must specify an integer number of credits to add in the POST body.",
		}
		marshaledResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error marshaling POST steal credits bad request: %v", err)
		}
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(marshaledResp)
		return
	}

	err = repo.AdjustStealCreditsForUser(guildID, userID, repository.OperationAdd, credits)
	if err != nil {
		log.Printf("Error adjusting steal credits: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(internalServerError("No changes were made."))
		return
	}
	rw.Write([]byte(`{"statusCode": 200, "message": "ok"}`))
}

func setStealCreditsHandler(rw http.ResponseWriter, r *http.Request) {
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
			ResponseCode: http.StatusBadRequest,
			Error:        "Bad Request",
			Message:      "You must specify an integer number of boost request credits to set.",
		}
		marshaledResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error marshaling PUT steal credits bad request: %v", err)
		}
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(marshaledResp)
		return
	}

	err = repo.UpdateStealCreditsForUser(guildID, userID, credits)
	if err != nil {
		log.Printf("Error setting steal credits: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(internalServerError("No changes were made."))
		return
	}
	rw.Write([]byte(`{"statusCode": 200, "message": "ok"}`))
}
