package courses

type CreateCourseDTO struct {
	Name             string  `json:"name" binding:"required,max=255"`
	ShortDescription string  `json:"short_description" binding:"required,max=255"`
	Goal             float64 `json:"goal" binding:"required"`
}
