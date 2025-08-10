package models

// Profile represents the publicly viewable information for a user.
// It includes the username and a list of their public decks.
type Profile struct {
	Username    string `json:"username"`
	PublicDecks []Deck `json:"public_decks"`
}
