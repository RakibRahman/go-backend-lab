package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetUserListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	limit := 20
	page := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	offset := page * limit
	term := r.URL.Query().Get("term")

	usersResponse, dbErr := h.repo.GetUsers(r.Context(), limit, offset, term)

	if dbErr != nil {
		log.Printf("failed to get users list: %v", dbErr)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to get users",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(usersResponse); err != nil {
		log.Printf("failed to encode error response: %v", err)
	}
}

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")

	for _, user := range UserList {
		if user.ID == id {
			if err := json.NewEncoder(w).Encode(user); err != nil {
				log.Printf("failed to encode response: %v", err)
			}
			return
		}
	}

	errMsg := map[string]string{
		"error": "user not found",
	}

	w.WriteHeader(http.StatusNotFound)

	if err := json.NewEncoder(w).Encode(errMsg); err != nil {
		log.Printf("failed to encode error response: %v", err)
	}
}

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if input.Name == "" || input.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "name and email are required",
		})
		return
	}

	createdUser, err := h.repo.CreateUser(r.Context(), input)

	if err != nil {
		log.Printf("failed to create user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to create user",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		log.Printf("failed to encode response: %v", err)
	}

}

func DeleteUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	for index, user := range UserList {
		if user.ID == id {
			UserList = append(UserList[:index], UserList[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	errMsg := map[string]string{
		"error": "user not found",
	}
	w.WriteHeader(http.StatusNotFound)

	if err := json.NewEncoder(w).Encode(errMsg); err != nil {
		log.Printf("failed to encode error response: %v", err)
	}

}

func UpdateUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	var input UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	for i := range UserList {
		if UserList[i].ID == id {
			if input.Name != nil {
				UserList[i].Name = *input.Name
			}

			if input.Email != nil {
				UserList[i].Email = *input.Email
			}

			if err := json.NewEncoder(w).Encode(UserList[i]); err != nil {
				log.Printf("failed to encode error response: %v", err)
			}
			return
		}
	}

	http.Error(w, "user not found", http.StatusNotFound)
}
