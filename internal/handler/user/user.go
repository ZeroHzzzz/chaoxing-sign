package user

import (
	"chaoxing/internal/pkg/utils"
	"chaoxing/internal/pkg/xerr"
	"chaoxing/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type registerReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type bindChaoxingReq struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		xerr.AbortWithException(c, xerr.ParamError, err)
		return
	}

	if err := services.Register(c, req.Username, req.Password); err != nil {
		xerr.AbortWithException(c, xerr.RegisterErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

func Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		xerr.AbortWithException(c, xerr.ParamError, err)
		return
	}

	user, err := services.GetUserByEmailPass(c, req.Email, req.Password)
	if err != nil {
		xerr.AbortWithException(c, xerr.UserNotFind, err)
		return
	}

	token, err := services.Login(c, user.ID, req.Password)
	if err != nil {
		xerr.AbortWithException(c, xerr.NotLogin, err)
		return
	}

	utils.JsonSuccessResponse(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func GetUserInfo(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	user, err := services.GetUserByID(c, userID)
	if err != nil {
		xerr.AbortWithException(c, xerr.UserNotFind, err)
		return
	}

	utils.JsonSuccessResponse(c, gin.H{
		"user": user,
	})
}

func GetUserList(c *gin.Context) {
	username := c.Query("username")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := services.GetUsersByUsername(c, username, page, pageSize)
	if err != nil {
		xerr.AbortWithException(c, xerr.UserNotFind, err)
		return
	}

	utils.JsonSuccessResponse(c, gin.H{
		"total": total,
		"list":  users,
	})
}

func DeleteUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	password := c.Query("password")

	if _, err := services.GetUserByID(c, userID); err != nil {
		xerr.AbortWithException(c, xerr.UserNotFind, err)
		return
	}
	if err := services.DeleteUserByPass(c, userID, password); err != nil {
		xerr.AbortWithException(c, xerr.DeleteUserErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

func BindChaoxingAccount(c *gin.Context) {
	var req bindChaoxingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		xerr.AbortWithException(c, xerr.ParamError, err)
		return
	}

	userID := c.GetInt("user_id") // 从JWT中获取
	if err := services.BindChaoxingAccount(c, userID, req.Phone, req.Password); err != nil {
		xerr.AbortWithException(c, xerr.ChaoxingOperateErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

func GetChaoxingAccount(c *gin.Context) {
	userID := c.GetInt("user_id") // 从JWT中获取

	account, err := services.GetUserChaoxingAccount(c, userID)
	if err != nil {
		xerr.AbortWithException(c, xerr.ChaoxingOperateErr, err)
		return
	}

	utils.JsonSuccessResponse(c, account)
}

func UpdateChaoxingAccount(c *gin.Context) {
	var req bindChaoxingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		xerr.AbortWithException(c, xerr.ParamError, err)
		return
	}

	userID := c.GetInt("user_id") // 从JWT中获取
	if err := services.UpdateChaoxingAccount(c, userID, req.Phone, req.Password); err != nil {
		xerr.AbortWithException(c, xerr.ChaoxingOperateErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

func UnbindChaoxingAccount(c *gin.Context) {
	userID := c.GetInt("user_id") // 从JWT中获取

	if err := services.UnbindChaoxingAccount(c, userID); err != nil {
		xerr.AbortWithException(c, xerr.ChaoxingOperateErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
