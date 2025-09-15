package courses

import (
	"gateway/guards"
	"gateway/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CourseController struct {
	service *CourseService
}

func NewCourseController(s *CourseService) *CourseController {
	return &CourseController{
		service: s,
	}
}

func (c *CourseController) RegisterRoutes(r *gin.RouterGroup) {
	courseRoutes := r.Group("/courses")
	courseRoutes.Use(middlewares.AuthMiddleware())
	{
		courseRoutes.POST("/", middlewares.RequirePermission(guards.PermCourseCRUD), c.CreateCourse)
	}
}

func (c *CourseController) CreateCourse(ctx *gin.Context) {
	var dto CreateCourseDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}
	course, err := c.service.CreateCourse(dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, course)
}
