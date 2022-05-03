//go:build integration
// +build integration

package api_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/oppzippy/BoostRequestBot/api"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository/database"
	"github.com/oppzippy/BoostRequestBot/initialization"
)

type stealCreditsGetResponse struct {
	GuildID *string `json:"guildId" validate:"required"`
	UserID  *string `json:"userId" validate:"required"`
	Credits *int    `json:"credits" validate:"required"`
}

type stealCreditsPatchRequest struct {
	Credits   int    `json:"credits"`
	Operation string `json:"operation"`
}

func TestStealCredits(t *testing.T) {
	db, err := initialization.GetDBC()
	if err != nil {
		t.Errorf("Error connecting to db: %v", err)
		return
	}
	repo := database.NewRepository(db)
	server := api.NewWebAPI(repo, nil, "localhost:8080")
	go server.ListenAndServe()
	client := &http.Client{}
	apiKey, err := repo.NewAPIKey("1")
	if err != nil {
		t.Errorf("Error creating api key: %v", err)
		return
	}

	t.Run("get credits", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://localhost:8080/v1/users/1/stealCredits", nil)
		if err != nil {
			t.Errorf("Error creating http request: %v", err)
			return
		}
		req.Header.Set("X-API-Key", apiKey.Key)
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("http error: %v", err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Got http status code %v, expected 200", resp.Status)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("error reading from server: %v", err)
			return
		}
		var r stealCreditsGetResponse
		err = json.Unmarshal(body, &r)
		if err != nil {
			t.Errorf("Error parsing json: %v", err)
			return
		}
		validate := validator.New()
		err = validate.Struct(r)
		if err != nil {
			t.Errorf("Invalid response schema: %v", err)
			return
		}
		if *r.Credits != 0 {
			t.Errorf("Expected 0 credits, got %d", r.Credits)
		}
	})
	t.Run("set credits", func(t *testing.T) {
		body := stealCreditsPatchRequest{
			Credits:   2,
			Operation: "=",
		}
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Error marshalling body: %v", err)
			return
		}
		req, err := http.NewRequest("PATCH", "http://localhost:8080/v1/users/2/stealCredits", bytes.NewBuffer(bodyJSON))
		if err != nil {
			t.Errorf("Error creating http request: %v", err)
			return
		}
		req.Header.Set("X-API-Key", apiKey.Key)
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("http error: %v", err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Got http status code code %v, expected 200", resp.Status)
		}
		var r stealCreditsGetResponse
		err = json.NewDecoder(resp.Body).Decode(&r)
		if err != nil {
			t.Errorf("Error fetching credits: %v", err)
			return
		}
		if *r.Credits != 2 {
			t.Errorf("Set credits to 2 but ended up with %v", *r.Credits)
		}
	})
	t.Run("add credits", func(t *testing.T) {

	})
	server.Close()
}
