package handlers

import (
	"context"
	"net/http"

	"mana-tomb/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// ... (RegisterUser function remains the same) ...
func RegisterUser(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser struct {
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		query := `
			INSERT INTO users (username, email, password_hash)
			VALUES ($1, $2, $3)
			RETURNING id, username, email, created_at, updated_at
		`

		var createdUser models.User
		err = dbpool.QueryRow(context.Background(), query, newUser.Username, newUser.Email, string(hashedPassword)).Scan(
			&createdUser.ID,
			&createdUser.Username,
			&createdUser.Email,
			&createdUser.CreatedAt,
			&createdUser.UpdatedAt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createdUser)
	}
}

// ... (LoginUser function remains the same) ...
func LoginUser(dbpool *pgxpool.Pool, store *sessions.CookieStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginDetails struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginDetails); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		var user models.User
		query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
		err := dbpool.QueryRow(context.Background(), query, loginDetails.Email).Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginDetails.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		session, _ := store.Get(c.Request, "mana-tomb-session")
		session.Values["user_id"] = user.ID.String()
		err = session.Save(c.Request, c.Writer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// GetCurrentUser fetches the details of the currently logged-in user.
func GetCurrentUser(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the user ID from the context (set by the AuthRequired middleware).
		userIDStr, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in session"})
			return
		}

		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}

		// Fetch user details from the database.
		var user models.User
		query := `SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1`
		err = dbpool.QueryRow(context.Background(), query, userID).Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// LogoutUser handles logging out by clearing the session cookie.
func LogoutUser(store *sessions.CookieStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := store.Get(c.Request, "mana-tomb-session")

		// Setting MaxAge to -1 immediately deletes the cookie.
		session.Options.MaxAge = -1
		err := session.Save(c.Request, c.Writer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	}
}
