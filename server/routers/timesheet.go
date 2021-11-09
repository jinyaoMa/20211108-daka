package routers

import (
	"this/server/controllers"
	"this/server/middlewares"

	"github.com/gin-gonic/gin"
)

func Timesheet(engine *gin.Engine) {
	g := engine.Group("/timesheet")
	{
		g.GET("/list", middlewares.Auth(), controllers.GetTimesheetList)
		g.GET("/download", middlewares.Auth(), controllers.DownloadExcel)

		g.POST("/update", middlewares.Auth(), controllers.UpdateTimesheet)
	}
}
