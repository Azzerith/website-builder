package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID         string `gorm:"primaryKey;type:char(36)"`
	Email      string `gorm:"unique;not null;size:255"`
	Password   string `gorm:"not null;size:255"`
	FullName   string `gorm:"not null;size:100"`
	AvatarURL  string `gorm:"size:255"`
	LastLogin  *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	TeamMember []TeamMember `gorm:"foreignKey:UserID"`
	Project    []Project    `gorm:"foreignKey:CreatedBy"`
	Comment    []Comment    `gorm:"foreignKey:UserID"`
	Revision   []Revision   `gorm:"foreignKey:UserID"`
	Session    []Session    `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "user"
}