package main

import (
	"log"
	"os"
	"strings"

	"website-builder/config"
	"website-builder/routes"
	"website-builder/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env (for development)
	// In production, these should be set in the environment
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using environment variables")
	}

	// Initialize database
	config.InitDB()

	// Set up Gin router
	r := gin.Default()

	// Get allowed origins from environment
	allowedOrigins := getOriginsFromEnv()

	// CORS configuration with WebSocket support
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Sec-WebSocket-Protocol"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Pass hub to routes
	routes.SetupRoutes(r, config.DB, hub)


	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run server
	log.Printf("Server running on port %s", port)
	log.Printf("Allowed origins: %v", allowedOrigins)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getOriginsFromEnv() []string {
	// Default origins for development
	defaultOrigins := []string{
		"http://localhost:5173", 
		"http://localhost:5174",
		"http://127.0.0.1:5173",
		"ws://localhost:5173",
	}

	// Get additional origins from environment
	envOrigins := os.Getenv("ALLOWED_ORIGINS")
	if envOrigins == "" {
		return defaultOrigins
	}

	// Split multiple origins separated by comma
	var origins []string
	for _, origin := range strings.Split(envOrigins, ",") {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			origins = append(origins, trimmed)
			// Also add WebSocket version if it's http
			if strings.HasPrefix(trimmed, "http://") {
				origins = append(origins, strings.Replace(trimmed, "http://", "ws://", 1))
			} else if strings.HasPrefix(trimmed, "https://") {
				origins = append(origins, strings.Replace(trimmed, "https://", "wss://", 1))
			}
		}
	}

	return append(defaultOrigins, origins...)
}