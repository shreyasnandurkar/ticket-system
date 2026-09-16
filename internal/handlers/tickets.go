package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"ticket-system/internal/models"
	"ticket-system/internal/store"
)

type createTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type updateStatusRequest struct {
	Status models.Status `json:"status"`
}

func (h *Handler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	var req createTicketRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	ticket := h.store.CreateTicket(userIDFrom(r), title, strings.TrimSpace(req.Description))
	writeJSON(w, http.StatusCreated, ticket)
}

func (h *Handler) ListTickets(w http.ResponseWriter, r *http.Request) {
	tickets := h.store.ListTickets(userIDFrom(r))
	writeJSON(w, http.StatusOK, tickets)
}

func (h *Handler) GetTicket(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	ticket, err := h.store.GetTicket(id, userIDFrom(r))
	if err != nil {
		writeError(w, http.StatusNotFound, "ticket not found")
		return
	}
	writeJSON(w, http.StatusOK, ticket)
}

func (h *Handler) UpdateTicketStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req updateStatusRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !req.Status.IsValid() {
		writeError(w, http.StatusBadRequest, "status must be one of: open, in_progress, closed")
		return
	}

	ticket, err := h.store.UpdateStatus(id, userIDFrom(r), req.Status)
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, "ticket not found")
	case errors.Is(err, models.ErrTicketClosed), errors.Is(err, models.ErrInvalidTransition):
		writeError(w, http.StatusConflict, err.Error())
	case err != nil:
		writeError(w, http.StatusInternalServerError, "could not update ticket")
	default:
		writeJSON(w, http.StatusOK, ticket)
	}
}

func parseID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, "ticket not found")
		return "", false
	}
	return id.String(), true
}
