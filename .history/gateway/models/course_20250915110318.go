package models

import "time"

type Course struct {
	ID          string       `json:"id" gorm:"primaryKey"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	CreatorID   string       `json:"creatorId"`
	CreatedAt   time.Time    `json:"createdAt"`
	Sections    []Section    `json:"sections" gorm:"foreignKey:CourseID"`
	VotingPlan  []VotingPlan `json:"votingPlan" gorm:"foreignKey:CourseID"`
	Investment  Investment   `json:"investment"`
	Status      string       `json:"status"` // pending, approved, upcoming, ongoing, ended
}

type Section struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	CourseID string `json:"courseId"`
	Title    string `json:"title"`
}

type VotingPlan struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	CourseID string    `json:"courseId"`
	Section  int       `json:"section"`
	Date     time.Time `json:"date"`
}

type Investment struct {
	TotalInvested float64 `json:"totalInvested"`
	Currency      string  `json:"currency"`
	Investors     int     `json:"investors"`
	Progress      float64 `json:"progress"` // 0% to 100%
	Goal          float64 `json:"goal"`     // Mục tiêu đầu tư
}
