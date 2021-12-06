package routers

import (
	"this/server/controllers"
	"this/server/middlewares"

	"github.com/gin-gonic/gin"
)

func User(engine *gin.Engine) {
	g := engine.Group("/user")
	{
		g.POST("/add", middlewares.Auth(), controllers.AddUser)
	}
}
