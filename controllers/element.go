package controllers

import (
	"net/http"

	"website-builder/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ElementController struct {
	db *gorm.DB
}

func NewElementController(db *gorm.DB) *ElementController {
	return &ElementController{db: db}
}

func (ec *ElementController) CreateElement(c *gin.Context) {
	var input struct {
		PageID          string             `json:"page_id" binding:"required"`
		Type            models.ElementType `json:"type" binding:"required"`
		Data            models.JSON        `json:"data" binding:"required"`
		PositionX       int                `json:"position_x" binding:"required"`
		PositionY       int                `json:"position_y" binding:"required"`
		Width           int                `json:"width" binding:"required"`
		Height          int                `json:"height" binding:"required"`
		ZIndex          int                `json:"z_index"`
		ParentElementID *string            `json:"parent_element_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	element := models.Element{
		ID:              uuid.New().String(),
		PageID:          input.PageID,
		Type:            input.Type,
		Data:            input.Data,
		PositionX:       input.PositionX,
		PositionY:       input.PositionY,
		Width:           input.Width,
		Height:          input.Height,
		ZIndex:          input.ZIndex,
		ParentElementID: input.ParentElementID,
	}

	if err := ec.db.Create(&element).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create element"})
		return
	}

	c.JSON(http.StatusCreated, element)
}

func (ec *ElementController) GetElement(c *gin.Context) {
	elementID := c.Param("id")

	var element models.Element
	if err := ec.db.Preload("Child").Preload("Comment").
		First(&element, "id = ?", elementID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Element not found"})
		return
	}

	c.JSON(http.StatusOK, element)
}

func (ec *ElementController) UpdateElement(c *gin.Context) {
	elementID := c.Param("id")

	var input struct {
		Data      models.JSON `json:"data"`
		PositionX *int        `json:"position_x"`
		PositionY *int        `json:"position_y"`
		Width     *int        `json:"width"`
		Height    *int        `json:"height"`
		ZIndex    *int        `json:"z_index"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var element models.Element
	if err := ec.db.First(&element, "id = ?", elementID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Element not found"})
		return
	}

	// Update only provided fields
	if input.Data != nil {
		element.Data = input.Data
	}
	if input.PositionX != nil {
		element.PositionX = *input.PositionX
	}
	if input.PositionY != nil {
		element.PositionY = *input.PositionY
	}
	if input.Width != nil {
		element.Width = *input.Width
	}
	if input.Height != nil {
		element.Height = *input.Height
	}
	if input.ZIndex != nil {
		element.ZIndex = *input.ZIndex
	}

	if err := ec.db.Save(&element).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update element"})
		return
	}

	c.JSON(http.StatusOK, element)
}

// DeleteElement deletes an element
func (ec *ElementController) DeleteElement(c *gin.Context) {
	elementID := c.Param("id")

	if err := ec.db.Delete(&models.Element{}, "id = ?", elementID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete element"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Element deleted successfully"})
}

// ListElements lists all elements for a page
func (ec *ElementController) ListElements(c *gin.Context) {
	pageID := c.Query("page_id")
	if pageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page_id query parameter is required"})
		return
	}

	var elements []models.Element
	if err := ec.db.Where("page_id = ?", pageID).Find(&elements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch elements"})
		return
	}

	c.JSON(http.StatusOK, elements)
}
