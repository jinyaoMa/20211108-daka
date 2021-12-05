package controllers

import (
	"net/http"
	"this/database"
	"this/database/models"
	"this/gate"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateTimesheetForm struct {
	StoreID     int       `form:"StoreID"`
	ID          uint      `form:"ID" binding:"required"`
	SigninTime  time.Time `form:"SigninTime" binding:"required"`
	SignoutTime time.Time `form:"SignoutTime" binding:"required"`
	Total       float64   `form:"Total" binding:"required"`
}

// @Summary UpdateTimesheet
// @Description Update Timesheet
// @Tags After Authorization
// @accept x-www-form-urlencoded
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @param StoreID query int false "Store ID"
// @param ID formData uint true "ID"
// @param SigninTime formData time.Time true "SigninTime"
// @param SignoutTime formData time.Time true "SignoutTime"
// @param Total formData float64 true "Total"
// @Success 200 "{ ok }"
// @Failure 401 "Auth failed"
// @Failure 404 "{ error }"
// @Router /timesheet/update/ [post]
func UpdateTimesheet(c *gin.Context) {
	temp, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Middleware Auth not exists!",
		})
		return
	}
	claims := temp.(*gate.Claims)

	if !(claims.Usertype == "Admin" || claims.Usertype == "Office") {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usertype " + claims.Usertype + " is unavailable!",
		})
		return
	}

	var form UpdateTimesheetForm
	if c.ShouldBind(&form) != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Check for empty fields!",
		})
		return
	}

	store, ok := database.GetStore(form.StoreID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Store error!",
		})
		return
	}

	var timesheet *models.Timesheet
	_, err := timesheet.UpdateById(*store, form.ID, form.SigninTime, form.SignoutTime, form.Total)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Timesheet error!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}
