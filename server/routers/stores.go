package routers

import (
	"this/server/controllers"

	"github.com/gin-gonic/gin"
)

func Stores(engine *gin.Engine) {
	g := engine.Group("/stores")
	{
		g.GET("/list", controllers.GetStoresList)
	}
}
