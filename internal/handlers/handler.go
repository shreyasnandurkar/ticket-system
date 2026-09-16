package handlers

import (
	"ticket-system/internal/auth"
	"ticket-system/internal/store"
)

// Handler holds the dependencies every HTTP handler needs.
type Handler struct {
	store  *store.Store
	tokens *auth.TokenManager
}

func New(s *store.Store, t *auth.TokenManager) *Handler {
	return &Handler{store: s, tokens: t}
}
