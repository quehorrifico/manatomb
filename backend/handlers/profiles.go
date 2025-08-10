package handlers

import (
	"context"
	"net/http"

	"mana-tomb/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetUserProfile fetches a user's public profile by their username.
func GetUserProfile(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		// 1. Find the user by username to get their ID.
		var userID string
		var user models.User
		userQuery := `SELECT id, username FROM users WHERE username = $1`
		err := dbpool.QueryRow(context.Background(), userQuery, username).Scan(&userID, &user.Username)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// 2. Fetch all decks that belong to this user AND are marked as public.
		decksQuery := `
			SELECT id, name, description, format, user_id, created_at, updated_at 
			FROM decks 
			WHERE user_id = $1 AND is_public = TRUE 
			ORDER BY updated_at DESC
		`
		rows, err := dbpool.Query(context.Background(), decksQuery, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve public decks"})
			return
		}
		defer rows.Close()

		publicDecks := make([]models.Deck, 0)
		for rows.Next() {
			var deck models.Deck
			if err := rows.Scan(&deck.ID, &deck.Name, &deck.Description, &deck.Format, &deck.UserID, &deck.CreatedAt, &deck.UpdatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan deck row"})
				return
			}
			publicDecks = append(publicDecks, deck)
		}

		profile := models.Profile{
			Username:    user.Username,
			PublicDecks: publicDecks,
		}

		c.JSON(http.StatusOK, profile)
	}
}
