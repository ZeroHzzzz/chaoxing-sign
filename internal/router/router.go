package router

import (
	"chaoxing/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.ErrHandler())

	// 用户相关路由
	// userGroup := r.Group("/user")
	// {
	// 	// 无需认证的接口
	// 	userGroup.POST("/send-code", user.SendVerificationCode) // 发送验证码
	// 	userGroup.POST("/register/email", user.RegisterByEmail) // 注册
	// 	userGroup.POST("/register/test", user.RegisterTest)     // 注册测试
	// 	userGroup.POST("/login/id", user.LoginByID)             // ID登录
	// 	userGroup.POST("/login/email", user.LoginByEmail)       // 邮箱登录
	// 	userGroup.GET("/info/:id", user.GetUserInfo)            // 获取用户信息
	// 	userGroup.GET("/list", user.GetUserList)                // 获取用户列表
	// 	userGroup.DELETE("/:id", user.DeleteUser)               // 删除用户

	// 	// 需要认证的接口
	// 	auth := userGroup.Group("/", middleware.JWTAuth())
	// 	{
	// 		// 超星账号相关
	// 		auth.POST("/chaoxing/bind", user.BindChaoxingAccount)
	// 		auth.GET("/chaoxing/info", user.GetChaoxingAccount)
	// 		auth.GET("/chaoxing/list", user.GetChaoxingAccounts) // 新增：获取超星账号列表
	// 		auth.PUT("/chaoxing/update", user.UpdateChaoxingAccount)
	// 		auth.DELETE("/chaoxing/unbind", user.UnbindChaoxingAccount)

	// 		// 分组相关
	// 		groupGroup := auth.Group("/group")
	// 		{
	// 			groupGroup.POST("", user.CreateGroup)                               // 创建分组
	// 			groupGroup.GET("/my", user.GetMyGroups)                             // 获取我的分组列表
	// 			groupGroup.GET("/:id", user.GetGroupInfo)                           // 获取分组信息
	// 			groupGroup.PUT("/:id", user.UpdateGroup)                            // 更新分组信息
	// 			groupGroup.DELETE("/:id", user.DeleteGroup)                         // 删除分组
	// 			groupGroup.POST("/:id/member", user.AddGroupMember)                 // 添加分组成员
	// 			groupGroup.DELETE("/:id/member/:member_id", user.RemoveGroupMember) // 移除分组成员
	// 			groupGroup.GET("/:id/members", user.GetGroupMembers)                // 获取分组成员列表
	// 			groupGroup.POST("/:id/transfer", user.TransferCaptain)              // 转移组长
	// 			groupGroup.GET("/:id/invite-code", user.GetGroupInviteCode)         // 获取群组邀请码
	// 			groupGroup.POST("/join", user.JoinGroupByInviteCode)                // 通过邀请码加入群组
	// 		}
	// 	}
	// }

	return r
}
