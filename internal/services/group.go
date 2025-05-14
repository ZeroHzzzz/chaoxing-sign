package services

import (
	"chaoxing/internal/models"
	"chaoxing/internal/pkg/xerr"
	"context"
)

func NewGroup(ctx context.Context, name string, captainID int) error {
	// 创建分组
	newGroup := &models.Group{
		Name:      name,
		CaptainID: captainID,
	}

	err := d.NewGroup(ctx, newGroup)
	if err != nil {
		return err
	}

	return nil
}

func GetGroupByGroupID(ctx context.Context, groupID int) (*models.Group, error) {
	// 获取分组信息
	group, err := d.GetGroupByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func GetGroupByCaptainID(ctx context.Context, captainID int) ([]*models.Group, error) {
	// 获取分组信息
	groups, err := d.GetGroupsByUserID(ctx, captainID)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func UpdateGroupProfile(ctx context.Context, groupID int, name string) error {
	group, err := d.GetGroupByGroupID(ctx, groupID)
	if err != nil {
		return err
	}

	group.Name = name

	err = d.UpdateGroup(ctx, group)
	if err != nil {
		return err
	}

	return nil
}

func DeleteGroup(ctx context.Context, groupID int) error {
	// 删除分组
	err := d.DeleteGroup(ctx, groupID)
	if err != nil {
		return err
	}

	return nil
}

func GetGroupsByUserID(ctx context.Context, userID int) ([]*models.Group, error) {
	// 获取用户所在的分组
	groups, err := d.GetGroupsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func AddGroupMembership(ctx context.Context, groupID, userID int) error {
	// 添加用户到分组
	membership := &models.GroupMembership{
		GroupID: groupID,
		UserID:  userID,
	}

	err := d.AddGroupMembership(ctx, membership)
	if err != nil {
		return err
	}

	return nil
}

func RemoveGroupMembership(ctx context.Context, groupID, userID int) error {
	// 退出分组
	err := d.RemoveGroupMembership(ctx, groupID, userID)
	if err != nil {
		return err
	}

	return nil
}

func GetGroupMemberships(ctx context.Context, groupID int) ([]*models.GroupMembership, error) {
	// 获取分组成员角色
	memberships, err := d.GetGroupMemberships(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return memberships, nil
}

func GetGroupMembership(ctx context.Context, groupID, userID int) (*models.GroupMembership, error) {
	// 获取用户在分组中的信息
	membership, err := d.GetGroupMembership(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}

	return membership, nil
}

func GetGroupMembersByGroupID(ctx context.Context, groupID int) ([]*models.User, error) {
	// 获取分组成员信息
	members, err := d.GetGroupMembersByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func TransferGroupCaptain(ctx context.Context, groupID, userID, newCaptainID int) error {
	// 获取分组信息
	group, err := d.GetGroupByGroupID(ctx, groupID)
	if err != nil {
		return err
	}

	if group.CaptainID != userID {
		return xerr.PremissionDenied
	}

	group.CaptainID = newCaptainID

	err = d.UpdateGroup(ctx, group)
	if err != nil {
		return err
	}

	return nil
}

func CheckGroupCaptain(ctx context.Context, groupID, userID int) (bool, error) {
	// 获取分组信息
	group, err := d.GetGroupByGroupID(ctx, groupID)
	if err != nil {
		return false, err
	}

	// 检查用户是否为分组队长
	if group.CaptainID == userID {
		return true, nil
	}

	return false, nil
}
