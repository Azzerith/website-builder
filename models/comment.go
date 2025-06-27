package models

import (
	"time"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID        string `gorm:"primaryKey;type:char(36)"`
	ElementID string `gorm:"not null;type:char(36)"`
	UserID    string `gorm:"not null;type:char(36)"`
	Content   string `gorm:"not null;type:text"`
	Resolved  bool   `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Element   Element        `gorm:"foreignKey:ElementID"`
	User      User           `gorm:"foreignKey:UserID"`
	Reply     []CommentReply `gorm:"foreignKey:CommentID"`
}

func (Comment) TableName() string {
	return "comment"
}

type CommentReply struct {
	gorm.Model
	ID        string `gorm:"primaryKey;type:char(36)"`
	CommentID string `gorm:"not null;type:char(36)"`
	UserID    string `gorm:"not null;type:char(36)"`
	Content   string `gorm:"not null;type:text"`
	CreatedAt time.Time
	Comment   Comment `gorm:"foreignKey:CommentID"`
	User      User    `gorm:"foreignKey:UserID"`
}

func (CommentReply) TableName() string {
	return "comment_reply"
}