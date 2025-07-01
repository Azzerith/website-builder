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
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using environment variables")
	}

	config.InitDB()

	r := gin.Default()

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
	defaultOrigins := []string{
		"http://localhost:5173",
		"http://localhost:5174",
		"http://127.0.0.1:5173",
	}

	envOrigins := os.Getenv("ALLOWED_ORIGINS")
	if envOrigins == "" {
		return defaultOrigins
	}

	var origins []string
	for _, origin := range strings.Split(envOrigins, ",") {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" && (strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") || trimmed == "*") {
			origins = append(origins, trimmed)
		} else {
			log.Printf("Skipping invalid CORS origin: %s", trimmed)
		}
	}

	return append(defaultOrigins, origins...)
}
