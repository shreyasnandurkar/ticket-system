package models

import (
	"errors"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in_progress"
	StatusClosed     Status = "closed"
)

type Ticket struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var (
	ErrTicketClosed      = errors.New("closed ticket cannot be updated")
	ErrInvalidTransition = errors.New("ticket status cannot move backwards")
)

// open -> in_progress -> closed
var order = map[Status]int{
	StatusOpen:       0,
	StatusInProgress: 1,
	StatusClosed:     2,
}

func (s Status) IsValid() bool {
	_, ok := order[s]
	return ok
}

// CheckTransition error is returned if a ticket cannot move from "from" to "to".
func CheckTransition(from, to Status) error {
	if from == StatusClosed {
		return ErrTicketClosed
	}
	if order[to] < order[from] {
		return ErrInvalidTransition
	}
	return nil
}
