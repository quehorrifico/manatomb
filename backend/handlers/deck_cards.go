package handlers

import (
	"context"
	"net/http"

	"mana-tomb/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AddCardToDeck now accepts a 'board' parameter in the request body.
func AddCardToDeck(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckIDStr := c.Param("deckId")
		deckID, err := uuid.Parse(deckIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID format"})
			return
		}

		var requestBody struct {
			Card  models.Card `json:"card"`
			Board string      `json:"board"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
			return
		}

		card := requestBody.Card
		board := requestBody.Board
		if board == "" {
			board = "main" // Default to main board
		}

		tx, err := dbpool.Begin(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
			return
		}
		defer tx.Rollback(context.Background())

		cardCacheQuery := `
			INSERT INTO cards (scryfall_id, name, image_uris, mana_cost, cmc, type_line, oracle_text, colors, color_identity)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (scryfall_id) DO NOTHING
		`
		_, err = tx.Exec(context.Background(), cardCacheQuery,
			card.ScryfallID, card.Name, card.ImageURIs, card.ManaCost, card.CMC, card.TypeLine, card.OracleText, card.Colors, card.ColorIdentity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache card data"})
			return
		}

		addCardQuery := `
			INSERT INTO deck_cards (deck_id, card_scryfall_id, board, quantity)
			VALUES ($1, $2, $3, 1)
			ON CONFLICT (deck_id, card_scryfall_id, board) DO UPDATE
			SET quantity = deck_cards.quantity + 1
		`
		// We now include the board in the conflict target and insert.
		_, err = tx.Exec(context.Background(), addCardQuery, deckID, card.ScryfallID, board)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add card to deck"})
			return
		}

		if err := tx.Commit(context.Background()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Card added successfully"})
	}
}

// RemoveCardFromDeck handles decrementing a card's quantity or removing it entirely.
func RemoveCardFromDeck(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckIDStr := c.Param("deckId")
		cardIDStr := c.Param("cardId")
		deckID, err := uuid.Parse(deckIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
			return
		}
		cardID, err := uuid.Parse(cardIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
			return
		}

		// First, check the current quantity.
		var quantity int
		checkQuery := `SELECT quantity FROM deck_cards WHERE deck_id = $1 AND card_scryfall_id = $2`
		err = dbpool.QueryRow(context.Background(), checkQuery, deckID, cardID).Scan(&quantity)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Card not found in deck"})
			return
		}

		if quantity > 1 {
			// If more than one, decrement the quantity.
			updateQuery := `UPDATE deck_cards SET quantity = quantity - 1 WHERE deck_id = $1 AND card_scryfall_id = $2`
			_, err = dbpool.Exec(context.Background(), updateQuery, deckID, cardID)
		} else {
			// If only one, delete the row.
			deleteQuery := `DELETE FROM deck_cards WHERE deck_id = $1 AND card_scryfall_id = $2`
			_, err = dbpool.Exec(context.Background(), deleteQuery, deckID, cardID)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove card"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Card removed successfully"})
	}
}
