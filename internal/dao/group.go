package dao

import (
	"chaoxing/internal/models"
	"context"
)

func (d *Dao) NewGroup(ctx context.Context, group *models.Group) error {
	return d.DB.Create(group).Error
}

func (d *Dao) GetGroupByID(ctx context.Context, id int) (*models.Group, error) {
	var group models.Group
	err := d.DB.Where("id = ?", id).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (d *Dao) UpdateGroup(ctx context.Context, group *models.Group) error {
	return d.DB.Model(group).Updates(group).Error
}

func (d *Dao) DeleteGroup(ctx context.Context, id int) error {
	return d.DB.Delete(&models.Group{}, id).Error
}

func (d *Dao) AddGroupMembership(ctx context.Context, member *models.GroupMembership) error {
	return d.DB.Create(member).Error
}

func (d *Dao) RemoveGroupMembership(ctx context.Context, groupID, userID int) error {
	return d.DB.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&models.GroupMembership{}).Error
}

func (d *Dao) UpdateGroupMembership(ctx context.Context, member *models.GroupMembership) error {
	return d.DB.Model(member).Updates(member).Error
}

func (d *Dao) GetGroupMemberships(ctx context.Context, groupID int) ([]*models.GroupMembership, error) {
	var memberships []*models.GroupMembership
	err := d.DB.Where("group_id = ?", groupID).Find(&memberships).Error
	if err != nil {
		return nil, err
	}
	return memberships, nil
}

func (d *Dao) GetGroupMembership(ctx context.Context, groupID, userID int) (*models.GroupMembership, error) {
	var membership models.GroupMembership
	err := d.DB.Where("group_id = ? AND user_id = ?", groupID, userID).First(&membership).Error
	if err != nil {
		return nil, err
	}
	return &membership, nil
}

func (d *Dao) GetGroupsByUserID(ctx context.Context, userID int) ([]*models.Group, error) {
	var groups []*models.Group
	err := d.DB.Joins("JOIN group_memberships ON groups.id = group_memberships.group_id").
		Where("group_memberships.user_id = ?", userID).
		Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (d *Dao) GetGroupMembersByGroupID(ctx context.Context, groupID int) ([]*models.User, error) {
	var members []*models.User
	err := d.DB.Joins("JOIN group_memberships ON users.id = group_memberships.user_id").
		Where("group_memberships.group_id = ?", groupID).
		Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}
