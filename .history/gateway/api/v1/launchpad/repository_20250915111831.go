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

func (c *CourseRepository) CourseCreate(course *models.Course) error {
	return c.DB.Create(course).Error
}

func (c *CourseRepository) FindById(id string) (*models.Course, error) {
	var course models.Course
}
