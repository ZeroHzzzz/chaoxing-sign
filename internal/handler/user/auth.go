package user

import (
	"chaoxing/internal/pkg/utils"
	"chaoxing/internal/pkg/xerr"
	"chaoxing/internal/services"

	"github.com/gin-gonic/gin"
)

type registerReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code"`
}

type sendCodeReq struct {
	Email string `json:"email" binding:"required,email"`
}

type loginByEmailReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginByIDReq struct {
	ID       int    `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SendVerificationCode 发送邮箱验证码
func SendVerificationCode(c *gin.Context) {
	var req sendCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		xerr.AbortWithException(c, xerr.ParamError, err)
		return
	}

	// 检查邮箱是否已被注册
	if user, _ := services.GetUserByEmail(c, req.Email); user != nil {
		xerr.AbortWithException(c, xerr.EmailVerifyErr, nil)
		return
	}

	if err := services.SendVerificationCode(c, req.Email); err != nil {
		xerr.AbortWithException(c, xerr.EmailVerifyErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

// Register 注册新用户
func RegisterByEmail(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		xerr.AbortWithException(c, xerr.ParamError, err)
		return
	}

	// 检查邮箱是否已被注册
	if user, _ := services.GetUserByEmail(c, req.Email); user != nil {
		xerr.AbortWithException(c, xerr.EmailVerifyErr, nil)
		return
	}

	if err := services.RegisterByEmail(c, req.Username, req.Email, req.Password, req.Code); err != nil {
		xerr.AbortWithException(c, xerr.RegisterErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

func RegisterTest(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		xerr.AbortWithException(c, xerr.ParamError, err)
		return
	}

	// 检查邮箱是否已被注册
	if user, _ := services.GetUserByEmail(c, req.Email); user != nil {
		xerr.AbortWithException(c, xerr.EmailVerifyErr, nil)
		return
	}

	if err := services.RegisterTest(c, req.Username, req.Email, req.Password); err != nil {
		xerr.AbortWithException(c, xerr.RegisterErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

// LoginByEmail 邮箱登录
func LoginByEmail(c *gin.Context) {
	var req loginByEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		xerr.AbortWithException(c, xerr.ParamError, err)
		return
	}

	token, user, err := services.LoginByEmail(c, req.Email, req.Password)
	if err != nil {
		xerr.AbortWithException(c, xerr.NotLogin, err)
		return
	}

	utils.JsonSuccessResponse(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// LoginByID ID登录
func LoginByID(c *gin.Context) {
	var req loginByIDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		xerr.AbortWithException(c, xerr.ParamError, err)
		return
	}

	token, user, err := services.LoginByID(c, req.ID, req.Password)
	if err != nil {
		xerr.AbortWithException(c, xerr.NotLogin, err)
		return
	}

	utils.JsonSuccessResponse(c, gin.H{
		"token": token,
		"user":  user,
	})
}
