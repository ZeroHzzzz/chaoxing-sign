package router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 用户相关路由
	// userGroup := r.Group("/user")
	{
		// 无需认证的接口
		// userGroup.POST("/register", handler.RegisterHandler)
		// userGroup.POST("/login", handler.LoginHandler)
		// userGroup.GET("/info/:id", handler.GetUserHandler)
		// userGroup.GET("/list", handler.GetUsersHandler)
		// userGroup.DELETE("/:id", handler.DeleteUserHandler)

		// 需要认证的接口
		// auth := userGroup.Group("/", middleware.JWTAuth())
		// {
		// 	// 超星账号相关
		// 	auth.POST("/chaoxing/bind", handler.BindChaoxingAccountHandler)
		// 	auth.GET("/chaoxing/info", handler.GetChaoxingAccountHandler)
		// 	auth.PUT("/chaoxing/update", handler.UpdateChaoxingAccountHandler)
		// 	auth.DELETE("/chaoxing/unbind", handler.UnbindChaoxingAccountHandler)
		// }
	}

	return r
}
