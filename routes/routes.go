package routes

import (
	"website-builder/controllers"
	"website-builder/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, hub *websocket.Hub) {
	// Initialize controllers
	projectController := controllers.NewProjectController(db, hub)

	// API routes with auth middleware
	api := r.Group("/api", middleware.AuthMiddleware())
	{
		// Project routes
		api.POST("/projects", projectController.CreateProject)
		api.GET("/projects", projectController.GetProjects)
		api.GET("/projects/:id", projectController.GetProject)
		api.PUT("/projects/:id", projectController.UpdateProject)
		api.DELETE("/projects/:id", projectController.DeleteProject)

		// WebSocket route
		api.GET("/ws", func(c *gin.Context) {
			controllers.WebSocketHandler(c, hub)
		})
	}
}