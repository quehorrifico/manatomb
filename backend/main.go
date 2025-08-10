package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"mana-tomb/backend/handlers"
	"mana-tomb/backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var store *sessions.CookieStore

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables from OS")
	}

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatal("SESSION_SECRET environment variable is required")
	}
	store = sessions.NewCookieStore([]byte(sessionSecret))

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	router := gin.Default()

	public := router.Group("/api/users")
	{
		public.POST("/register", handlers.RegisterUser(dbpool))
		public.POST("/login", handlers.LoginUser(dbpool, store))
	}

	// --- Public Routes ---
	// Profile routes are public so anyone can view them.
	profiles := router.Group("/api/profiles")
	{
		profiles.GET("/:username", handlers.GetUserProfile(dbpool))
	}

	protected := router.Group("/api")
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
			// Card-related routes
			decks.POST("/:deckId/cards", handlers.AddCardToDeck(dbpool))
			decks.DELETE("/:deckId/cards/:cardId", handlers.RemoveCardFromDeck(dbpool))
			decks.PUT("/:deckId/visibility", handlers.SetDeckVisibility(dbpool))
		}
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Backend is running!",
		})
	})

	log.Println("Starting server on :8080")
	router.Run(":8080")
}
