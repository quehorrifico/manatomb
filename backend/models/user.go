package models

import (
	"time"

	"github.com/google/uuid"
)

// User defines the structure for a user in the application.
// Note the `json:"-"` tag on PasswordHash to prevent it from being sent in API responses.
type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose this field
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
