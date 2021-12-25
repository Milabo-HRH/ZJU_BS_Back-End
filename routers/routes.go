package routers

import (
	"awesomeProject/controller"
	"awesomeProject/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/user/register", controller.Register)
	r.POST("/user/login", controller.Login)
	r.GET("/user/me", middleware.AuthMiddleware(), controller.Info)
	return r
}
