package main

import (
	"encoding/json"
	"log"
	"net/http"

	tollbooth "github.com/didip/tollbooth/v7"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func endpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	message := Message{
		Status: "Successful",
		Body:   "Hi! You've reached the API.",
	}
	err := json.NewEncoder(w).Encode(&message)
	if err != nil {
		return
	}
}

func main() {
	failedMessage := Message{
		Status: "Request Failed",
		Body:   "The API is at capacity, try again later!",
	}
	jsonMessage, _ := json.Marshal(failedMessage)

	limiter := tollbooth.NewLimiter(4, nil)
	limiter.SetMessageContentType("application/json")
	limiter.SetMessage(string(jsonMessage))

	http.Handle("/ping", tollbooth.LimitFuncHandler(limiter, endpointHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error listening on port :8080", err)
	}
}
