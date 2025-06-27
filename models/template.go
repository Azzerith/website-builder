package models

import (
	"time"
	"gorm.io/gorm"
)

type Template struct {
	gorm.Model
	ID          string `gorm:"primaryKey;type:char(36)"`
	Name        string `gorm:"not null;size:100"`
	Data        JSON   `gorm:"type:json;not null"`
	Category    string `gorm:"not null;size:50"`
	ThumbnailURL string `gorm:"not null;size:255"`
	IsPremium   bool   `gorm:"default:false"`
	CreatedAt   time.Time
}

func (Template) TableName() string {
	return "template"
}