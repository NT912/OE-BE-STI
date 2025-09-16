package launchpad

import (
	"net/http"
	"strconv"

	"gateway/guards"
	"gateway/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		lpRoutes.GET("/", c.GetLaunchpads)
		lpRoutes.GET("/:id", c.GetLaunchpadByID)
	}

	// authentication actions
	lpRoutesAuth := r.Group("/launchpads")
	lpRoutesAuth.Use(middlewares.AuthMiddleware())
	{
		// create & admin actions require permission
		lpRoutesAuth.POST("/", middlewares.RequirePermission(guards.PermCourseCRUD), c.CreateLaunchpad)
		lpRoutesAuth.POST("/:id/approve", middlewares.RequirePermission(guards.PermCourseCRUD), c.ApproveLaunchpad)

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

func (c *LaunchpadController) GetLaunchpadByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid launchpad",
		})
		return
	}

	launchpad, err := c.service.GetLaunchpadByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "launchpad not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, launchpad)
}

func (c *LaunchpadController) GetLaunchpads(ctx *gin.Context) {
	launchpads, err := c.service.GetLaunchpads()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, launchpads)
}

func (c *LaunchpadController) ApproveLaunchpad(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid launchpad",
		})
		return
	}

	launchpad, err := c.service.ApproveLaunchpad(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "launchpad not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, launchpad)
}
