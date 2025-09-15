package launchpad

import (
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
