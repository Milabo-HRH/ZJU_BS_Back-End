package routers

import (
	"ZJU_BS_Back-End/controller"
	"ZJU_BS_Back-End/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {

	r.Use(middleware.CorsMiddleware())

	r.POST("/user/register", controller.Register)
	r.POST("/user/login", controller.Login)
	r.GET("/user/me", middleware.AuthMiddleware, controller.Info)

	r.POST("/pics", middleware.AuthMiddleware, middleware.AuthPrivilege("normal"), controller.FileUpload)
	r.GET("/pics/:filename", middleware.AuthPrivilege("guest"), controller.GetPicture)

	r.GET("/tasks", middleware.AuthPrivilege("guest"), controller.GetTasks)
	r.GET("/annotations", middleware.AuthPrivilege("guest"), controller.GetAnnotations)
	r.POST("/tasks/publish", middleware.AuthMiddleware, middleware.AuthPrivilege("normal"), controller.PublishTask)
	r.POST("/annotations/publish", middleware.AuthMiddleware, middleware.AuthPrivilege("normal"), controller.PublishAnnotation)
	r.POST("/annotations/review", middleware.AuthMiddleware, middleware.AuthPrivilege("important"), controller.ReviewAnnotation)
	r.DELETE("/tasks", middleware.AuthMiddleware, middleware.AuthPrivilege("normal"), controller.DeleteTask)
	return r
}
