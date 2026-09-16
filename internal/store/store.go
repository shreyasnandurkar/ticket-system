package store

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"ticket-system/internal/models"
)

var (
	ErrEmailTaken = errors.New("email already registered")
	ErrNotFound   = errors.New("not found")
)

// Store - I preferred to use in-memory storage.
// The mutex makes it safe to use from many requests at the same time.
type Store struct {
	mu           sync.RWMutex
	usersByEmail map[string]models.User
	usersByID    map[string]models.User
	tickets      map[string]models.Ticket
}

func New() *Store {
	return &Store{
		usersByEmail: make(map[string]models.User),
		usersByID:    make(map[string]models.User),
		tickets:      make(map[string]models.Ticket),
	}
}

func (s *Store) CreateUser(email, passwordHash string) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.usersByEmail[email]; exists {
		return models.User{}, ErrEmailTaken
	}

	user := models.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
	}
	s.usersByEmail[email] = user
	s.usersByID[user.ID] = user
	return user, nil
}

func (s *Store) GetUserByEmail(email string) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.usersByEmail[email]
	if !ok {
		return models.User{}, ErrNotFound
	}
	return user, nil
}

// UserExists is used by the auth middleware to reject tokens
// for users that are no longer in the store (e.g. after a restart).
func (s *Store) UserExists(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.usersByID[id]
	return ok
}

func (s *Store) CreateTicket(userID, title, description string) models.Ticket {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	ticket := models.Ticket{
		ID:          uuid.NewString(),
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      models.StatusOpen,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tickets[ticket.ID] = ticket
	return ticket
}

// ListTickets returns only the tickets owned by userID, oldest first.
func (s *Store) ListTickets(userID string) []models.Ticket {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Ticket
	for _, t := range s.tickets {
		if t.UserID == userID {
			result = append(result, t)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result
}

func (s *Store) GetTicket(id, userID string) (models.Ticket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tickets[id]
	if !ok || t.UserID != userID {
		return models.Ticket{}, ErrNotFound
	}
	return t, nil
}

// UpdateStatus does the ownership check, the transition check, and the update.
func (s *Store) UpdateStatus(id, userID string, status models.Status) (models.Ticket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tickets[id]
	if !ok || t.UserID != userID {
		return models.Ticket{}, ErrNotFound
	}
	if err := models.CheckTransition(t.Status, status); err != nil {
		return models.Ticket{}, err
	}

	t.Status = status
	t.UpdatedAt = time.Now().UTC()
	s.tickets[id] = t
	return t, nil
}
