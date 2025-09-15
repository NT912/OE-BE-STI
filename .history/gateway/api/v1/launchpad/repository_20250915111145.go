package lauchpad

import "gorm.io/gorm"

type CourseRepository struct {
	DB *gorm.DB
}
