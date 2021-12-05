package routers

import (
	"this/server/controllers"
	"this/server/middlewares"

	"github.com/gin-gonic/gin"
)

func Stores(engine *gin.Engine) {
	g := engine.Group("/stores")
	{
		g.GET("/list", middlewares.Auth(), controllers.GetStoresList)
	}
}
