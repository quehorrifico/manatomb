package models

import (
	"time"

	"github.com/google/uuid"
)

// The Deck struct is updated to hold separate slices for different boards.
// This makes it easier to send structured data to the frontend.
type Deck struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Format      string    `json:"format"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Mainboard   []Card    `json:"mainboard,omitempty"`
	Maybeboard  []Card    `json:"maybeboard,omitempty"`
}
