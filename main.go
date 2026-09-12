package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Golang!")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "About page")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errMsg := map[string]string{
			"error": "invalid request body",
		}

		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(errMsg); err != nil {
			log.Printf("failed to encode error response: %v", err)
		}

		return
	}

	if input.Name == "" || input.Email == "" {
		errMsg := map[string]string{
			"error": "name and email are required",
		}

		w.WriteHeader(http.StatusBadRequest)

		if err := json.NewEncoder(w).Encode(errMsg); err != nil {
			log.Printf("failed to encode error response: %v", err)
		}

		return
	}

	id := uuid.NewString()

	response := map[string]string{
		"id":    id,
		"name":  input.Name,
		"email": input.Email,
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"status": "OK",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
	}

}

func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/users", createUserHandler)

	log.Println("Server running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
