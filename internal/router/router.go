package router

import (
	"chaoxing/internal/handler/user"
	"chaoxing/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 用户相关路由
	userGroup := r.Group("/user")
	{
		// 无需认证的接口
		userGroup.POST("/register", user.Register)
		userGroup.POST("/login", user.Login)
		userGroup.GET("/info/:id", user.GetUserInfo)
		userGroup.GET("/list", user.GetUserList)
		userGroup.DELETE("/:id", user.DeleteUser)

		// 需要认证的接口
		auth := userGroup.Group("/", middleware.JWTAuth())
		{
			// 超星账号相关
			auth.POST("/chaoxing/bind", user.BindChaoxingAccount)
			auth.GET("/chaoxing/info", user.GetChaoxingAccount)
			auth.PUT("/chaoxing/update", user.UpdateChaoxingAccount)
			auth.DELETE("/chaoxing/unbind", user.UnbindChaoxingAccount)
		}
	}

	return r
}
