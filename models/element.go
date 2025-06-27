package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
	"gorm.io/gorm"
)

type ElementType string

const (
	TextElement    ElementType = "text"
	ImageElement   ElementType = "image"
	ButtonElement  ElementType = "button"
	VideoElement   ElementType = "video"
	FormElement    ElementType = "form"
	SectionElement ElementType = "section"
	DividerElement ElementType = "divider"
	MapElement     ElementType = "map"
	SocialElement  ElementType = "social"
)

type Element struct {
	gorm.Model
	ID              string      `gorm:"primaryKey;type:char(36)"`
	PageID          string      `gorm:"not null;type:char(36)"`
	Type            ElementType `gorm:"type:enum('text','image','button','video','form','section','divider','map','social');not null"`
	Data            JSON        `gorm:"type:json;not null"`
	PositionX       int         `gorm:"not null"`
	PositionY       int         `gorm:"not null"`
	Width           int         `gorm:"not null"`
	Height          int         `gorm:"not null"`
	ZIndex          int         `gorm:"default:0"`
	ParentElementID *string     `gorm:"type:char(36)"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Page            Page       `gorm:"foreignKey:PageID"`
	Parent          *Element   `gorm:"foreignKey:ParentElementID"`
	Child           []Element  `gorm:"foreignKey:ParentElementID"`
	Comment         []Comment  `gorm:"foreignKey:ElementID"`
}

func (Element) TableName() string {
	return "element"
}

type JSON map[string]interface{}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, j)
}

func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}