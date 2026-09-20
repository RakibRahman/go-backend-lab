package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mini-crud/internal/platform/database"
	"mini-crud/internal/user"
	"net/http"
	"os"
)

const apiVersion = " /api/v1"

func main() {
	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")

	db, err := database.NewPostgresPool(ctx, databaseURL)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("connected to PostgreSQL")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", homeHandler)
	mux.HandleFunc("GET /about", aboutHandler)
	mux.HandleFunc("GET /health", healthHandler)
	userRoutes(mux)

	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func userRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET "+apiVersion+"/users", user.GetUserListHandler)
	mux.HandleFunc("GET "+apiVersion+"/users/{id}", user.GetUserByIDHandler)
	mux.HandleFunc("POST "+apiVersion+"/users", user.CreateUserHandler)
	mux.HandleFunc("PATCH "+apiVersion+"/users/{id}", user.UpdateUserByIDHandler)
	mux.HandleFunc("DELETE "+apiVersion+"/users/{id}", user.DeleteUserByIDHandler)
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
