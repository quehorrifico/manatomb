package handlers

import (
	"context"
	"net/http"

	"mana-tomb/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ... (CreateDeck, GetUserDecks, UpdateDeck, DeleteDeck functions are unchanged) ...
func CreateDeck(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newDeckData struct {
			Name        string `json:"name" binding:"required"`
			Description string `json:"description"`
			Format      string `json:"format"`
		}

		if err := c.ShouldBindJSON(&newDeckData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		userIDStr, _ := c.Get("userID")
		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			return
		}

		query := `
			INSERT INTO decks (name, description, format, user_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id, name, description, format, user_id, created_at, updated_at
		`
		var createdDeck models.Deck
		err = dbpool.QueryRow(context.Background(), query, newDeckData.Name, newDeckData.Description, newDeckData.Format, userID).Scan(
			&createdDeck.ID,
			&createdDeck.Name,
			&createdDeck.Description,
			&createdDeck.Format,
			&createdDeck.UserID,
			&createdDeck.CreatedAt,
			&createdDeck.UpdatedAt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deck: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createdDeck)
	}
}

func GetUserDecks(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, _ := c.Get("userID")
		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			return
		}

		query := `SELECT id, name, description, format, user_id, created_at, updated_at FROM decks WHERE user_id = $1 ORDER BY updated_at DESC`
		rows, err := dbpool.Query(context.Background(), query, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve decks"})
			return
		}
		defer rows.Close()

		decks := make([]models.Deck, 0)
		for rows.Next() {
			var deck models.Deck
			if err := rows.Scan(&deck.ID, &deck.Name, &deck.Description, &deck.Format, &deck.UserID, &deck.CreatedAt, &deck.UpdatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan deck row"})
				return
			}
			decks = append(decks, deck)
		}

		c.JSON(http.StatusOK, decks)
	}
}

func UpdateDeck(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckIDStr := c.Param("deckId")
		deckID, err := uuid.Parse(deckIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID format"})
			return
		}

		var deckData struct {
			Name        string `json:"name" binding:"required"`
			Description string `json:"description"`
		}
		if err := c.ShouldBindJSON(&deckData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		query := `
			UPDATE decks
			SET name = $1, description = $2, updated_at = NOW()
			WHERE id = $3
			RETURNING id, name, description, format, user_id, created_at, updated_at
		`
		var updatedDeck models.Deck
		err = dbpool.QueryRow(context.Background(), query, deckData.Name, deckData.Description, deckID).Scan(
			&updatedDeck.ID, &updatedDeck.Name, &updatedDeck.Description, &updatedDeck.Format,
			&updatedDeck.UserID, &updatedDeck.CreatedAt, &updatedDeck.UpdatedAt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update deck: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, updatedDeck)
	}
}

func DeleteDeck(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckIDStr := c.Param("deckId")
		deckID, err := uuid.Parse(deckIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID format"})
			return
		}

		query := `DELETE FROM decks WHERE id = $1`
		cmdTag, err := dbpool.Exec(context.Background(), query, deckID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete deck"})
			return
		}

		if cmdTag.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found or not owned by user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Deck deleted successfully"})
	}
}

// GetDeckByID is updated to fetch cards and group them by board.
func GetDeckByID(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckIDStr := c.Param("deckId")
		deckID, err := uuid.Parse(deckIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID format"})
			return
		}

		var deck models.Deck
		deckQuery := `SELECT id, name, description, format, user_id, created_at, updated_at FROM decks WHERE id = $1`
		err = dbpool.QueryRow(context.Background(), deckQuery, deckID).Scan(
			&deck.ID, &deck.Name, &deck.Description, &deck.Format, &deck.UserID, &deck.CreatedAt, &deck.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
			return
		}

		cardsQuery := `
			SELECT c.scryfall_id, c.name, c.image_uris, c.mana_cost, c.cmc, c.type_line, c.oracle_text, c.colors, c.color_identity, dc.quantity, dc.board
			FROM cards c
			JOIN deck_cards dc ON c.scryfall_id = dc.card_scryfall_id
			WHERE dc.deck_id = $1
		`
		rows, err := dbpool.Query(context.Background(), cardsQuery, deckID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cards for deck"})
			return
		}
		defer rows.Close()

		deck.Mainboard = make([]models.Card, 0)
		deck.Maybeboard = make([]models.Card, 0)

		for rows.Next() {
			var card models.Card
			var board string
			if err := rows.Scan(&card.ScryfallID, &card.Name, &card.ImageURIs, &card.ManaCost, &card.CMC, &card.TypeLine, &card.OracleText, &card.Colors, &card.ColorIdentity, &card.Quantity, &board); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan card row"})
				return
			}
			// Sort cards into the correct slice based on the board.
			if board == "maybeboard" {
				deck.Maybeboard = append(deck.Maybeboard, card)
			} else {
				deck.Mainboard = append(deck.Mainboard, card)
			}
		}

		c.JSON(http.StatusOK, deck)
	}
}

func SetDeckVisibility(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckIDStr := c.Param("deckId")
		deckID, err := uuid.Parse(deckIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID format"})
			return
		}

		var payload struct {
			IsPublic bool `json:"is_public"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		userIDStr, _ := c.Get("userID")
		// Important: Ensure the user owns the deck they are trying to modify.
		query := `
			UPDATE decks 
			SET is_public = $1 
			WHERE id = $2 AND user_id = $3
		`
		cmdTag, err := dbpool.Exec(context.Background(), query, payload.IsPublic, deckID, userIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update deck visibility"})
			return
		}

		if cmdTag.RowsAffected() == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Deck not found or you do not have permission to edit it"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Deck visibility updated successfully"})
	}
}
