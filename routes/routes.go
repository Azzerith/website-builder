package routes

import (
	"website-builder/controllers"
	"website-builder/middleware"
	"website-builder/websocket"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, hub *websocket.Hub) {
	// Initialize controllers
	authController := controllers.NewAuthController(db)
	projectController := controllers.NewProjectController(db, hub)
	elementController := controllers.NewElementController(db)

	// Public routes (no auth required)
	api := r.Group("/api")
	{
		api.POST("/login", authController.Login)
		api.POST("/register", authController.Register)
	}

	// Protected routes (require auth)
	protected := r.Group("/api", middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/me", authController.GetCurrentUser)

		// Project routes
		protected.POST("/projects", projectController.CreateProject)
		protected.GET("/projects", projectController.GetProjects)
		protected.GET("/projects/:id", projectController.GetProject)
		protected.PUT("/projects/:id", projectController.UpdateProject)
		protected.DELETE("/projects/:id", projectController.DeleteProject)

		protected.POST("/elements", elementController.CreateElement)
	protected.GET("/elements/:id", elementController.GetElement)
	protected.PUT("/elements/:id", elementController.UpdateElement)
	protected.DELETE("/elements/:id", elementController.DeleteElement)
	protected.GET("/elements", elementController.ListElements)

		// WebSocket route
		protected.GET("/ws", func(c *gin.Context) {
			controllers.WebSocketHandler(c, hub)
		})
	}
}
