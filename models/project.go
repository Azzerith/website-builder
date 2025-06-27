package models

import (
	"time"
	"gorm.io/gorm"
)

type ProjectStatus string

const (
	Draft     ProjectStatus = "draft"
	Published ProjectStatus = "published"
	Archived  ProjectStatus = "archived"
)

type Project struct {
	gorm.Model
	ID          string        `gorm:"primaryKey;type:char(36)"`
	Name        string        `gorm:"not null;size:100"`
	TeamID      string        `gorm:"not null;type:char(36)"`
	CreatedBy   string        `gorm:"not null;type:char(36)"`
	TemplateID  *string       `gorm:"type:char(36)"`
	Domain      string        `gorm:"size:255"`
	PublishedURL string       `gorm:"size:255"`
	Status      ProjectStatus `gorm:"type:enum('draft','published','archived');default:'draft'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Team        Team         `gorm:"foreignKey:TeamID"`
	Creator     User         `gorm:"foreignKey:CreatedBy"`
	Page        []Page       `gorm:"foreignKey:ProjectID"`
	Revision    []Revision   `gorm:"foreignKey:ProjectID"`
	Session     []Session    `gorm:"foreignKey:ProjectID"`
}

func (Project) TableName() string {
	return "project"
}