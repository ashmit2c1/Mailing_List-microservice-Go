package server

import (
	"encoding/json"
	"log"
	"mailinglist/db"
	"net/http"
)

type JSONServer struct {
	DB *db.DB
}

func (s *JSONServer) AddSubscriberHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := s.DB.AddSubscriber(request.Email); err != nil {
		http.Error(w, "Could not add subscriber", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Subscriber added"})
}

func (s *JSONServer) ListSubscribersHandler(w http.ResponseWriter, r *http.Request) {
	emails, err := s.DB.ListSubscribers()
	if err != nil {
		http.Error(w, "Could not list subscribers", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(emails)
}

func StartJSONServer(db *db.DB) {
	server := &JSONServer{DB: db}

	http.HandleFunc("/add", server.AddSubscriberHandler)
	http.HandleFunc("/list", server.ListSubscribersHandler)

	log.Println("Starting JSON server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
