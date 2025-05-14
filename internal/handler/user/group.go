package user

import (
	"chaoxing/internal/pkg/utils"
	"chaoxing/internal/pkg/xerr"
	"chaoxing/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createGroupReq struct {
	Name string `json:"name" binding:"required"`
}

type updateGroupReq struct {
	Name string `json:"name" binding:"required"`
}

type addMemberReq struct {
	UserID int `json:"user_id" binding:"required"`
}

type transferCaptainReq struct {
	NewCaptainID int `json:"new_captain_id" binding:"required"`
}

// CreateGroup 创建分组
func CreateGroup(c *gin.Context) {
	var req createGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		xerr.AbortWithException(c, xerr.ParamError, err)
		return
	}

	userID := c.GetInt("userID")
	if err := services.NewGroup(c, req.Name, userID); err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

// GetGroupInfo 获取分组信息
func GetGroupInfo(c *gin.Context) {
	groupID, _ := strconv.Atoi(c.Param("id"))

	group, err := services.GetGroupByGroupID(c, groupID)
	if err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}

	utils.JsonSuccessResponse(c, group)
}

// GetMyGroups 获取我的分组列表
func GetMyGroups(c *gin.Context) {
	userID := c.GetInt("userID")

	groups, err := services.GetGroupsByUserID(c, userID)
	if err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}

	utils.JsonSuccessResponse(c, groups)
}

// UpdateGroup 更新分组信息
func UpdateGroup(c *gin.Context) {
	groupID, _ := strconv.Atoi(c.Param("id"))
	var req updateGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		xerr.AbortWithException(c, xerr.ParamError, err)
		return
	}

	// 验证权限
	userID := c.GetInt("userID")
	isCaptain, err := services.CheckGroupCaptain(c, groupID, userID)
	if err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}
	if !isCaptain {
		xerr.AbortWithException(c, xerr.PremissionDenied, nil)
		return
	}

	if err := services.UpdateGroupProfile(c, groupID, req.Name); err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

// DeleteGroup 删除分组
func DeleteGroup(c *gin.Context) {
	groupID, _ := strconv.Atoi(c.Param("id"))

	// 验证权限
	userID := c.GetInt("userID")
	isCaptain, err := services.CheckGroupCaptain(c, groupID, userID)
	if err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}
	if !isCaptain {
		xerr.AbortWithException(c, xerr.PremissionDenied, nil)
		return
	}

	if err := services.DeleteGroup(c, groupID); err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

// AddGroupMember 添加分组成员
func AddGroupMember(c *gin.Context) {
	groupID, _ := strconv.Atoi(c.Param("id"))
	var req addMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		xerr.AbortWithException(c, xerr.ParamError, err)
		return
	}

	if err := services.AddGroupMembership(c, groupID, req.UserID); err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

// RemoveGroupMember 移除分组成员
func RemoveGroupMember(c *gin.Context) {
	groupID, _ := strconv.Atoi(c.Param("id"))
	memberID, _ := strconv.Atoi(c.Param("member_id"))

	// 验证权限
	userID := c.GetInt("userID")
	isCaptain, err := services.CheckGroupCaptain(c, groupID, userID)
	if err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}
	if !isCaptain && userID != memberID { // 组长可以移除任何人，成员只能移除自己
		xerr.AbortWithException(c, xerr.PremissionDenied, nil)
		return
	}

	if err := services.RemoveGroupMembership(c, groupID, memberID); err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

// GetGroupMembers 获取分组成员列表
func GetGroupMembers(c *gin.Context) {
	groupID, _ := strconv.Atoi(c.Param("id"))

	members, err := services.GetGroupMembersByGroupID(c, groupID)
	if err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}

	utils.JsonSuccessResponse(c, members)
}

// TransferCaptain 转移组长
func TransferCaptain(c *gin.Context) {
	groupID, _ := strconv.Atoi(c.Param("id"))
	var req transferCaptainReq
	if err := c.ShouldBindJSON(&req); err != nil {
		xerr.AbortWithException(c, xerr.ParamError, err)
		return
	}

	userID := c.GetInt("userID")
	if err := services.TransferGroupCaptain(c, groupID, userID, req.NewCaptainID); err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

// GetGroupInviteCode 获取群组邀请码
func GetGroupInviteCode(c *gin.Context) {
	groupID, _ := strconv.Atoi(c.Param("id"))

	// 验证权限
	userID := c.GetInt("userID")
	isCaptain, err := services.CheckGroupCaptain(c, groupID, userID)
	if err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}
	if !isCaptain {
		xerr.AbortWithException(c, xerr.PremissionDenied, nil)
		return
	}

	group, err := services.GetGroupByGroupID(c, groupID)
	if err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}

	utils.JsonSuccessResponse(c, gin.H{
		"invite_code": group.InviteCode,
	})
}

// JoinGroupByInviteCode 通过邀请码加入群组
func JoinGroupByInviteCode(c *gin.Context) {
	var req struct {
		InviteCode string `json:"invite_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		xerr.AbortWithException(c, xerr.ParamError, err)
		return
	}

	userID := c.GetInt("userID")
	if err := services.JoinGroupByInviteCode(c, req.InviteCode, userID); err != nil {
		xerr.AbortWithException(c, xerr.ServerErr, err)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
