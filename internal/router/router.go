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
		userGroup.POST("/send-code", user.SendVerificationCode) // 发送验证码
		userGroup.POST("/register", user.Register)              // 注册
		userGroup.POST("/login/id", user.LoginByID)             // ID登录
		userGroup.POST("/login/email", user.LoginByEmail)       // 邮箱登录
		userGroup.GET("/info/:id", user.GetUserInfo)            // 获取用户信息
		userGroup.GET("/list", user.GetUserList)                // 获取用户列表
		userGroup.DELETE("/:id", user.DeleteUser)               // 删除用户

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
