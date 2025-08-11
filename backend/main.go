package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"mana-tomb/backend/handlers"
	"mana-tomb/backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var store *sessions.CookieStore

func main() {
	// --- Load Environment Variables ---
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables from OS")
	}

	// --- Session Store Setup ---
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatal("SESSION_SECRET environment variable is required")
	}
	store = sessions.NewCookieStore([]byte(sessionSecret))

	// --- Database Connection Setup ---
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// Check for the production database certificate.
	dbCert := os.Getenv("DB_CERT")
	if dbCert != "" {
		// If the certificate exists, we're in production.
		// Write the certificate to a temporary file.
		certPath := filepath.Join(os.TempDir(), "db-ca-certificate.crt")
		err := os.WriteFile(certPath, []byte(dbCert), 0644)
		if err != nil {
			log.Fatalf("Unable to write database certificate to temp file: %v", err)
		}
		// Append SSL options to the connection string.
		connStr = fmt.Sprintf("%s sslmode=require sslrootcert=%s", connStr, certPath)
	} else {
		// Otherwise, we're in local development.
		connStr = fmt.Sprintf("%s sslmode=disable", connStr)
	}

	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	// --- Router Setup ---
	router := gin.Default()

	// Set Gin to release mode in production
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// API Routes (unchanged)
	api := router.Group("/api")
	{
		authPublic := api.Group("/users")
		{
			authPublic.POST("/register", handlers.RegisterUser(dbpool))
			authPublic.POST("/login", handlers.LoginUser(dbpool, store))
		}
		profiles := api.Group("/profiles")
		{
			profiles.GET("/:username", handlers.GetUserProfile(dbpool))
		}
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired(store))
		{
			protected.GET("/users/me", handlers.GetCurrentUser(dbpool))
			protected.POST("/users/logout", handlers.LogoutUser(store))
			decks := protected.Group("/decks")
			{
				decks.POST("/", handlers.CreateDeck(dbpool))
				decks.GET("/", handlers.GetUserDecks(dbpool))
				decks.GET("/:deckId", handlers.GetDeckByID(dbpool))
				decks.PUT("/:deckId", handlers.UpdateDeck(dbpool))
				decks.DELETE("/:deckId", handlers.DeleteDeck(dbpool))
				decks.PUT("/:deckId/visibility", handlers.SetDeckVisibility(dbpool))
				decks.POST("/:deckId/cards", handlers.AddCardToDeck(dbpool))
				decks.DELETE("/:deckId/cards/:cardId", handlers.RemoveCardFromDeck(dbpool))
			}
		}
	}

	log.Println("Starting server on :8080")
	router.Run(":8080")
}
