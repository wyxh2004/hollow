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

	// 公开的话题箱路由
	api := r.Group("/api")
	{
		// 公开接口，不需要认证
		api.GET("/boxes", boxController.ListBoxes)
		api.GET("/boxes/:id", boxController.GetBox)
		api.POST("/boxes/:id/messages", boxController.CreateMessage)

		// 需要认证的接口
		authorized := api.Group("")
		authorized.Use(middlewares.AuthMiddleware())
		{
			authorized.POST("/boxes", boxController.CreateBox)
			authorized.POST("/messages/:id/like", boxController.LikeMessage)
		}
	}
}
