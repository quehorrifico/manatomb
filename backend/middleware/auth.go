package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// AuthRequired is a middleware to ensure a user is authenticated.
func AuthRequired(store *sessions.CookieStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, "mana-tomb-session")
		if err != nil {
			// This could happen if the cookie secret changes, for example.
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			return
		}

		// Check if the user_id exists in the session.
		// If it's not there, the user is not logged in.
		if userID, ok := session.Values["user_id"].(string); !ok || userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		// If the user_id exists, we can add it to the request context
		// so that subsequent handlers can access it.
		c.Set("userID", session.Values["user_id"])

		// Continue to the next handler.
		c.Next()
	}
}
