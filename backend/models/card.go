package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Card represents the data for a single Magic: The Gathering card
// that we cache in our database from Scryfall.
type Card struct {
	ScryfallID    uuid.UUID       `json:"id"` // Note: This is the Scryfall ID
	Name          string          `json:"name"`
	ImageURIs     json.RawMessage `json:"image_uris"`
	ManaCost      string          `json:"mana_cost"`
	CMC           float32         `json:"cmc"`
	TypeLine      string          `json:"type_line"`
	OracleText    string          `json:"oracle_text"`
	Colors        []string        `json:"colors"`
	ColorIdentity []string        `json:"color_identity"`
	Quantity      int             `json:"quantity,omitempty"` // Used when returning cards in a deck
}
