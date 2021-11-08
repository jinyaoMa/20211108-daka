package routers

import (
	"this/server/controllers"
	"this/server/middlewares"

	"github.com/gin-gonic/gin"
)

func TestRouter(engine *gin.Engine) {
	g := engine.Group("/test")
	{
		g.GET("/getToken", controllers.TestJwtGetToken)
		g.GET("/checkToken", middlewares.Auth(), controllers.TestJwtCheckToken)
	}
}
