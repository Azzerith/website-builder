package models

import (
	"time"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	ID         string `gorm:"primaryKey;type:char(36)"`
	ProjectID  string `gorm:"not null;type:char(36)"`
	UserID     string `gorm:"not null;type:char(36)"`
	SocketID   string `gorm:"not null;size:255"`
	LastActive time.Time
	Project    Project `gorm:"foreignKey:ProjectID"`
	User       User    `gorm:"foreignKey:UserID"`
}

func (Session) TableName() string {
	return "session"
}