package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const userIDKey contextKey = "userID"

func (h *Handler) RequireAuth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		tokenString, found := strings.CutPrefix(header, "Bearer ")
		if !found || strings.TrimSpace(tokenString) == "" {
			writeError(w, http.StatusUnauthorized, "missing or malformed Authorization header")
			return
		}

		userID, err := h.tokens.Parse(strings.TrimSpace(tokenString))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		if !h.store.UserExists(userID) {
			writeError(w, http.StatusUnauthorized, "user no longer exists")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	})
}

func userIDFrom(r *http.Request) string {
	id, _ := r.Context().Value(userIDKey).(string)
	return id
}

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
