package models

import (
	"time"
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model
	ID            string `gorm:"primaryKey;type:char(36)"`
	ProjectID     string `gorm:"not null;type:char(36)"`
	Name          string `gorm:"not null;size:100"`
	Path          string `gorm:"not null;size:100"`
	IsHomepage    bool   `gorm:"default:false"`
	SEOTitle      string `gorm:"size:255"`
	SEODescription string `gorm:"type:text"`
	SEOKeywords   string `gorm:"size:255"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Project       Project   `gorm:"foreignKey:ProjectID"`
	Element       []Element `gorm:"foreignKey:PageID"`
}

func (Page) TableName() string {
	return "page"
}