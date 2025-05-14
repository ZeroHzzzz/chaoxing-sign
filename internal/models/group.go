package models

type RoleType int

var (
	RoleTypeMember  RoleType = 0 // 普通成员
	RoleTypeCaptain RoleType = 1 // 组长
)

type Group struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"unique;not null"`
	CaptainID int    `json:"captain_id" gorm:"not null;index"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

type GroupMembership struct {
	ID      int      `json:"id" gorm:"primaryKey"`
	GroupID int      `json:"group_id" gorm:"not null;index"`
	UserID  int      `json:"user_id" gorm:"not null;index"`
	Role    RoleType `json:"role" gorm:"not null"`

	// 添加唯一索引
	_ string `gorm:"uniqueIndex:idx_group_user"`
}
