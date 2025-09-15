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
