package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	nameQuery := r.URL.Query().Get("name")

	if nameQuery != "" {
		log.Printf("%s", nameQuery)

		var matches []User
		for _, u := range UserList {
			if strings.EqualFold(u.Name, nameQuery) {
				matches = append(matches, u)
			}
		}

		if len(matches) == 0 {
			errMsg := map[string]string{"error": "no user found with that name"}

			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(errMsg); err != nil {
				log.Printf("failed to encode error response: %v", err)
				return

			}
			return
		}

		if err := json.NewEncoder(w).Encode(matches); err != nil {
			log.Printf("failed to encode error response: %v", err)
			return
		}
		return
	}
	if err := json.NewEncoder(w).Encode(UserList); err != nil {
		log.Printf("failed to encode error response: %v", err)
	}
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
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
