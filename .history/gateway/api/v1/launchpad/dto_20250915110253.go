package lauchpad

import "time"

type CreateCourseDTO struct {
	Title          string          `json:"title" validate:"required"`
	Description    string          `json:"description" validate:"required"`
	Sections       []SectionDTO    `json:"sections" validate:"required"`
	VotingPlan     []VotingPlanDTO `json:"votingPlan" validate:"required"`
	InvestmentGoal float64         `json:"investmentGoal" validate:"required,gt=0"`
}

type SectionDTO struct {
	Title string `json:"title" validate:"required"`
}

type VotingPlanDTO struct {
	Section int       `json:"section" validate:"required"`
	Date    time.Time `json:"date" validate:"required"`
}

type InvestDTO struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}
