package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CourseStatus string

const (
	CourseDraft   CourseStatus = "draft"
	CoursePublish CourseStatus = "publish"
	CoursePreview CourseStatus = "preview"
)

type Investment struct {
	TotalInvestment float64 `gorm:"default:0" json:"total_investment"`
	Currency        string  `gorm:"default:'USDT'" json:"currency"`
	Investors       int     `gorm:"default:0" json:"investors"`
	Progress        float64 `gorm:"default:0" json:"progress"` // 0% -> 100%
	Goal            float64 `gorm:"default:0" json:"goal"`
}

type Course struct {
	gorm.Model
	Name             string         `gorm:"size:255;not null" json:"name"`
	ShortDescription string         `gorm:"size:255;not null" json:"short_description"`
	Description      string         `gorm:"type:text" json:"description"`
	Content          string         `gorm:"type:text" json:"content"`
	Banner           datatypes.JSON `gorm:"type:jsonb" json:"banner"`
	VideoPreview     datatypes.JSON `gorm:"type:jsonb" json:"video_preview"`
	CategoryId       string         `gorm:"size:255" json:"category_id"`
	Level            string         `gorm:"size:100" json:"level"`
	IsCompleted      bool           `gorm:"default:false" json:"is_completed"`
	LecturerId       uint           `json:"lecturer_id"`
	Lecturer         *User          `gorm:"foreignKey=LecturerId" json:"lecturer"`
	Status           CourseStatus   `sql:"type:ENUM('draft', 'publish', 'preview')" gorm:"column:status"`

	Version   string `gorm:"size:50" json:"version"`
	IsPublish bool   `gorm:"default:false" json:"is_publish"`
}
