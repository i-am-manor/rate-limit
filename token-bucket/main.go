package main

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/time/rate"
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

func rateLimiter(next func (w http.ResponseWriter, r *http.Request)) http.Handler {
	limiter := rate.NewLimiter(2, 4)
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			message := Message{
				Status: "Request Failed",
				Body:   "The API is at capacity, try again later!",
			}
			json.NewEncoder(w).Encode(&message)
			return
		} else {
			next(w, r)
		}
	})
}

func main () {
	http.Handle("/ping", rateLimiter(endpointHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error listening on port :8080", err)
	}
}
