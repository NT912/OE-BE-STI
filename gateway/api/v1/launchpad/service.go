package launchpad

import (
	"errors"
	"time"

	"gateway/models"
)

type LaunchpadService struct {
	repo *LaunchpadRepository
}

func NewLaunchpadService(r *LaunchpadRepository) *LaunchpadService {
	return &LaunchpadService{repo: r}
}

func (s *LaunchpadService) CreateLaunchpad(dto CreateLaunchpadDTO) (*models.Launchpad, error) {
	// Check course exist
	ok, err := s.repo.CourseExists(dto.CourseID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("Course not found")
	}

	// Create Launchpad object
	lp := &models.Launchpad{
		CourseID:    dto.CourseID,
		Title:       dto.Title,
		Description: dto.Description,
		FundingGoal: dto.FundingGoal,
		Funded:      0,
		Backers:     0,
		Approved:    false, // admin will approve
		Status:      models.LaunchpadUpcoming,
	}

	if err := s.repo.Create(lp); err != nil {
		return nil, err
	}

	// create voting plans if provided
	if len(dto.VotingPlans) > 0 {
		var plans []models.VotingPlan
		for _, p := range dto.VotingPlans {
			t, parseErr := parseDateFlexible(p.ScheduleAt)
			if parseErr != nil {
				return nil, parseErr
			}
			plans = append(plans, models.VotingPlan{
				LaunchpadID: lp.ID,
				Step:        p.Step,
				Sections:    p.Sections,
				ScheduleAt:  t,
				Title:       p.Title,
			})
		}
		if err := s.repo.CreateVotingPlans(lp.ID, plans); err != nil {
			return nil, err
		}
		// optionally set next voting date
		if len(plans) > 0 {
			t := plans[0].ScheduleAt
			_ = s.repo.UpdateNextVotingAt(lp.ID, &t)
		}
	}

	return s.repo.FindByID(lp.ID)
}

func (s *LaunchpadService) GetLaunchpadByID(id uint) (*models.Launchpad, error) {
	return s.repo.FindByID(id)
}

func (s *LaunchpadService) GetLaunchpads() ([]models.Launchpad, error) {
	return s.repo.FindAll(true)
}

// Goal will be input by investor
func calculateStatus(goal, funded float64) models.LaunchpadStatus {
	if funded <= 0 || goal <= 0 {
		return models.LaunchpadUpcoming
	}
	percent := (funded / goal) * 100
	if percent >= 80 {
		return models.LaunchpadSuccess
	}
	return models.LaunchpadFeaturing
}

func parseDateFlexible(s string) (time.Time, error) {
	// try ISO YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}

	// try M/D/YYYY like "3/6/2025"
	if t, err := time.Parse("2/1/2006", s); err == nil {
		return t, nil
	}

	// try RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	return time.Time{}, errors.New("invalid date format, use YYYY-MM-DD")
}
