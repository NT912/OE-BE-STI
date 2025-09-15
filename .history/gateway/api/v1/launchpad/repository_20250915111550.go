package lauchpad

import (
	"gateway/models"

	"gorm.io/gorm"
)

type CourseRepository struct {
	DB *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{DB: db}
}

func (c *CourseRepository) CourseCreate(course *models.Course)
