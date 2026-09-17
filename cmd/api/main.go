package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mini-crud/internal/user"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", homeHandler)
	mux.HandleFunc("GET /about", aboutHandler)
	mux.HandleFunc("GET /health", healthHandler)
	userRoutings(mux)

	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func userRoutings(mux *http.ServeMux) {
	mux.HandleFunc("GET /users", user.GetUserListHandler)
	mux.HandleFunc("GET /users/{id}", user.GetUserByIDHandler)
	mux.HandleFunc("POST /users", user.CreateUserHandler)
	mux.HandleFunc("PATCH /users/{id}", user.UpdateUserByIDHandler)
	mux.HandleFunc("DELETE /users/{id}", user.DeleteUserByIDHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
	}
	fmt.Fprintf(w, "Hello from Golang!")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "About page")
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
