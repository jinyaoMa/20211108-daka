package routers

import (
	"this/server/controllers"

	"github.com/gin-gonic/gin"
)

func Auth(engine *gin.Engine) {
	g := engine.Group("/auth")
	{
		g.POST("/login", controllers.Login)
	}
}
