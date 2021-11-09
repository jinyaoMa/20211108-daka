package controllers

import (
	"net/http"
	"this/database"
	"this/database/models"
	"this/gate"

	"github.com/gin-gonic/gin"
)

type GetTimesheetListQuery struct {
	Offset  int    `form:"offset"`
	Limit   int    `form:"limit"`
	Order   string `form:"order"`
	OrderBy string `form:"orderby"`
}

// @Summary GetTimesheetList
// @Description Get Timesheet List
// @Tags After Authorization
// @accept json
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @param offset query int false "Offset"
// @param limit query int false "Limit"
// @param order query int false "Order"
// @param orderby query int false "OrderBy"
// @Success 200 "{ ok , data , count }"
// @Failure 401 "Auth failed"
// @Failure 404 "{ error }"
// @Router /timesheet/list [get]
func GetTimesheetList(c *gin.Context) {
	temp, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Middleware Auth not exists!",
		})
		return
	}
	claims := temp.(*gate.Claims)

	var form GetTimesheetListQuery
	if c.ShouldBind(&form) != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Check for empty fields!",
		})
		return
	}

	store, ok := database.GetStore(claims.StoreID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Store error!",
		})
		return
	}

	var timesheet *models.Timesheet
	list, err := timesheet.GetListByOffsetLimit(*store, form.Offset, form.Limit, form.OrderBy, form.Order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Timesheet error!",
		})
		return
	}

	count, errCount := timesheet.GetCount(*store)
	if errCount != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Timesheet error!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":    true,
		"data":  list,
		"count": count,
	})
}
