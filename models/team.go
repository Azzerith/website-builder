package models

import (
	"time"
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	ID        string `gorm:"primaryKey;type:char(36)"`
	Name      string `gorm:"not null;size:100"`
	CreatedBy string `gorm:"not null;type:char(36)"`
	CreatedAt time.Time
	TeamMember []TeamMember `gorm:"foreignKey:TeamID"`
	Project    []Project    `gorm:"foreignKey:TeamID"`
}

func (Team) TableName() string {
	return "team"
}

type TeamMemberRole string

const (
	Owner  TeamMemberRole = "owner"
	Admin  TeamMemberRole = "admin"
	Editor TeamMemberRole = "editor"
	Viewer TeamMemberRole = "viewer"
)

type TeamMember struct {
	gorm.Model
	TeamID    string        `gorm:"primaryKey;type:char(36)"`
	UserID    string        `gorm:"primaryKey;type:char(36)"`
	Role      TeamMemberRole `gorm:"type:enum('owner','admin','editor','viewer');default:'editor'"`
	JoinedAt  time.Time
	Team      Team `gorm:"foreignKey:TeamID"`
	User      User `gorm:"foreignKey:UserID"`
}

func (TeamMember) TableName() string {
	return "team_member"
}