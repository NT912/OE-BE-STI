package launchpad

import (
	"time"

	"gateway/models"

	"gorm.io/gorm"
)

type LaunchpadRepository struct {
	db *gorm.DB
}

func NewLaunchpadRepository(db *gorm.DB) *LaunchpadRepository {
	return &LaunchpadRepository{db: db}
}

func (r *LaunchpadRepository) Create(lp *models.Launchpad) error {
	return r.db.Create(lp).Error
}

func (r *LaunchpadRepository) CreateVotingPlans(lpID uint, plans []models.VotingPlan) error {
	for i := range plans {
		plans[i].LaunchpadID = lpID
	}

	return r.db.Create(&plans).Error
}

func (r *LaunchpadRepository) FindByID(id uint) (*models.Launchpad, error) {
	var lp models.Launchpad
	err := r.db.Preload("Course").Preload("VotingPlans", func(db *gorm.DB) *gorm.DB {
		return db.Order("step asc")
	}).First(&lp, id).Error
	return &lp, err
}

func (r *LaunchpadRepository) Update(lp *models.Launchpad) error {
	// update UpdatedAt automatically via Save
	return r.db.Save(lp).Error
}

func (r *LaunchpadRepository) FindAllByStatus(status models.LaunchpadStatus) ([]models.Launchpad, error) {
	var list []models.Launchpad
	err := r.db.Preload("Course").Preload("VotingPlans", func(db *gorm.DB) *gorm.DB {
		return db.Order("step asc")
	}).Where("status = ? AND approved = ?", status, true).Find(&list).Error

	return list, err
}

func (r *LaunchpadRepository) CourseExists(courseID uint) (bool, error) {
	var c models.Course
	err := r.db.Select("id").First(&c, courseID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *LaunchpadRepository) UpdateNextVotingAt(lpID uint, t *time.Time) error {
	return r.db.Model(&models.Launchpad{}).Where("id=?", lpID).Update("next_voting_at", t).Error
}
