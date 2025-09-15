package courses

import "gateway/models"

type CourseService struct {
	repo *CourseRepository
}

func NewCourseService(repo *CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) CreateCourse(dto CreateCourseDTO, lecturerID uint) (*models.Course, error) {
	course := &models.Course{
		Name:             dto.Name,
		ShortDescription: dto.ShortDescription,
		Goal:             dto.Goal,
		LecturerId:       lecturerID,
	}
	if err := s.repo.Create(course); err != nil {
		return nil, err
	}
	return course, nil
}
