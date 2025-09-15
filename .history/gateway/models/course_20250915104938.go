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
	Funding     Funding      `json:"funding"`
	Status      string       `json:"status"` // pending, approved, upcoming, ongoing, ended, success, rejected
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

type Funding struct {
	FundedAmount float64 `json:"fundedAmount"`
	Currency     string  `json:"currency"`
	Backers      int     `json:"backers"`
	Progress     float64 `json:"progress"`
	Goal         float64 `json:"goal"`
}
