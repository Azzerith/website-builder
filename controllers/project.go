package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"website-builder/models"
	ws "website-builder/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectController struct {
	db  *gorm.DB
	hub *ws.Hub
}

func NewProjectController(db *gorm.DB, hub *ws.Hub) *ProjectController {
	return &ProjectController{db: db, hub: hub}
}

func (pc *ProjectController) CreateProject(c *gin.Context) {
	var input struct {
		Name       string `json:"name" binding:"required"`
		TeamID     string `json:"team_id" binding:"required"`
		TemplateID string `json:"template_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	project := models.Project{
		ID:         uuid.New().String(),
		Name:       input.Name,
		TeamID:     input.TeamID,
		CreatedBy:  userID.(string),
		TemplateID: &input.TemplateID,
		Status:     models.Draft,
	}

	if err := pc.db.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	pc.notifyTeam(project.TeamID, "project_created", project)
	c.JSON(http.StatusCreated, project)
}

func (pc *ProjectController) GetProject(c *gin.Context) {
	projectID := c.Param("id")

	var project models.Project
	if err := pc.db.Preload("Pages").Preload("Pages.Elements").
		First(&project, "id = ?", projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if !pc.hasTeamAccess(c, project.TeamID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (pc *ProjectController) UpdateProject(c *gin.Context) {
	projectID := c.Param("id")

	var input struct {
		Name   string               `json:"name"`
		Status models.ProjectStatus `json:"status"`
		Domain string               `json:"domain"`
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

	if !pc.hasTeamAccess(c, project.TeamID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

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

	pc.notifyTeam(project.TeamID, "project_updated", project)
	c.JSON(http.StatusOK, project)
}

func (pc *ProjectController) DeleteProject(c *gin.Context) {
	projectID := c.Param("id")

	var project models.Project
	if err := pc.db.First(&project, "id = ?", projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if !pc.hasTeamAccess(c, project.TeamID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := pc.db.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	pc.notifyTeam(project.TeamID, "project_deleted", project.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

func (pc *ProjectController) GetProjects(c *gin.Context) {
	teamID := c.Query("team_id")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team_id parameter is required"})
		return
	}

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

	pc.hub.Broadcast <- jsonMessage
}
