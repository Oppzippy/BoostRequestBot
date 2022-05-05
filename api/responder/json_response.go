package responder

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondJSON(rw http.ResponseWriter, response any) {
	responseJSON, err := json.Marshal(response)

	if err != nil {
		log.Printf("error marshalling JSON: %v", err)
		rw.WriteHeader(500)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	_, err = rw.Write(responseJSON)
	if err != nil {
		log.Printf("error sending http response: %v", err)
	}
}
