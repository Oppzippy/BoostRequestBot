// +build integration

package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/oppzippy/BoostRequestBot/api"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository/database"
	"github.com/oppzippy/BoostRequestBot/initialization"
)

func TestStealCredits(t *testing.T) {
	db, err := initialization.GetDBC()
	if err != nil {
		t.Errorf("Error connecting to db: %v", err)
		return
	}
	repo := database.NewRepository(db)
	server := api.NewWebAPI(repo, "localhost:8080")
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
		req.Header.Add("X-API-Key", apiKey.Key)
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("http error: %v", err)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("error reading from server: %v", err)
			return
		}
		var r api.StealCreditsGetResponse
		err = json.Unmarshal(body, &r)
		if err != nil {
			t.Errorf("Error parsing json: %v", err)
			return
		}
		if r.Credits != 0 {
			t.Errorf("Expected 0 credits, got %d", r.Credits)
		}
	})
	t.Run("set credits", func(t *testing.T) {

	})
	t.Run("add credits", func(t *testing.T) {

	})
	server.Close()
}
