package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"ticket-system/internal/auth"
	"ticket-system/internal/handlers"
	"ticket-system/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
		log.Println("WARNING: JWT_SECRET is not set, using an insecure default")
	}

	st := store.New()
	tokens := auth.NewTokenManager(secret, 24*time.Hour)
	h := handlers.New(st, tokens)

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)

	// Protected routes (that mean the user must be authenticated)
	mux.Handle("POST /tickets", h.RequireAuth(h.CreateTicket))
	mux.Handle("GET /tickets", h.RequireAuth(h.ListTickets))
	mux.Handle("GET /tickets/{id}", h.RequireAuth(h.GetTicket))
	mux.Handle("PATCH /tickets/{id}/status", h.RequireAuth(h.UpdateTicketStatus))

	mux.HandleFunc("/health", h.MethodNotAllowed("GET"))
	mux.HandleFunc("/auth/register", h.MethodNotAllowed("POST"))
	mux.HandleFunc("/auth/login", h.MethodNotAllowed("POST"))
	mux.HandleFunc("/tickets", h.MethodNotAllowed("GET, POST"))
	mux.HandleFunc("/tickets/{id}", h.MethodNotAllowed("GET"))
	mux.HandleFunc("/tickets/{id}/status", h.MethodNotAllowed("PATCH"))
	mux.HandleFunc("/", h.NotFound)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handlers.LogRequests(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	log.Printf("server listening on :%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
