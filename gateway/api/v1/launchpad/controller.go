package launchpad

import (
	"net/http"

	"gateway/guards"
	"gateway/middlewares"

	"github.com/gin-gonic/gin"
)

type LaunchpadController struct {
	service *LaunchpadService
}

func NewLaunchpadController(s *LaunchpadService) *LaunchpadController {
	return &LaunchpadController{service: s}
}

func (c *LaunchpadController) RegisterRoutes(r *gin.RouterGroup) {
	lpRoutes := r.Group("/launchpads")
	{
		lpRoutes.GET("/")
	}

	// authentication actions
	lpRoutesAuth := r.Group("/launchpads")
	lpRoutesAuth.Use(middlewares.AuthMiddleware())
	{
		// create & admin actions require permission
		lpRoutesAuth.POST("/", middlewares.RequirePermission(guards.PermCourseCRUD), c.CreateLaunchpad)
		lpRoutesAuth.POST("/:id/approve", middlewares.RequirePermission(guards.PermCourseCRUD))

		// invest: any logged-in user can invest (as you wanted)
	}
}

func (c *LaunchpadController) CreateLaunchpad(ctx *gin.Context) {
	var dto CreateLaunchpadDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	lp, err := c.service.CreateLaunchpad(dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, lp)
}
