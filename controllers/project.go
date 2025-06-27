package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"website-builder/models"
	"website-builder/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectController struct {
	db  *gorm.DB
	hub *websocket.Hub
}

func NewProjectController(db *gorm.DB, hub *websocket.Hub) *ProjectController {
	return &ProjectController{db: db, hub: hub}
}

// CreateProject handles project creation
func (pc *ProjectController) CreateProject(c *gin.Context) {
	var input struct {
		Name      string `json:"name" binding:"required"`
		TeamID    string `json:"team_id" binding:"required"`
		TemplateID string `json:"template_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	project := models.Project{
		ID:        uuid.New().String(),
		Name:      input.Name,
		TeamID:    input.TeamID,
		CreatedBy: userID.(string),
		TemplateID: &input.TemplateID,
		Status:    models.Draft,
	}

	if err := pc.db.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	// Broadcast project creation to team members
	pc.notifyTeam(project.TeamID, "project_created", project)

	c.JSON(http.StatusCreated, project)
}

// GetProject retrieves a single project
func (pc *ProjectController) GetProject(c *gin.Context) {
	projectID := c.Param("id")

	var project models.Project
	if err := pc.db.Preload("Pages").Preload("Pages.Elements").
		First(&project, "id = ?", projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Check if user has access to this project
	if !pc.hasProjectAccess(c, project.TeamID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// UpdateProject handles project updates
func (pc *ProjectController) UpdateProject(c *gin.Context) {
	projectID := c.Param("id")

	var input struct {
		Name   string           `json:"name"`
		Status models.ProjectStatus `json:"status"`
		Domain string           `json:"domain"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project models.Project
	if err := pc.db.First(&project, "id = ?", projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Check if user has access to this project
	if !pc.hasProjectAccess(c, project.TeamID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update fields
	if input.Name != "" {
		project.Name = input.Name
	}
	if input.Status != "" {
		project.Status = input.Status
	}
	if input.Domain != "" {
		project.Domain = input.Domain
	}

	if err := pc.db.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	// Broadcast update to team members
	pc.notifyTeam(project.TeamID, "project_updated", project)

	c.JSON(http.StatusOK, project)
}

// DeleteProject handles project deletion
func (pc *ProjectController) DeleteProject(c *gin.Context) {
	projectID := c.Param("id")

	var project models.Project
	if err := pc.db.First(&project, "id = ?", projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Check if user has access to this project
	if !pc.hasProjectAccess(c, project.TeamID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := pc.db.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	// Broadcast deletion to team members
	pc.notifyTeam(project.TeamID, "project_deleted", project.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// GetProjects lists all projects for a team
func (pc *ProjectController) GetProjects(c *gin.Context) {
	teamID := c.Query("team_id")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team_id parameter is required"})
		return
	}

	// Check if user has access to this team
	if !pc.hasTeamAccess(c, teamID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var projects []models.Project
	if err := pc.db.Where("team_id = ?", teamID).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// Helper function to check project access
func (pc *ProjectController) hasProjectAccess(c *gin.Context, teamID string) bool {
	return pc.hasTeamAccess(c, teamID)
}

// Helper function to check team access
func (pc *ProjectController) hasTeamAccess(c *gin.Context, teamID string) bool {
	userID, exists := c.Get("userID")
	if !exists {
		return false
	}

	var count int64
	pc.db.Model(&models.TeamMember{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Count(&count)

	return count > 0
}

// Helper function to notify team members
func (pc *ProjectController) notifyTeam(teamID string, event string, data interface{}) {
	message := map[string]interface{}{
		"event": event,
		"data":  data,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal notification message: %v", err)
		return
	}

	// Broadcast to all connected clients in the team
	// In a real implementation, you would need to track which clients belong to which team
	pc.hub.Broadcast(jsonMessage)
}