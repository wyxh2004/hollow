package routes

import (
	"hollow/controllers"
	"hollow/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	userController := &controllers.UserController{}
	boxController := &controllers.BoxController{}

	// 用户认证路由
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", userController.Register)
		auth.POST("/login", userController.Login)
	}

	// 留言箱路由
	api := r.Group("/api")
	api.Use(middlewares.AuthMiddleware())
	{
		api.POST("/boxes", boxController.CreateBox)
		api.GET("/boxes", boxController.ListBoxes)
		api.GET("/boxes/:id", boxController.GetBox)
		api.POST("/boxes/:id/messages", boxController.CreateMessage)
		api.POST("/messages/:id/like", boxController.LikeMessage)
	}
}
