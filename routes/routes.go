package routes

import (
	"website-builder/controllers"
	"website-builder/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, hub *websocket.Hub) {
	api := r.Group("/api")
	{
		// Project routes
		api.GET("/projects", controllers.GetProjects)
		api.POST("/projects", controllers.CreateProject)
		api.GET("/projects/:id", controllers.GetProject)
		api.PUT("/projects/:id", controllers.UpdateProject)
		api.DELETE("/projects/:id", controllers.DeleteProject)

		// WebSocket route
		api.GET("/ws", func(c *gin.Context) {
			controllers.WebSocketHandler(c, hub)
		})
	}

	// Static files for production
	r.Static("/static", "./static")
}