package courses

import (
	"fmt"
	"net/http"

	"gateway/guards"
	"gateway/middlewares"

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

	lecturerID := ctx.GetUint("user_id")
	fmt.Println("LecturerID from context:", lecturerID)

	course, err := c.service.CreateCourse(dto, lecturerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, course)
}
