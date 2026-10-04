package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func queryIntDefault(r *http.Request, key string, fallback int) int {
	if v := r.URL.Query().Get(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return fallback
}

func (h *Handler) GetUserListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	limit := queryIntDefault(r, "limit", 20)
	page := queryIntDefault(r, "page", 0)
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

func (h *Handler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")

	user, err := h.repo.GetUserByID(r.Context(), id)

	if err != nil {
		log.Printf("failed to get users list: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "user not found",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
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

func (h *Handler) UpdateUserByIDHandler(w http.ResponseWriter, r *http.Request) {
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

	if input.Name != nil && *input.Name == "" {
		http.Error(w, "name cannot be empty", http.StatusBadRequest)
		return
	}

	if input.Email != nil && *input.Email == "" {
		http.Error(w, "email cannot be empty", http.StatusBadRequest)
		return
	}

	if input.Email != nil {
		if _, err := mail.ParseAddress(*input.Email); err != nil {
			http.Error(w, "email is not a valid format", http.StatusBadRequest)
			return
		}
	}

	updatedUser, err := h.repo.UpdateUserByID(r.Context(), id, input)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "user not found",
			})
			return
		}

		log.Printf("failed to update user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to update user",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *Handler) DeleteUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")

	if err := h.repo.DeleteUserByID(r.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "user not found",
			})
			return
		}

		log.Printf("failed to delete user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to delete user",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
