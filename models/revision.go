package models

import (
	"time"
	"gorm.io/gorm"
)

type Revision struct {
	gorm.Model
	ID        string `gorm:"primaryKey;type:char(36)"`
	ProjectID string `gorm:"not null;type:char(36)"`
	UserID    string `gorm:"not null;type:char(36)"`
	Data      JSON   `gorm:"type:json;not null"`
	Message   string `gorm:"size:255"`
	CreatedAt time.Time
	Project   Project `gorm:"foreignKey:ProjectID"`
	User      User    `gorm:"foreignKey:UserID"`
}

func (Revision) TableName() string {
	return "revision"
}