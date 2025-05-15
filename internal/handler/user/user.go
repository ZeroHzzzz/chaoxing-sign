package user

import (
	"chaoxing/internal/pkg/utils"
	"chaoxing/internal/pkg/xerr"
	"chaoxing/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type bindChaoxingReq struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
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

func GetChaoxingAccounts(c *gin.Context) {
	userID := c.GetInt("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	accounts, total, err := services.GetUserBindChaoxingAccounts(c, userID, page, pageSize)
	if err != nil {
		xerr.AbortWithException(c, xerr.ChaoxingOperateErr, err)
		return
	}

	utils.JsonSuccessResponse(c, gin.H{
		"total": total,
		"list":  accounts,
	})
}

func UnbindChaoxingAccount(c *gin.Context) {
	accountID := c.Param("id")
	userID := c.GetInt("userID")

	if err := services.UnbindChaoxingAccount(c, userID, accountID); err != nil {
		xerr.AbortWithException(c, xerr.ChaoxingOperateErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

func GetUserChaoxingAccountByPhone(c *gin.Context) {
	phone := c.Param("phone")
	userID := c.GetInt("userID")

	account, err := services.GetUserChaoxingAccountByPhone(c, phone, userID)
	if err != nil {
		xerr.AbortWithException(c, xerr.ChaoxingOperateErr, err)
		return
	}

	utils.JsonSuccessResponse(c, account)
}

func DeleteUserChaoxingAccounts(c *gin.Context) {
	userID := c.GetInt("userID")

	if err := services.DeleteUserChaoxingAccounts(c, userID); err != nil {
		xerr.AbortWithException(c, xerr.ChaoxingOperateErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
